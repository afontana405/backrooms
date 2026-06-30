package engine

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ── Diagnostic Entry ────────────────────────────────────────────────────────

type DiagLevel int

const (
	DiagInfo  DiagLevel = 0
	DiagWarn  DiagLevel = 1
	DiagError DiagLevel = 2
)

func (l DiagLevel) String() string {
	switch l {
	case DiagInfo:
		return "info"
	case DiagWarn:
		return "warn"
	case DiagError:
		return "error"
	}
	return "info"
}

type DiagEntry struct {
	Level    DiagLevel `json:"level"`
	Category string    `json:"cat"`
	Message  string    `json:"msg"`
	Detail   string    `json:"detail,omitempty"`
	Time     string    `json:"time"`
}

// ── Collector ───────────────────────────────────────────────────────────────

type DiagCollector struct {
	Entries []DiagEntry
	mu      sync.Mutex
	start   time.Time
}

func NewDiagCollector() *DiagCollector {
	return &DiagCollector{
		start: time.Now(),
	}
}

func (d *DiagCollector) add(level DiagLevel, category, message, detail string) {
	if d == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	elapsed := time.Since(d.start).Milliseconds()
	d.Entries = append(d.Entries, DiagEntry{
		Level:    level,
		Category: category,
		Message:  message,
		Detail:   detail,
		Time:     fmt.Sprintf("%dms", elapsed),
	})

	// Feed into global flight recorder
	GlobalFlight.Record("server", category, level, message, detail, "")
}

func (d *DiagCollector) Info(category, message string) {
	d.add(DiagInfo, category, message, "")
}

func (d *DiagCollector) InfoDetail(category, message, detail string) {
	d.add(DiagInfo, category, message, detail)
}

func (d *DiagCollector) Warn(category, message string) {
	d.add(DiagWarn, category, message, "")
}

func (d *DiagCollector) WarnDetail(category, message, detail string) {
	d.add(DiagWarn, category, message, detail)
}

func (d *DiagCollector) Error(category, message string) {
	d.add(DiagError, category, message, "")
}

func (d *DiagCollector) ErrorDetail(category, message, detail string) {
	d.add(DiagError, category, message, detail)
}

// Counts returns error, warning, info counts.
func (d *DiagCollector) Counts() (int, int, int) {
	if d == nil {
		return 0, 0, 0
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	var errs, warns, infos int
	for _, e := range d.Entries {
		switch e.Level {
		case DiagError:
			errs++
		case DiagWarn:
			warns++
		case DiagInfo:
			infos++
		}
	}
	return errs, warns, infos
}

// ToJSON serializes entries for embedding in the page.
func (d *DiagCollector) ToJSON() string {
	if d == nil {
		return "[]"
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	b, err := json.Marshal(d.Entries)
	if err != nil {
		return "[]"
	}
	return string(b)
}

// ── Page JSON Validation ────────────────────────────────────────────────────

// ValidatePageJSON checks the raw page JSON for structural issues.
func (d *DiagCollector) ValidatePageJSON(raw []byte) {
	if d == nil {
		return
	}

	// Check valid JSON
	var parsed interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		d.Error("json", fmt.Sprintf("Invalid JSON: %v", err))
		return
	}
	d.Info("json", "Page JSON is valid")

	// Check page structure
	page, ok := parsed.(map[string]interface{})
	if !ok {
		d.Error("json", "Page must be a JSON object")
		return
	}

	if _, ok := page["body"]; !ok {
		d.Error("json", "Page missing required 'body' array")
		return
	}

	body, ok := page["body"].([]interface{})
	if !ok {
		d.Error("json", "'body' must be an array")
		return
	}

	d.Info("json", fmt.Sprintf("Page has %d top-level atoms", len(body)))
}

// ── Atom Validation ─────────────────────────────────────────────────────────

// ValidateAtom checks an atom during render.
func (d *DiagCollector) ValidateAtom(tag string, props map[string]interface{}, e *Engine, parent string, depth int) {
	if d == nil {
		return
	}

	// Check if tag exists in registry
	_, inRegistry := e.Registry[tag]
	isHTMLTag := isKnownHTMLTag(tag)

	if !inRegistry && !isHTMLTag {
		d.WarnDetail("render", fmt.Sprintf("Unknown tag: '%s'", tag),
			"Not in component registry or HTML tag list. Will render as raw HTML element.")
	}

	// Check structural rules
	if (tag == "modal" || tag == "drawer" || tag == "confirm") && depth > 1 {
		d.ErrorDetail("rule",
			fmt.Sprintf("'%s' must be at top level of body, not inside '%s'", tag, parent),
			"Modals, drawers, and confirms position themselves fixed — nesting breaks layout.")
	}

	if tag == "accordion-item" && parent != "accordion" {
		d.WarnDetail("rule",
			fmt.Sprintf("'accordion-item' should be a direct child of 'accordion', found inside '%s'", parent),
			"accordion-item relies on accordion parent for expand/collapse behavior.")
	}

	if tag == "breadcrumb-item" && parent != "breadcrumb" {
		d.WarnDetail("rule",
			"'breadcrumb-item' should be a direct child of 'breadcrumb'",
			fmt.Sprintf("Found inside '%s'", parent))
	}

	if tag == "tab" && parent != "tabs" {
		d.WarnDetail("rule",
			"'tab' should be a direct child of 'tabs'",
			fmt.Sprintf("Found inside '%s'", parent))
	}

	if tag == "stepper-step" && parent != "stepper" {
		d.WarnDetail("rule",
			"'stepper-step' should be a direct child of 'stepper'",
			fmt.Sprintf("Found inside '%s'", parent))
	}

	// Log successful render
	d.InfoDetail("render", fmt.Sprintf("Rendered: %s", tag),
		fmt.Sprintf("depth=%d parent=%s props=%d", depth, parent, len(props)))
}

// ── Template Var Validation ─────────────────────────────────────────────────

// LogTemplateNil logs when a template variable resolves to nil.
func (d *DiagCollector) LogTemplateNil(path string) {
	if d == nil {
		return
	}
	d.WarnDetail("template",
		fmt.Sprintf("{{%s}} resolved to empty", path),
		"This template variable has no value in the page context. Check page loader injects this data.")
}

// LogTemplateResolved logs a successful template resolution.
func (d *DiagCollector) LogTemplateResolved(path string, valueType string) {
	if d == nil {
		return
	}
	d.InfoDetail("template",
		fmt.Sprintf("{{%s}} → %s", path, valueType),
		"")
}

// ── Known HTML Tags ─────────────────────────────────────────────────────────

var knownHTMLTags = map[string]bool{
	"a": true, "abbr": true, "address": true, "area": true, "article": true,
	"aside": true, "audio": true, "b": true, "base": true, "bdi": true,
	"bdo": true, "blockquote": true, "body": true, "br": true, "button": true,
	"canvas": true, "caption": true, "cite": true, "code": true, "col": true,
	"colgroup": true, "data": true, "datalist": true, "dd": true, "del": true,
	"details": true, "dfn": true, "dialog": true, "div": true, "dl": true,
	"dt": true, "em": true, "embed": true, "fieldset": true, "figcaption": true,
	"figure": true, "footer": true, "form": true, "h1": true, "h2": true,
	"h3": true, "h4": true, "h5": true, "h6": true, "head": true,
	"header": true, "hgroup": true, "hr": true, "html": true, "i": true,
	"iframe": true, "img": true, "input": true, "ins": true, "kbd": true,
	"label": true, "legend": true, "li": true, "link": true, "main": true,
	"map": true, "mark": true, "menu": true, "meta": true, "meter": true,
	"nav": true, "noscript": true, "object": true, "ol": true, "optgroup": true,
	"option": true, "output": true, "p": true, "picture": true, "pre": true,
	"progress": true, "q": true, "rp": true, "rt": true, "ruby": true,
	"s": true, "samp": true, "script": true, "section": true, "select": true,
	"slot": true, "small": true, "source": true, "span": true, "strong": true,
	"style": true, "sub": true, "summary": true, "sup": true, "table": true,
	"tbody": true, "td": true, "template": true, "textarea": true, "tfoot": true,
	"th": true, "thead": true, "time": true, "title": true, "tr": true,
	"track": true, "u": true, "ul": true, "var": true, "video": true, "wbr": true,
}

func isKnownHTMLTag(tag string) bool {
	return knownHTMLTags[strings.ToLower(tag)]
}

// ── Action Validation ───────────────────────────────────────────────────────

// LogAction logs an API action call.
func (d *DiagCollector) LogAction(action string, found bool) {
	if d == nil {
		return
	}
	if found {
		d.Info("action", fmt.Sprintf("Action '%s' found", action))
	} else {
		d.Error("action", fmt.Sprintf("Action '%s' not found — no handler registered", action))
	}
}
