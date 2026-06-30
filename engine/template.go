package engine

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var templateVarRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// PageContext holds data available to template variable substitution
type PageContext struct {
	User     map[string]interface{}
	Data     map[string]interface{}
	Session  map[string]interface{}
	Flash    string
	Redirect string // if set, server redirects instead of rendering
}

// NewPageContext returns an empty PageContext with initialized maps
func NewPageContext() *PageContext {
	return &PageContext{
		User:    map[string]interface{}{},
		Data:    map[string]interface{}{},
		Session: map[string]interface{}{},
	}
}

// ApplyContext substitutes {{var.path}} tokens in the page JSON.
// Works on the parsed JSON tree so complex types (arrays, objects) are preserved
// when a prop value is entirely a single template variable.
func ApplyContext(raw []byte, ctx *PageContext) ([]byte, error) {
	return ApplyContextDiag(raw, ctx, nil)
}

// ApplyContextDiag substitutes template vars with diagnostics logging.
func ApplyContextDiag(raw []byte, ctx *PageContext, diag *DiagCollector) ([]byte, error) {
	var parsed interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return raw, nil
	}

	bound := walkAndReplace(parsed, ctx, diag)

	result, err := json.Marshal(bound)
	if err != nil {
		return raw, nil
	}
	return result, nil
}

// walkAndReplace recursively walks a parsed JSON value and resolves template vars.
func walkAndReplace(v interface{}, ctx *PageContext, diag *DiagCollector) interface{} {
	switch val := v.(type) {
	case string:
		return resolveTemplateString(val, ctx, diag)

	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, v2 := range val {
			out[k] = walkAndReplace(v2, ctx, diag)
		}
		return out

	case []interface{}:
		out := make([]interface{}, len(val))
		for i, v2 := range val {
			out[i] = walkAndReplace(v2, ctx, diag)
		}
		return out

	default:
		return v
	}
}

// resolveTemplateString handles string values that may contain template vars.
// If the entire string is a single {{path}}, the resolved value keeps its type
// (array, object, number, etc). If mixed with other text, string interpolation.
func resolveTemplateString(s string, ctx *PageContext, diag *DiagCollector) interface{} {
	trimmed := strings.TrimSpace(s)

	// Fast path: no template vars at all
	if !strings.Contains(trimmed, "{{") {
		return s
	}

	// Check if the entire string is a single template variable
	matches := templateVarRe.FindAllStringIndex(trimmed, -1)
	if len(matches) == 1 && matches[0][0] == 0 && matches[0][1] == len(trimmed) {
		// Entire value is one template var — preserve type
		path := strings.TrimSpace(trimmed[2 : len(trimmed)-2])
		resolved := resolveContextPath(path, ctx)
		if resolved != nil {
			if diag != nil {
				diag.LogTemplateResolved(path, fmt.Sprintf("%T", resolved))
			}
			return resolved
		}
		if diag != nil {
			diag.LogTemplateNil(path)
		}
		return ""
	}

	// Mixed content — string interpolation
	return templateVarRe.ReplaceAllStringFunc(s, func(match string) string {
		path := strings.TrimSpace(match[2 : len(match)-2])
		resolved := resolveContextPath(path, ctx)
		if resolved == nil {
			if diag != nil {
				diag.LogTemplateNil(path)
			}
			return ""
		}
		if diag != nil {
			diag.LogTemplateResolved(path, fmt.Sprintf("%T", resolved))
		}
		return fmt.Sprintf("%v", resolved)
	})
}

func resolveContextPath(path string, ctx *PageContext) interface{} {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) == 0 {
		return nil
	}

	var root map[string]interface{}
	switch strings.ToLower(parts[0]) {
	case "user":
		root = ctx.User
	case "data":
		root = ctx.Data
	case "session":
		root = ctx.Session
	case "flash":
		return ctx.Flash
	default:
		return nil
	}

	if root == nil {
		return nil
	}
	if len(parts) == 1 {
		return root
	}
	return resolveMapPath(root, parts[1])
}

func resolveMapPath(m map[string]interface{}, path string) interface{} {
	parts := strings.SplitN(path, ".", 2)
	val, ok := m[parts[0]]
	if !ok {
		return nil
	}
	if len(parts) == 1 {
		return val
	}
	if nested, ok := val.(map[string]interface{}); ok {
		return resolveMapPath(nested, parts[1])
	}
	return nil
}
