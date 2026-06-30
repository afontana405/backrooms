package engine

import (
	"fmt"
	"strings"
)

// ── Alert ───────────────────────────────────────────────────────────────
// ["alert", { "variant": "success", "title": "Saved" }, "Your changes were saved."]

func renderAlert(props map[string]interface{}, children string, e *Engine) (string, error) {
	variant := propStr(props, "variant", "info")
	title := propStr(props, "title", "")
	dismissible := propBool(props, "dismissible", false)
	if variant == "danger" {
		variant = "error"
	}

	titleHTML := ""
	if title != "" {
		titleHTML = fmt.Sprintf(`<span class="cs-alert__title">%s</span>`, title)
	}

	dismissHTML := ""
	if dismissible {
		dismissHTML = `<button class="cs-alert__dismiss" onclick="this.closest('.cs-alert').remove()" aria-label="Close">&#10005;</button>`
	}

	return fmt.Sprintf(`<div%s role="alert"><div class="cs-alert__body">%s%s</div>%s</div>`,
		userAttrs(props, "cs-alert cs-alert--"+variant), titleHTML, children, dismissHTML), nil
}

// ── Badge ───────────────────────────────────────────────────────────────
// ["badge", { "variant": "success" }, "Active"]

func renderBadge(props map[string]interface{}, children string, e *Engine) (string, error) {
	variant := propStr(props, "variant", "default")
	dataID := propStr(props, "data-id", "badge")

	return fmt.Sprintf(`<span class="cs-badge cs-badge--%s" data-id="%s">%s</span>`,
		variant, dataID, children), nil
}

// ── Chip ────────────────────────────────────────────────────────────────
// ["chip", { "label": "React", "dismissible": true }]

func renderChip(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	dismissible := propBool(props, "dismissible", false)
	color := propStr(props, "color", "")
	dataID := propStr(props, "data-id", "chip")

	cls := "cs-chip"
	if color != "" {
		cls += " cs-chip--" + color
	}

	dismissHTML := ""
	if dismissible {
		dismissHTML = `<button class="cs-chip__dismiss" onclick="this.closest('.cs-chip').remove()" aria-label="Remove">&#10005;</button>`
	}

	return fmt.Sprintf(`<span class="%s" data-id="%s">%s%s</span>`, cls, dataID, label, dismissHTML), nil
}

// ── Spinner ─────────────────────────────────────────────────────────────
// ["spinner", { "size": "md" }]

func renderSpinner(props map[string]interface{}, children string, e *Engine) (string, error) {
	size := propSize(props, "md")
	dataID := propStr(props, "data-id", "spinner")

	return fmt.Sprintf(`<span class="cs-spinner cs-spinner--%s" role="status" aria-label="Loading" data-id="%s"></span>`,
		size, dataID), nil
}

// ── Skeleton ────────────────────────────────────────────────────────────
// ["skeleton", { "lines": 3, "avatar": true }]

func renderSkeleton(props map[string]interface{}, children string, e *Engine) (string, error) {
	lines := int(propFloat(props, "lines", 1))
	avatar := propBool(props, "avatar", false)
	dataID := propStr(props, "data-id", "skeleton")

	var html strings.Builder
	html.WriteString(fmt.Sprintf(`<div class="cs-skeleton" data-id="%s">`, dataID))

	if avatar {
		html.WriteString(`<div class="cs-skeleton__avatar cs-skeleton__block"></div>`)
	}

	html.WriteString(`<div class="cs-skeleton__lines">`)
	for i := 0; i < lines; i++ {
		width := "100%"
		if i == lines-1 && lines > 1 {
			width = "60%"
		}
		html.WriteString(fmt.Sprintf(`<div class="cs-skeleton__block" style="width:%s"></div>`, width))
	}
	html.WriteString(`</div></div>`)

	return html.String(), nil
}

// ── Progress ────────────────────────────────────────────────────────────
// ["progress", { "value": 65, "max": 100, "label": "Upload" }]

func renderProgress(props map[string]interface{}, children string, e *Engine) (string, error) {
	value := propFloat(props, "value", 0)
	max := propFloat(props, "max", 100)
	label := propStr(props, "label", "")
	color := propStr(props, "color", "")
	dataID := propStr(props, "data-id", "progress")

	pct := 0.0
	if max > 0 {
		pct = (value / max) * 100
	}
	if pct > 100 {
		pct = 100
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<div class="cs-progress__label"><span>%s</span><span>%.0f%%</span></div>`, label, pct)
	}

	fillStyle := fmt.Sprintf("width:%.1f%%", pct)
	if color != "" {
		fillStyle += fmt.Sprintf(";background:%s", color)
	}

	return fmt.Sprintf(`<div class="cs-progress" data-id="%s">
  %s
  <div class="cs-progress__track">
    <div class="cs-progress__fill" style="%s"></div>
  </div>
</div>`, dataID, labelHTML, fillStyle), nil
}

// ── Banner ──────────────────────────────────────────────────────────────
// ["banner", { "variant": "warning", "title": "Maintenance" }, "Message"]

func renderBanner(props map[string]interface{}, children string, e *Engine) (string, error) {
	variant := propStr(props, "variant", "info")
	title := propStr(props, "title", "")
	dismissible := propBool(props, "dismissible", true)
	dataID := propStr(props, "data-id", "banner--"+variant)

	icons := map[string]string{
		"info":    "&#9432;",
		"success": "&#10003;",
		"warning": "&#9888;",
		"danger":  "&#128683;",
	}
	icon := icons[variant]
	if icon == "" {
		icon = icons["info"]
	}

	titleHTML := ""
	if title != "" {
		titleHTML = fmt.Sprintf(`<strong class="cs-banner__title">%s</strong> `, title)
	}

	dismissHTML := ""
	if dismissible {
		dismissHTML = `<button class="cs-banner__dismiss" onclick="this.closest('.cs-banner').remove()" aria-label="Dismiss">&#10005;</button>`
	}

	return fmt.Sprintf(`<div class="cs-banner cs-banner--%s" role="alert" data-id="%s">
  <span class="cs-banner__icon">%s</span>
  <div class="cs-banner__body">%s%s</div>
  %s
</div>`, variant, dataID, icon, titleHTML, children, dismissHTML), nil
}

// ── Tooltip ─────────────────────────────────────────────────────────────
// ["tooltip", { "text": "More info" }, ["button", {}, "Hover me"]]

func renderTooltip(props map[string]interface{}, children string, e *Engine) (string, error) {
	text := propStr(props, "text", "")
	position := propStr(props, "position", "top")
	dataID := propStr(props, "data-id", "tooltip")

	return fmt.Sprintf(`<span class="cs-tooltip cs-tooltip--%s" data-id="%s">
  %s
  <span class="cs-tooltip__text">%s</span>
</span>`, position, dataID, children, text), nil
}
