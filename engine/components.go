package engine

import "fmt"

// Component is the interface every registered component implements
type Component interface {
	Render(props map[string]interface{}, children string, e *Engine) (string, error)
}

// ComponentFunc is a convenience type so you can register plain functions
type ComponentFunc func(props map[string]interface{}, children string, e *Engine) (string, error)

func (f ComponentFunc) Render(props map[string]interface{}, children string, e *Engine) (string, error) {
	return f(props, children, e)
}

// RegisterDefaults loads the starter component library
func RegisterDefaults(e *Engine) {
	// ── Atom: Layout ────────────────────────────────────────────────
	e.Register("header", ComponentFunc(renderHeader))
	e.Register("footer", ComponentFunc(renderFooter))
	e.Register("card", ComponentFunc(renderCard))
	e.Register("sidebar", ComponentFunc(renderSidebar))
	e.Register("section", ComponentFunc(renderSection))
	e.Register("divider", ComponentFunc(renderDivider))
	e.Register("split-view", ComponentFunc(renderSplitView))
	e.Register("split-pane", ComponentFunc(renderSplitPane))
	e.Register("paper", ComponentFunc(renderPaper))
	e.Register("app-bar", ComponentFunc(renderAppBar))
	e.Register("box", ComponentFunc(renderBox))
	e.Register("grid", ComponentFunc(renderGrid))
	e.Register("stack", ComponentFunc(renderStack))
	e.Register("image-list", ComponentFunc(renderImageList))
	e.Register("image-list-item", ComponentFunc(renderImageListItem))
	e.Register("masonry", ComponentFunc(renderMasonry))

	// ── Atom: Navigation ────────────────────────────────────────────
	e.Register("nav", ComponentFunc(renderNav))
	e.Register("nav-link", ComponentFunc(renderNavLink))
	e.Register("tabs", ComponentFunc(renderTabs))
	e.Register("tab", ComponentFunc(renderTab))
	e.Register("accordion", ComponentFunc(renderAccordion))
	e.Register("accordion-item", ComponentFunc(renderAccordionItem))
	e.Register("breadcrumb", ComponentFunc(renderBreadcrumb))
	e.Register("breadcrumb-item", ComponentFunc(renderBreadcrumbItem))
	e.Register("pagination", ComponentFunc(renderPagination))
	e.Register("stepper", ComponentFunc(renderStepper))
	e.Register("stepper-step", ComponentFunc(renderStepperStep))
	e.Register("toolbar", ComponentFunc(renderToolbar))
	e.Register("bottom-nav", ComponentFunc(renderBottomNav))
	e.Register("bottom-nav-action", ComponentFunc(renderBottomNavAction))
	e.Register("menubar", ComponentFunc(renderMenubar))
	e.Register("menubar-item", ComponentFunc(renderMenubarItem))
	e.Register("speed-dial", ComponentFunc(renderSpeedDial))
	e.Register("speed-dial-action", ComponentFunc(renderSpeedDialAction))

	// ── Atom: Inputs ────────────────────────────────────────────────
	e.Register("input", ComponentFunc(renderInput))
	e.Register("textarea", ComponentFunc(renderTextarea))
	e.Register("select", ComponentFunc(renderSelect))
	e.Register("native-select", ComponentFunc(renderNativeSelect))
	e.Register("autocomplete", ComponentFunc(renderAutocomplete))
	e.Register("form-field", ComponentFunc(renderFormField))
	e.Register("multi-select", ComponentFunc(renderMultiSelect))
	e.Register("checkbox", ComponentFunc(renderCheckbox))
	e.Register("radio", ComponentFunc(renderRadio))
	e.Register("switch", ComponentFunc(renderSwitch))
	e.Register("form", ComponentFunc(renderForm))
	e.Register("slider", ComponentFunc(renderSlider))
	e.Register("number-input", ComponentFunc(renderNumberInput))
	e.Register("file-upload", ComponentFunc(renderFileUpload))
	e.Register("tag-input", ComponentFunc(renderTagInput))
	e.Register("date-input", ComponentFunc(renderDateInput))
	e.Register("search", ComponentFunc(renderSearch))
	e.Register("color-input", ComponentFunc(renderColorInput))
	e.Register("fab", ComponentFunc(renderFab))
	e.Register("toggle-group", ComponentFunc(renderToggleGroup))
	e.Register("toggle-button", ComponentFunc(renderToggleButton))
	e.Register("transfer-list", ComponentFunc(renderTransferList))

	// ── Atom: Actions ───────────────────────────────────────────────
	e.Register("button", ComponentFunc(renderButton))
	e.Register("icon", ComponentFunc(renderIcon))
	e.Register("icon-button", ComponentFunc(renderIconButton))
	e.Register("button-group", ComponentFunc(renderButtonGroup))
	e.Register("copy-button", ComponentFunc(renderCopyButton))

	// ── Atom: Display ───────────────────────────────────────────────
	e.Register("heading", ComponentFunc(renderHeading))
	e.Register("text", ComponentFunc(renderText))
	e.Register("avatar", ComponentFunc(renderAvatar))
	e.Register("avatar-group", ComponentFunc(renderAvatarGroup))
	e.Register("empty-state", ComponentFunc(renderEmptyState))
	e.Register("kbd", ComponentFunc(renderKbd))
	e.Register("code", ComponentFunc(renderCode))
	e.Register("code-block", ComponentFunc(renderCodeBlock))
	e.Register("timeline", ComponentFunc(renderTimeline))
	e.Register("timeline-item", ComponentFunc(renderTimelineItem))
	e.Register("rating", ComponentFunc(renderRating))
	e.Register("callout", ComponentFunc(renderCallout))
	e.Register("image", ComponentFunc(renderImage))
	e.Register("link", ComponentFunc(renderLink))
	e.Register("tag", ComponentFunc(renderTag))

	// ── Atom: Data ──────────────────────────────────────────────────
	e.Register("stat-card", ComponentFunc(renderStatCard))
	e.Register("table", ComponentFunc(renderTable))
	e.Register("list", ComponentFunc(renderList))
	e.Register("list-item", ComponentFunc(renderListItem))
	e.Register("kv-list", ComponentFunc(renderKvList))
	e.Register("kv-item", ComponentFunc(renderKvItem))
	e.Register("data-grid", ComponentFunc(renderDataGrid))
	e.Register("tree", ComponentFunc(renderTree))
	e.Register("tree-item", ComponentFunc(renderTreeItem))
	e.Register("virtual-list", ComponentFunc(renderVirtualList))
	e.Register("chart", ComponentFunc(renderChart))
	e.Register("calendar", ComponentFunc(renderCalendar))

	// ── Atom: Feedback ──────────────────────────────────────────────
	e.Register("alert", ComponentFunc(renderAlert))
	e.Register("badge", ComponentFunc(renderBadge))
	e.Register("chip", ComponentFunc(renderChip))
	e.Register("spinner", ComponentFunc(renderSpinner))
	e.Register("skeleton", ComponentFunc(renderSkeleton))
	e.Register("progress", ComponentFunc(renderProgress))
	e.Register("tooltip", ComponentFunc(renderTooltip))
	e.Register("banner", ComponentFunc(renderBanner))

	// ── Atom: Overlay ───────────────────────────────────────────────
	e.Register("menu", ComponentFunc(renderMenu))
	e.Register("menu-item", ComponentFunc(renderMenuItem))
	e.Register("popover", ComponentFunc(renderPopover))
	e.Register("drawer", ComponentFunc(renderDrawer))
	e.Register("snackbar", ComponentFunc(renderSnackbar))
	e.Register("confirm", ComponentFunc(renderConfirm))
	e.Register("notification", ComponentFunc(renderNotification))
	e.Register("notification-item", ComponentFunc(renderNotificationItem))
	e.Register("command", ComponentFunc(renderCommand))
	e.Register("command-item", ComponentFunc(renderCommandItem))
	e.Register("context-menu", ComponentFunc(renderContextMenu))
	e.Register("hover-card", ComponentFunc(renderHoverCard))
	e.Register("backdrop", ComponentFunc(renderBackdrop))
	e.Register("dialog", ComponentFunc(renderDialog))
	e.Register("dialog-title", ComponentFunc(renderDialogTitle))
	e.Register("dialog-content", ComponentFunc(renderDialogContent))
	e.Register("dialog-actions", ComponentFunc(renderDialogActions))

	// ── Atom: Media ─────────────────────────────────────────────────
	e.Register("video", ComponentFunc(renderVideo))
	e.Register("audio", ComponentFunc(renderAudio))
	e.Register("iframe", ComponentFunc(renderIframe))
	e.Register("aspect-ratio", ComponentFunc(renderAspectRatio))
	e.Register("carousel", ComponentFunc(renderCarousel))
	e.Register("rich-text", ComponentFunc(renderRichText))

	// ── Atom: Chat ──────────────────────────────────────────────────
	e.Register("chat-widget", ComponentFunc(renderChatWidget))
	e.Register("data-chat", ComponentFunc(renderDataChat))
}

// --- Prop helpers ---

// userAttrs extracts style, id, data-id, class (appended) from props and returns an HTML attr string.
// baseClass is the component's own class — user "class" prop is appended to it.
func userAttrs(props map[string]interface{}, baseClass string) string {
	cls := baseClass
	if extra := propStr(props, "class", ""); extra != "" {
		cls += " " + extra
	}
	out := fmt.Sprintf(` class="%s"`, cls)
	if id := propStr(props, "id", ""); id != "" {
		out += fmt.Sprintf(` id="%s"`, id)
	}
	if did := propStr(props, "data-id", ""); did != "" {
		out += fmt.Sprintf(` data-id="%s"`, did)
	}
	if style := propStr(props, "style", ""); style != "" {
		out += fmt.Sprintf(` style="%s"`, style)
	}
	return out
}

func propStr(props map[string]interface{}, key, fallback string) string {
	if v, ok := props[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return fallback
}

// propAnyStr converts any prop value to string — handles int, float, bool, string.
// Returns "" if key is missing. Never silently drops non-string values.
func propAnyStr(props map[string]interface{}, key string) string {
	v, ok := props[key]
	if !ok || v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case int:
		return fmt.Sprintf("%d", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func propFloat(props map[string]interface{}, key string, fallback float64) float64 {
	if v, ok := props[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return fallback
}

func propBool(props map[string]interface{}, key string, fallback bool) bool {
	if v, ok := props[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return fallback
}

// ── Helpers for MUI atoms ───────────────────────────────────────────────

// esc escapes a string for safe HTML attribute output.
func esc(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			out = append(out, []byte("&amp;")...)
		case '<':
			out = append(out, []byte("&lt;")...)
		case '>':
			out = append(out, []byte("&gt;")...)
		case '"':
			out = append(out, []byte("&quot;")...)
		default:
			out = append(out, s[i])
		}
	}
	return string(out)
}

// toString converts any value to string.
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// propInt converts any numeric value to int.
func propInt(props map[string]interface{}, key string, fallback int) int {
	if v, ok := props[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		case int64:
			return int(n)
		}
	}
	return fallback
}

// toInterfaceSlice converts a value to []interface{}.
func toInterfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return nil
}

// bsIcon returns an inline Bootstrap Icon HTML element.
func bsIcon(name string, size int) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("<i class=\"bi bi-%s\" style=\"font-size:%dpx\"></i>", name, size)
}

func propVariant(props map[string]interface{}, fallback string) string {
	v := propStr(props, "variant", fallback)
	if buttonVariants[v] {
		return v
	}
	return fallback
}

func propSize(props map[string]interface{}, fallback string) string {
	v := propStr(props, "size", fallback)
	if validSizes[v] {
		return v
	}
	return fallback
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
