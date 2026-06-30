package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// ── Nav ─────────────────────────────────────────────────────────────────
// ["nav", { "brand": "MyApp" }, ["nav-link", ...]]

func renderNav(props map[string]interface{}, children string, e *Engine) (string, error) {
	brand := propStr(props, "brand", "")
	brandHTML := ""
	if brand != "" {
		brandHTML = fmt.Sprintf(`<span class="nav-brand">%s</span>`, brand)
	}
	return fmt.Sprintf(`<nav class="cs-nav">%s<div class="nav-links">%s</div></nav>`,
		brandHTML, children), nil
}

func renderNavLink(props map[string]interface{}, children string, e *Engine) (string, error) {
	href := propStr(props, "href", "#")
	return fmt.Sprintf(`<a class="cs-nav-link" href="%s">%s</a>`, href, children), nil
}

// ── Tabs ────────────────────────────────────────────────────────────────
// ["tabs", { "id": "t1" }, ["tab", ...], ["tab", ...]]

func renderTabs(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "tabs")
	return fmt.Sprintf(`<div class="cs-tabs" data-id="%s">%s</div>`, dataID, children), nil
}

func renderTab(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "Tab")
	panel := propStr(props, "panel", label)
	active := propBool(props, "active", false)
	dataID := propStr(props, "data-id", "tab--"+panel)

	triggerClass := "cs-tab"
	if active {
		triggerClass += " cs-tab--active"
	}

	panelClass := "cs-tab-panel"
	if active {
		panelClass += " cs-tab-panel--active"
	}

	return fmt.Sprintf(`<button class="%s" data-tab-trigger="%s" data-id="%s">%s</button>
<div class="%s" data-tab-panel="%s">%s</div>`,
		triggerClass, panel, dataID, label,
		panelClass, panel, children), nil
}

// ── Accordion ───────────────────────────────────────────────────────────

func renderAccordion(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "accordion")
	return fmt.Sprintf(`<div class="cs-accordion" data-id="%s">%s</div>`, dataID, children), nil
}

func renderAccordionItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	title := propStr(props, "title", "")
	open := propBool(props, "open", false)
	dataID := propStr(props, "data-id", "accordion-item")

	cls := "cs-accordion-item"
	bodyStyle := "max-height:0;overflow:hidden;"
	if open {
		cls += " cs-accordion-item--open"
		bodyStyle = "max-height:none;overflow:hidden;"
	}

	return fmt.Sprintf(`<div class="%s" data-id="%s">
  <button class="cs-accordion-trigger" data-accordion-trigger>
    <span>%s</span>
    <span class="cs-accordion-icon">&#9660;</span>
  </button>
  <div class="cs-accordion-body" style="%s">
    <div class="cs-accordion-content">%s</div>
  </div>
</div>`, cls, dataID, title, bodyStyle, children), nil
}

// ── Breadcrumb ──────────────────────────────────────────────────────────

func renderBreadcrumb(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "breadcrumb")
	return fmt.Sprintf(`<nav class="cs-breadcrumb" aria-label="Breadcrumb" data-id="%s">
  <ol class="cs-breadcrumb__list">%s</ol>
</nav>`, dataID, children), nil
}

func renderBreadcrumbItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	href := propStr(props, "href", "")
	dataID := propStr(props, "data-id", "breadcrumb-item")

	content := children
	if href != "" {
		content = fmt.Sprintf(`<a class="cs-breadcrumb__link" href="%s">%s</a>`, href, children)
	}

	return fmt.Sprintf(`<li class="cs-breadcrumb__item" data-id="%s">%s</li>`, dataID, content), nil
}

// ── Pagination ──────────────────────────────────────────────────────────

func renderPagination(props map[string]interface{}, children string, e *Engine) (string, error) {
	total := int(propFloat(props, "total", 0))
	page := int(propFloat(props, "page", 1))
	perPage := int(propFloat(props, "per-page", 10))
	action := propStr(props, "on:change", "")
	dataID := propStr(props, "data-id", "pagination")

	if perPage <= 0 {
		perPage = 10
	}
	totalPages := (total + perPage - 1) / perPage
	if totalPages <= 0 {
		totalPages = 1
	}

	onclick := func(p int) string {
		if action != "" {
			return fmt.Sprintf(` onclick="csAction('%s:%d',this)"`, action, p)
		}
		return ""
	}

	var pages strings.Builder

	prevDisabled := ""
	if page <= 1 {
		prevDisabled = " cs-pagination__btn--disabled"
	}
	pages.WriteString(fmt.Sprintf(`<button class="cs-pagination__btn%s"%s>&#8249;</button>`, prevDisabled, onclick(page-1)))

	for i := 1; i <= totalPages; i++ {
		if totalPages > 7 {
			if i != 1 && i != totalPages && (i < page-2 || i > page+2) {
				if i == page-3 || i == page+3 {
					pages.WriteString(`<span class="cs-pagination__ellipsis">…</span>`)
				}
				continue
			}
		}
		activeClass := ""
		if i == page {
			activeClass = " cs-pagination__btn--active"
		}
		pages.WriteString(fmt.Sprintf(`<button class="cs-pagination__btn%s"%s>%s</button>`,
			activeClass, onclick(i), strconv.Itoa(i)))
	}

	nextDisabled := ""
	if page >= totalPages {
		nextDisabled = " cs-pagination__btn--disabled"
	}
	pages.WriteString(fmt.Sprintf(`<button class="cs-pagination__btn%s"%s>&#8250;</button>`, nextDisabled, onclick(page+1)))

	return fmt.Sprintf(`<nav class="cs-pagination" data-id="%s">%s</nav>`, dataID, pages.String()), nil
}

// ── Stepper ─────────────────────────────────────────────────────────────

func renderStepper(props map[string]interface{}, children string, e *Engine) (string, error) {
	direction := propStr(props, "direction", "vertical")
	return fmt.Sprintf(`<div%s>%s</div>`, userAttrs(props, "cs-stepper cs-stepper--"+direction), children), nil
}

func renderStepperStep(props map[string]interface{}, children string, e *Engine) (string, error) {
	title := propStr(props, "title", "")
	step := propStr(props, "step", "")
	active := propBool(props, "active", false)
	done := propBool(props, "done", false)
	dataID := propStr(props, "data-id", "stepper-step")

	cls := "cs-stepper-step"
	if active {
		cls += " cs-stepper-step--active"
	}
	if done {
		cls += " cs-stepper-step--done"
	}

	circleContent := step
	if done {
		circleContent = "&#10003;"
	}

	titleHTML := ""
	if title != "" {
		titleHTML = fmt.Sprintf(`<div class="cs-stepper-step__title">%s</div>`, title)
	}

	descHTML := ""
	if children != "" {
		descHTML = fmt.Sprintf(`<div class="cs-stepper-step__desc">%s</div>`, children)
	}

	return fmt.Sprintf(`<div class="%s" data-id="%s">
  <div class="cs-stepper-step__track">
    <div class="cs-stepper-step__circle">%s</div>
    <div class="cs-stepper-step__connector"></div>
  </div>
  <div class="cs-stepper-step__content">%s%s</div>
</div>`, cls, dataID, circleContent, titleHTML, descHTML), nil
}

// ── Toolbar ─────────────────────────────────────────────────────────────

func renderToolbar(props map[string]interface{}, children string, e *Engine) (string, error) {
	title := propStr(props, "title", "")
	dataID := propStr(props, "data-id", "toolbar")
	bordered := propBool(props, "bordered", false)

	titleHTML := ""
	if title != "" {
		titleHTML = fmt.Sprintf(`<span class="cs-toolbar__title">%s</span>`, title)
	}

	cls := "cs-toolbar"
	if bordered {
		cls += " cs-toolbar--bordered"
	}

	return fmt.Sprintf(`<div class="%s" data-id="%s">
  <div class="cs-toolbar__start">%s</div>
  <div class="cs-toolbar__end">%s</div>
</div>`, cls, dataID, titleHTML, children), nil
}

// ── Bottom Navigation ───────────────────────────────────────────────────
// MUI: BottomNavigation — mobile-style bottom nav bar with icon+label actions.
// Props: value (selected), showLabels (bool)
// Children: bottom-nav-action atoms
//
// ["bottom-nav", { "value": "home", "showLabels": true },
//   ["bottom-nav-action", { "label": "Home", "icon": "home", "value": "home" }],
//   ["bottom-nav-action", { "label": "Search", "icon": "search", "value": "search" }],
//   ["bottom-nav-action", { "label": "Profile", "icon": "person", "value": "profile" }]
// ]

func renderBottomNav(props map[string]interface{}, children string, e *Engine) (string, error) {
	value := propStr(props, "value", "")
	showLabels, _ := props["showLabels"].(bool)
	dataID := propStr(props, "data-id", "")

	classes := "cs-bottom-nav"
	if showLabels {
		classes += " cs-bottom-nav--show-labels"
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<nav class="%s" data-value="%s"%s>%s</nav>`,
		classes, esc(value), idAttr, children), nil
}

func renderBottomNavAction(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	icon := propStr(props, "icon", "")
	value := propStr(props, "value", "")
	href := propStr(props, "href", "")
	selected, _ := props["selected"].(bool)

	classes := "cs-bottom-nav__action"
	if selected {
		classes += " cs-bottom-nav__action--selected"
	}

	iconHTML := ""
	if icon != "" {
		iconHTML = bsIcon(icon, 24)
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<span class="cs-bottom-nav__label">%s</span>`, esc(label))
	}

	if href != "" {
		return fmt.Sprintf(`<a class="%s" href="%s" data-value="%s">%s%s</a>`,
			classes, esc(href), esc(value), iconHTML, labelHTML), nil
	}

	return fmt.Sprintf(`<button type="button" class="%s" data-value="%s" onclick="csBottomNav(this)">%s%s</button>`,
		classes, esc(value), iconHTML, labelHTML), nil
}

// ── Menubar ─────────────────────────────────────────────────────────────
// MUI: Menubar — horizontal menu bar with dropdown submenus.
// Props: data-id
// Children: menubar-item atoms
//
// ["menubar", {},
//   ["menubar-item", { "label": "File" },
//     ["menu-item", { "label": "New", "on:click": "file/new" }],
//     ["menu-item", { "label": "Open", "on:click": "file/open" }],
//     ["divider"],
//     ["menu-item", { "label": "Exit" }]
//   ],
//   ["menubar-item", { "label": "Edit" },
//     ["menu-item", { "label": "Undo" }],
//     ["menu-item", { "label": "Redo" }]
//   ]
// ]

func renderMenubar(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "")
	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}
	return fmt.Sprintf(`<nav class="cs-menubar" role="menubar"%s>%s</nav>`, idAttr, children), nil
}

func renderMenubarItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")

	return fmt.Sprintf(`<div class="cs-menubar__item" role="menuitem">`+
		`<button type="button" class="cs-menubar__trigger" onclick="csMenubarToggle(this)">%s</button>`+
		`<div class="cs-menubar__dropdown" role="menu">%s</div>`+
		`</div>`, esc(label), children), nil
}

// ── Speed Dial ──────────────────────────────────────────────────────────
// MUI: SpeedDial — floating button that expands to reveal actions.
// Props: icon, openIcon, direction (up/down/left/right), open (bool), ariaLabel
// Children: speed-dial-action atoms
//
// ["speed-dial", { "icon": "add", "ariaLabel": "Actions" },
//   ["speed-dial-action", { "icon": "copy", "tooltipTitle": "Copy" }],
//   ["speed-dial-action", { "icon": "save", "tooltipTitle": "Save" }],
//   ["speed-dial-action", { "icon": "print", "tooltipTitle": "Print" }]
// ]

func renderSpeedDial(props map[string]interface{}, children string, e *Engine) (string, error) {
	icon := propStr(props, "icon", "add")
	direction := propStr(props, "direction", "up")
	ariaLabel := propStr(props, "ariaLabel", "Speed Dial")
	dataID := propStr(props, "data-id", "")

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	iconHTML := bsIcon(icon, 24)

	return fmt.Sprintf(`<div class="cs-speed-dial cs-speed-dial--%s" aria-label="%s"%s>`+
		`<button type="button" class="cs-fab cs-fab--medium cs-fab--primary cs-speed-dial__trigger" onclick="csSpeedDialToggle(this)">%s</button>`+
		`<div class="cs-speed-dial__actions">%s</div>`+
		`</div>`,
		esc(direction), esc(ariaLabel), idAttr, iconHTML, children), nil
}

func renderSpeedDialAction(props map[string]interface{}, children string, e *Engine) (string, error) {
	icon := propStr(props, "icon", "")
	tooltipTitle := propStr(props, "tooltipTitle", "")
	action := propStr(props, "on:click", "")

	onclick := ""
	if action != "" {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, esc(action))
	}

	iconHTML := bsIcon(icon, 20)
	tooltip := ""
	if tooltipTitle != "" {
		tooltip = fmt.Sprintf(`<span class="cs-speed-dial__tooltip">%s</span>`, esc(tooltipTitle))
	}

	return fmt.Sprintf(`<div class="cs-speed-dial__action">%s<button type="button" class="cs-fab cs-fab--small cs-fab--secondary"%s>%s</button></div>`,
		tooltip, onclick, iconHTML), nil
}
