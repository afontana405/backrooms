package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Page is the top-level JSON structure
type Page struct {
	Title    string            `json:"title"`
	Theme    string            `json:"theme"`
	GameMode bool              `json:"gameMode"`
	PageType string            `json:"pageType"`
	Body     []json.RawMessage `json:"body"`
}

// Engine reads JSON and produces HTML
type Engine struct {
	Registry   map[string]Component
	PublicPath string
}

// New creates an engine with the default component registry
func New() *Engine {
	// Resolve public path relative to executable location
	exe, _ := os.Executable()
	publicPath := filepath.Join(filepath.Dir(exe), "public")

	e := &Engine{
		Registry:   make(map[string]Component),
		PublicPath: publicPath,
	}
	RegisterDefaults(e)
	RegisterScreenshot()
	return e
}

// Register adds a component to the registry
func (e *Engine) Register(name string, c Component) {
	e.Registry[name] = c
}

// RenderFile reads a JSON file and returns a complete HTML page
func (e *Engine) RenderFile(path string) (string, error) {
	data, err := ReadEmbedFile(strings.ReplaceAll(path, "\\", "/"))
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}
	return e.Render(data)
}

// Render takes raw JSON bytes and returns a complete HTML page
func (e *Engine) Render(data []byte) (string, error) {
	return e.RenderDiag(data, nil)
}

// RenderDiag renders with diagnostics collection.
func (e *Engine) RenderDiag(data []byte, diag *DiagCollector) (string, error) {
	// Validate page JSON
	if diag != nil {
		diag.ValidatePageJSON(data)
	}

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return "", fmt.Errorf("parsing page: %w", err)
	}

	var bodyHTML strings.Builder
	for _, raw := range page.Body {
		html, err := e.renderAtomDiag(raw, diag, "body", 1)
		if err != nil {
			if diag != nil {
				diag.Error("render", fmt.Sprintf("Render error: %v", err))
			}
			return "", err
		}
		bodyHTML.WriteString(html)
	}

	theme := page.Theme
	if theme == "" {
		theme = "light"
	}

	title := page.Title
	if title == "" {
		title = "ChefScript"
	}

	body := bodyHTML.String()

	// Inject diagnostics panel
	if diag != nil {
		errs, warns, infos := diag.Counts()
		body += DiagPanelHTML(diag.ToJSON(), errs, warns, infos)
	}

	return e.wrapPage(title, theme, body, page.GameMode, page.PageType), nil
}

// RenderPartial renders just the body content + title, no HTML shell.
// Used for client-side navigation where the shell is already loaded.
func (e *Engine) RenderPartial(data []byte) (string, string, error) {
	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return "", "", fmt.Errorf("parsing page: %w", err)
	}

	var bodyHTML strings.Builder
	for _, raw := range page.Body {
		html, err := e.renderAtom(raw)
		if err != nil {
			return "", "", err
		}
		bodyHTML.WriteString(html)
	}

	title := page.Title
	if title == "" {
		title = "ChefScript"
	}

	// Inject diagnostics panel so it survives partial navigation
	body := bodyHTML.String()
	body += DiagPanelHTML("[]", 0, 0, 0)

	return body, title, nil
}

// renderAtom takes a raw JSON value and renders it to HTML
func (e *Engine) renderAtom(raw json.RawMessage) (string, error) {
	return e.renderAtomDiag(raw, nil, "", 0)
}

// renderAtomDiag renders an atom with diagnostics collection.
func (e *Engine) renderAtomDiag(raw json.RawMessage, diag *DiagCollector, parent string, depth int) (string, error) {
	// Try string first — raw text child
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s, nil
	}

	// Must be an array — an atom [tag, props?, ...children]
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		return "", fmt.Errorf("atom must be a string or array, got: %s", string(raw))
	}

	if len(arr) == 0 {
		return "", fmt.Errorf("empty atom")
	}

	// Position 0: tag name
	var tag string
	if err := json.Unmarshal(arr[0], &tag); err != nil {
		return "", fmt.Errorf("atom tag must be a string: %w", err)
	}

	// Position 1+: props and/or children
	props := make(map[string]interface{})
	childStart := 1

	if len(arr) > 1 {
		// Try to parse position 1 as an object (props)
		var obj map[string]interface{}
		if err := json.Unmarshal(arr[1], &obj); err == nil {
			props = obj
			childStart = 2
		}
	}

	// Resolve schema props — direct binary → component pipe
	ResolveSchemaProps(props, diag)

	// Validate atom
	if diag != nil {
		diag.ValidateAtom(tag, props, e, parent, depth)
	}

	// Remaining positions are children
	var children []string
	for i := childStart; i < len(arr); i++ {
		child, err := e.renderAtomDiag(arr[i], diag, tag, depth+1)
		if err != nil {
			return "", fmt.Errorf("in <%s> child %d: %w", tag, i-childStart, err)
		}
		children = append(children, child)
	}

	childrenHTML := strings.Join(children, "\n")

	// Check registry first, then fall back to HTML passthrough
	if comp, ok := e.Registry[tag]; ok {
		return comp.Render(props, childrenHTML, e)
	}

	return e.htmlPassthrough(tag, props, childrenHTML), nil
}

// htmlPassthrough renders unknown tags as raw HTML elements
func (e *Engine) htmlPassthrough(tag string, props map[string]interface{}, children string) string {
	attrs := propsToAttrs(props)
	// Void elements can self-close; all others need open+close tags
	voidElements := map[string]bool{
		"area": true, "base": true, "br": true, "col": true, "embed": true,
		"hr": true, "img": true, "input": true, "link": true, "meta": true,
		"source": true, "track": true, "wbr": true,
	}
	if children == "" && voidElements[tag] {
		if attrs != "" {
			return fmt.Sprintf("<%s %s/>", tag, attrs)
		}
		return fmt.Sprintf("<%s/>", tag)
	}
	if attrs != "" {
		return fmt.Sprintf("<%s %s>%s</%s>", tag, attrs, children, tag)
	}
	return fmt.Sprintf("<%s>%s</%s>", tag, children, tag)
}

// RenderChildren is exported so components can render nested atoms
func (e *Engine) RenderChildren(raw json.RawMessage) (string, error) {
	return e.renderAtom(raw)
}

// propsToAttrs converts a props map to HTML attribute string
func propsToAttrs(props map[string]interface{}) string {
	if len(props) == 0 {
		return ""
	}
	var parts []string
	for k, v := range props {
		switch val := v.(type) {
		case string:
			parts = append(parts, fmt.Sprintf(`%s="%s"`, k, val))
		case float64:
			if val == float64(int(val)) {
				parts = append(parts, fmt.Sprintf(`%s="%d"`, k, int(val)))
			} else {
				parts = append(parts, fmt.Sprintf(`%s="%g"`, k, val))
			}
		case bool:
			if val {
				parts = append(parts, k)
			}
		}
	}
	return strings.Join(parts, " ")
}

// wrapPage assembles the final HTML document
func (e *Engine) wrapPage(title, theme, body string, gameMode bool, pageType string) string {
	css := GetThemeCSS(theme)
	bsTheme := ""
	if theme == "dark" {
		bsTheme = ` data-bs-theme="dark"`
	}
	pageTypeAttr := ""
	if pageType != "" {
		pageTypeAttr = fmt.Sprintf(` data-page-type="%s"`, pageType)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en"%s>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>%s</title>
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
  <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
  <style>
%s
  </style>
</head>
<body%s>
%s
%s
<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>`, bsTheme, title, css, pageTypeAttr, body, csRuntime)
}
