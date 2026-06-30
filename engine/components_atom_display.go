package engine

import (
	"fmt"
	"strings"
)

// ── Heading ─────────────────────────────────────────────────────────────
// ["heading", { "level": 2 }, "Page Title"]

func renderHeading(props map[string]interface{}, children string, e *Engine) (string, error) {
	level := int(propFloat(props, "level", 2))
	if level < 1 || level > 6 {
		level = 2
	}
	return fmt.Sprintf(`<h%d%s>%s</h%d>`, level, userAttrs(props, "cs-heading"), children, level), nil
}

// ── Text ────────────────────────────────────────────────────────────────
// ["text", {}, "Paragraph content"]

func renderText(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<p%s>%s</p>`, userAttrs(props, "cs-text"), children), nil
}

// ── Avatar ──────────────────────────────────────────────────────────────

func renderAvatar(props map[string]interface{}, children string, e *Engine) (string, error) {
	name := propStr(props, "name", "")
	src := propStr(props, "src", "")
	size := propStr(props, "size", "md")

	if src != "" {
		return fmt.Sprintf(`<div class="cs-avatar cs-avatar--%s"><img class="cs-avatar__img" src="%s" alt="%s"></div>`, size, src, name), nil
	}

	// Generate initials
	initials := ""
	if name != "" {
		parts := strings.Fields(name)
		for _, p := range parts {
			if len(p) > 0 {
				initials += string(p[0])
			}
		}
		if len(initials) > 2 {
			initials = initials[:2]
		}
	}
	if initials == "" {
		initials = `<i class="bi bi-person"></i>`
	}

	return fmt.Sprintf(`<div class="cs-avatar cs-avatar--%s">%s</div>`, size, initials), nil
}

// ── AvatarGroup ─────────────────────────────────────────────────────────

func renderAvatarGroup(props map[string]interface{}, children string, e *Engine) (string, error) {
	max := int(propFloat(props, "max", 4))
	size := propStr(props, "size", "md")
	dataID := propStr(props, "data-id", "avatar-group")
	_ = max
	_ = size
	return fmt.Sprintf(`<div class="cs-avatar-group" data-id="%s">%s</div>`, dataID, children), nil
}

// ── EmptyState ──────────────────────────────────────────────────────────

func renderEmptyState(props map[string]interface{}, children string, e *Engine) (string, error) {
	icon := propStr(props, "icon", "inbox")
	title := propStr(props, "title", "Nothing here yet")
	description := propStr(props, "description", "")
	action := propStr(props, "action", "")
	onClick := propStr(props, "on:click", "")
	dataID := propStr(props, "data-id", "empty-state")

	iconHTML, _ := renderIcon(map[string]interface{}{"name": icon, "size": float64(40)}, "", e)

	descHTML := ""
	if description != "" {
		descHTML = fmt.Sprintf(`<p class="cs-empty-state__desc">%s</p>`, description)
	}

	btnHTML := ""
	if action != "" {
		onclick := ""
		if onClick != "" {
			onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, onClick)
		}
		btnHTML = fmt.Sprintf(`<button class="cs-button cs-button--outline cs-button--md" data-id="%s--action"%s>%s</button>`,
			dataID, onclick, action)
	}

	return fmt.Sprintf(`<div class="cs-empty-state" data-id="%s">
  <div class="cs-empty-state__icon">%s</div>
  <h3 class="cs-empty-state__title">%s</h3>
  %s
  %s
</div>`, dataID, iconHTML, title, descHTML, btnHTML), nil
}

// ── Kbd ─────────────────────────────────────────────────────────────────

func renderKbd(props map[string]interface{}, children string, e *Engine) (string, error) {
	keys := propStr(props, "keys", children)
	dataID := propStr(props, "data-id", "kbd")

	parts := strings.Split(keys, "+")
	html := ""
	for i, k := range parts {
		if i > 0 {
			html += `<span class="cs-kbd__sep">+</span>`
		}
		html += fmt.Sprintf(`<kbd class="cs-kbd">%s</kbd>`, strings.TrimSpace(k))
	}
	return fmt.Sprintf(`<span class="cs-kbd-group" data-id="%s">%s</span>`, dataID, html), nil
}

// ── Code ────────────────────────────────────────────────────────────────

func renderCode(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<code class="cs-code">%s</code>`, children), nil
}

// ── CodeBlock ───────────────────────────────────────────────────────────

func renderCodeBlock(props map[string]interface{}, children string, e *Engine) (string, error) {
	lang := propStr(props, "lang", "")
	content := propStr(props, "content", children)
	dataID := propStr(props, "data-id", "code-block")

	langBadge := ""
	if lang != "" {
		langBadge = fmt.Sprintf(`<span class="cs-code-block__lang">%s</span>`, lang)
	}

	id := fmt.Sprintf("cb-%s", dataID)
	return fmt.Sprintf(`<div class="cs-code-block" data-id="%s">
  <div class="cs-code-block__header">
    %s
    <button class="cs-code-block__copy" onclick="csCopyCode('%s')" data-id="%s--copy">Copy</button>
  </div>
  <pre id="%s" class="cs-code-block__pre"><code>%s</code></pre>
</div>`, dataID, langBadge, id, dataID, id, content), nil
}

// ── Timeline ────────────────────────────────────────────────────────────

func renderTimeline(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "timeline")
	return fmt.Sprintf(`<div class="cs-timeline" data-id="%s">%s</div>`, dataID, children), nil
}

func renderTimelineItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	time := propStr(props, "time", "")
	title := propStr(props, "title", "")
	description := propStr(props, "description", "")
	color := propStr(props, "color", "default")
	dataID := propStr(props, "data-id", "timeline-item")

	timeHTML := ""
	if time != "" {
		timeHTML = fmt.Sprintf(`<div class="cs-timeline-item__time">%s</div>`, time)
	}

	descHTML := ""
	if description != "" {
		descHTML = fmt.Sprintf(`<div class="cs-timeline-item__desc">%s</div>`, description)
	}

	return fmt.Sprintf(`<div class="cs-timeline-item" data-id="%s">
  <div class="cs-timeline-item__track">
    <div class="cs-timeline-item__dot cs-timeline-item__dot--%s"></div>
    <div class="cs-timeline-item__line"></div>
  </div>
  <div class="cs-timeline-item__content">
    %s
    <div class="cs-timeline-item__title">%s</div>
    %s
    %s
  </div>
</div>`, dataID, color, timeHTML, title, descHTML, children), nil
}

// ── Callout ─────────────────────────────────────────────────────────────
// ["callout", { "variant": "warning", "title": "Heads up" }, "Message"]

func renderCallout(props map[string]interface{}, children string, e *Engine) (string, error) {
	variant := propStr(props, "variant", "info")
	if variant == "error" {
		variant = "danger"
	}

	return fmt.Sprintf(`<div%s><div class="cs-callout__body">%s</div></div>`, userAttrs(props, "cs-callout cs-callout--"+variant), children), nil
}

// ── Image ───────────────────────────────────────────────────────────────
// ["image", { "src": "/public/logo.png", "alt": "Logo" }]

func renderImage(props map[string]interface{}, children string, e *Engine) (string, error) {
	src := propStr(props, "src", "")
	alt := propStr(props, "alt", "")
	width := propStr(props, "width", "")
	height := propStr(props, "height", "")
	rounded := propBool(props, "rounded", false)
	dataID := propStr(props, "data-id", "image")

	cls := "cs-image"
	if rounded {
		cls += " cs-image--rounded"
	}

	sizeAttrs := ""
	if width != "" {
		sizeAttrs += fmt.Sprintf(` width="%s"`, width)
	}
	if height != "" {
		sizeAttrs += fmt.Sprintf(` height="%s"`, height)
	}

	return fmt.Sprintf(`<img class="%s" src="%s" alt="%s"%s data-id="%s" />`,
		cls, src, alt, sizeAttrs, dataID), nil
}

// ── Link ────────────────────────────────────────────────────────────────
// ["link", { "href": "/page/dashboard", "label": "Go to dashboard" }]

func renderLink(props map[string]interface{}, children string, e *Engine) (string, error) {
	href := propStr(props, "href", "#")
	target := propStr(props, "target", "")
	label := propStr(props, "label", children)
	variant := propStr(props, "variant", "default")
	dataID := propStr(props, "data-id", "link")

	targetAttr := ""
	if target != "" {
		targetAttr = fmt.Sprintf(` target="%s"`, target)
	}

	return fmt.Sprintf(`<a class="cs-link cs-link--%s" href="%s"%s data-id="%s">%s</a>`,
		variant, href, targetAttr, dataID, label), nil
}

// ── Rating ──────────────────────────────────────────────────────────────

func renderRating(props map[string]interface{}, children string, e *Engine) (string, error) {
	value := int(propFloat(props, "value", 0))
	max := int(propFloat(props, "max", 5))
	readonly := propBool(props, "readonly", false)
	name := propStr(props, "name", "rating")
	dataID := propStr(props, "data-id", "rating")

	stars := ""
	for i := 1; i <= max; i++ {
		cls := "cs-rating__star"
		if i <= value {
			cls += " cs-rating__star--filled"
		}
		if readonly {
			stars += fmt.Sprintf(`<span class="%s" data-id="%s--star-%d">★</span>`, cls, dataID, i)
		} else {
			stars += fmt.Sprintf(`<button class="%s" data-rating-value="%d" data-id="%s--star-%d" type="button">★</button>`,
				cls, i, dataID, i)
		}
	}

	readonlyAttr := ""
	if readonly {
		readonlyAttr = ` data-rating-readonly="true"`
	}

	return fmt.Sprintf(`<div class="cs-rating" data-id="%s" data-rating="%d" data-rating-name="%s"%s>
  %s
  <input type="hidden" name="%s" value="%d" data-rating-input />
</div>`, dataID, value, name, readonlyAttr, stars, name, value), nil
}

// ── Tag ─────────────────────────────────────────────────────────────────
// ["tag", { "label": "Go", "color": "cyan" }]

func renderTag(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	variant := propStr(props, "variant", "secondary")
	bsVariant := variant
	if variant == "error" {
		bsVariant = "danger"
	}
	return fmt.Sprintf(`<span%s>%s</span>`, userAttrs(props, "badge text-bg-"+bsVariant), label), nil
}
