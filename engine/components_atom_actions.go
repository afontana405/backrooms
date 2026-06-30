package engine

import (
	"fmt"
	"strings"
)

// ── Icon ────────────────────────────────────────────────────────────────
// ["icon", { "name": "search", "size": 20, "color": "currentColor" }]

func renderIcon(props map[string]interface{}, children string, e *Engine) (string, error) {
	name := propStr(props, "name", "")
	if name == "" {
		return "", fmt.Errorf("icon: 'name' prop is required")
	}

	size := int(propFloat(props, "size", 20))
	color := propStr(props, "color", "currentColor")
	class := propStr(props, "class", "")
	dataID := propStr(props, "data-id", fmt.Sprintf("icon--%s", name))

	cls := fmt.Sprintf("bi bi-%s", name)
	if class != "" {
		cls += " " + class
	}

	style := fmt.Sprintf("font-size:%dpx;color:%s", size, color)

	return fmt.Sprintf(`<i class="%s" data-id="%s" style="%s"></i>`,
		cls, dataID, style), nil
}

// ── Button ──────────────────────────────────────────────────────────────
// ["button", { "label": "Submit", "variant": "solid", "on:click": "my/action" }]

func renderButton(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	action := propStr(props, "on:click", "")
	rawOnclick := propStr(props, "onclick", "")
	variant := propStr(props, "variant", "solid")
	color := propStr(props, "color", "")
	size := propStr(props, "size", "md")

	cls := fmt.Sprintf("cs-button cs-button--%s cs-button--%s", variant, size)
	if color != "" {
		cls += fmt.Sprintf(" cs-button--color-%s", color)
	}
	props["class"] = cls

	onclick := ""
	typeAttr := ""
	if rawOnclick != "" {
		onclick = fmt.Sprintf(` onclick="%s"`, rawOnclick)
		typeAttr = ` type="button"`
	} else if action != "" {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, action)
		typeAttr = ` type="button"`
	}

	disabledAttr := ""
	if d, ok := props["disabled"]; ok && d == true {
		disabledAttr = " disabled"
	}

	return fmt.Sprintf(`<button%s%s%s%s>%s</button>`, typeAttr, userAttrs(props, ""), onclick, disabledAttr, label), nil
}

// ── IconButton ──────────────────────────────────────────────────────────
// ["icon-button", { "icon": "trash", "aria-label": "Delete", "on:click": "items/delete" }]

func renderIconButton(props map[string]interface{}, children string, e *Engine) (string, error) {
	icon := propStr(props, "icon", "")
	ariaLabel := propStr(props, "aria-label", icon)
	size := propSize(props, "md")
	variant := propVariant(props, "ghost")
	action := propStr(props, "on:click", "")
	color := propStr(props, "color", "")
	dataID := propStr(props, "data-id", "icon-button")

	iconHTML := ""
	if icon != "" {
		iconHTML, _ = renderIcon(map[string]interface{}{"name": icon}, "", e)
	} else {
		iconHTML = children
	}

	cls := fmt.Sprintf("cs-icon-button cs-icon-button--%s cs-icon-button--%s", variant, size)
	if color != "" {
		cls += fmt.Sprintf(" cs-icon-button--color-%s", color)
	}

	onclick := ""
	if action != "" {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, action)
	}

	return fmt.Sprintf(`<button class="%s" aria-label="%s" type="button" data-id="%s"%s>%s</button>`,
		cls, ariaLabel, dataID, onclick, iconHTML), nil
}

// ── ButtonGroup ─────────────────────────────────────────────────────────
// ["button-group", {}, ["button", ...], ["button", ...]]

func renderButtonGroup(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "button-group")
	return fmt.Sprintf(`<div class="cs-button-group" data-id="%s">%s</div>`, dataID, children), nil
}

// ── CopyButton ──────────────────────────────────────────────────────────
// ["copy-button", { "value": "text-to-copy", "label": "Copy key" }]

func renderCopyButton(props map[string]interface{}, children string, e *Engine) (string, error) {
	value := propStr(props, "value", "")
	label := propStr(props, "label", "Copy")
	size := propSize(props, "md")
	variant := propVariant(props, "outline")
	dataID := propStr(props, "data-id", "copy-button")

	cls := fmt.Sprintf("cs-button cs-button--%s cs-button--%s cs-copy-button", variant, size)
	escaped := strings.ReplaceAll(value, `"`, `&quot;`)
	escaped = strings.ReplaceAll(escaped, `'`, `&#39;`)

	onclick := fmt.Sprintf(
		`navigator.clipboard.writeText('%s').then(function(){var b=document.querySelector('[data-id="%s"]');var orig=b.textContent;b.textContent='✓ Copied';setTimeout(function(){b.textContent=orig},2000)})`,
		escaped, dataID)

	return fmt.Sprintf(`<button class="%s" type="button" data-id="%s" onclick="%s">%s</button>`,
		cls, dataID, onclick, label), nil
}
