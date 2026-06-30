package engine

import "fmt"

// ── Header ──────────────────────────────────────────────────────────────
// ["header", { "title": "Dashboard", "subtitle": "Overview" }]

func renderHeader(props map[string]interface{}, children string, e *Engine) (string, error) {
	title := propStr(props, "title", "")
	subtitle := propStr(props, "subtitle", "")

	inner := children
	if title != "" {
		sub := ""
		if subtitle != "" {
			sub = fmt.Sprintf(`<span class="header-subtitle">%s</span>`, subtitle)
		}
		inner = fmt.Sprintf(`<span class="header-title">%s</span>%s`, title, sub)
	}

	return fmt.Sprintf(`<div%s>%s</div>`, userAttrs(props, "cs-header"), inner), nil
}

// ── Footer ──────────────────────────────────────────────────────────────

func renderFooter(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<footer%s>%s</footer>`, userAttrs(props, "cs-footer"), children), nil
}

// ── Card ────────────────────────────────────────────────────────────────

func renderCard(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<div%s>%s</div>`, userAttrs(props, "cs-card"), children), nil
}

// ── Sidebar ─────────────────────────────────────────────────────────────
// ["sidebar", { "brand": "MyApp" }, ...nav-links]

func renderSidebar(props map[string]interface{}, children string, e *Engine) (string, error) {
	brand := propStr(props, "brand", "")
	dataID := propStr(props, "data-id", "sidebar")

	brandHTML := ""
	if brand != "" {
		brandHTML = fmt.Sprintf(`<div class="cs-sidebar__brand">%s</div>`, brand)
	}

	return fmt.Sprintf(`<aside class="cs-sidebar" data-id="%s">
  %s
  <nav class="cs-sidebar__nav">%s</nav>
</aside>`, dataID, brandHTML, children), nil
}

// ── Section ─────────────────────────────────────────────────────────────
// ["section", { "title": "Overview", "description": "Key metrics" }, ...children]

func renderSection(props map[string]interface{}, children string, e *Engine) (string, error) {
	title := propStr(props, "title", "")
	description := propStr(props, "description", "")
	dataID := propStr(props, "data-id", "section")

	titleHTML := ""
	if title != "" {
		titleHTML = fmt.Sprintf(`<h2 class="cs-section__title">%s</h2>`, title)
	}
	descHTML := ""
	if description != "" {
		descHTML = fmt.Sprintf(`<p class="cs-section__desc">%s</p>`, description)
	}
	headerHTML := ""
	if title != "" || description != "" {
		headerHTML = fmt.Sprintf(`<div class="cs-section__header">%s%s</div>`, titleHTML, descHTML)
	}

	return fmt.Sprintf(`<section class="cs-section" data-id="%s">%s<div class="cs-section__body">%s</div></section>`,
		dataID, headerHTML, children), nil
}

// ── Divider ─────────────────────────────────────────────────────────────
// ["divider"] or ["divider", { "label": "OR" }]

func renderDivider(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	if label != "" {
		return fmt.Sprintf(`<div class="cs-divider--labeled"><span class="cs-divider__line"></span><span class="cs-divider__label">%s</span><span class="cs-divider__line"></span></div>`, label), nil
	}
	return `<hr class="cs-divider">`, nil
}

// ── SplitView ───────────────────────────────────────────────────────────
// ["split-view", { "direction": "horizontal" }, ["split-pane", ...], ["split-pane", ...]]

func renderSplitView(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<div%s>%s</div>`, userAttrs(props, "cs-split-view"), children), nil
}

func renderSplitPane(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<div%s>%s</div>`, userAttrs(props, "cs-split-view__pane"), children), nil
}

// ── Paper ───────────────────────────────────────────────────────────────

func renderPaper(props map[string]interface{}, children string, e *Engine) (string, error) {
	elevation := propInt(props, "elevation", 0)
	if elevation == 0 {
		if _, ok := props["elevation"]; !ok {
			elevation = 1
		}
	}
	square, _ := props["square"].(bool)
	variant := propStr(props, "variant", "elevation")
	dataID := propStr(props, "data-id", "")

	classes := "cs-paper"
	if variant == "outlined" {
		classes += " cs-paper--outlined"
	} else {
		classes += fmt.Sprintf(" cs-paper--elevation-%d", clampInt(elevation, 0, 24))
	}
	if square {
		classes += " cs-paper--square"
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<div class="%s"%s>%s</div>`, classes, idAttr, children), nil
}

// ── AppBar ──────────────────────────────────────────────────────────────

func renderAppBar(props map[string]interface{}, children string, e *Engine) (string, error) {
	position := propStr(props, "position", "fixed")
	color := propStr(props, "color", "primary")
	elevation := propInt(props, "elevation", 0)
	if _, ok := props["elevation"]; !ok {
		elevation = 4
	}
	dataID := propStr(props, "data-id", "")

	classes := fmt.Sprintf("cs-app-bar cs-app-bar--position-%s cs-app-bar--color-%s", esc(position), esc(color))
	if elevation > 0 {
		classes += fmt.Sprintf(" cs-paper--elevation-%d", clampInt(elevation, 0, 24))
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<header class="%s"%s>%s</header>`, classes, idAttr, children), nil
}

// ── Box ─────────────────────────────────────────────────────────────────
// MUI: Box — generic wrapper. The building block for layout composition.
// Props: component (tag name), display, p, m, sx (inline style override)
//
// ["box", { "p": 4, "m": 2, "display": "flex" }, "Content"]

func renderBox(props map[string]interface{}, children string, e *Engine) (string, error) {
	component := propStr(props, "component", "div")
	dataID := propStr(props, "data-id", "")

	// Build inline style from spacing/layout props
	style := ""
	if p := propInt(props, "p", 0); p > 0 {
		style += fmt.Sprintf("padding:var(--spacing-%d);", p)
	}
	if px := propInt(props, "px", 0); px > 0 {
		style += fmt.Sprintf("padding-left:var(--spacing-%d);padding-right:var(--spacing-%d);", px, px)
	}
	if py := propInt(props, "py", 0); py > 0 {
		style += fmt.Sprintf("padding-top:var(--spacing-%d);padding-bottom:var(--spacing-%d);", py, py)
	}
	if m := propInt(props, "m", 0); m > 0 {
		style += fmt.Sprintf("margin:var(--spacing-%d);", m)
	}
	if mx := propInt(props, "mx", 0); mx > 0 {
		style += fmt.Sprintf("margin-left:var(--spacing-%d);margin-right:var(--spacing-%d);", mx, mx)
	}
	if my := propInt(props, "my", 0); my > 0 {
		style += fmt.Sprintf("margin-top:var(--spacing-%d);margin-bottom:var(--spacing-%d);", my, my)
	}
	if display := propStr(props, "display", ""); display != "" {
		style += fmt.Sprintf("display:%s;", display)
	}
	if gap := propInt(props, "gap", 0); gap > 0 {
		style += fmt.Sprintf("gap:var(--spacing-%d);", gap)
	}
	if flexDir := propStr(props, "flexDirection", ""); flexDir != "" {
		style += fmt.Sprintf("flex-direction:%s;", flexDir)
	}
	if alignItems := propStr(props, "alignItems", ""); alignItems != "" {
		style += fmt.Sprintf("align-items:%s;", alignItems)
	}
	if justifyContent := propStr(props, "justifyContent", ""); justifyContent != "" {
		style += fmt.Sprintf("justify-content:%s;", justifyContent)
	}
	if sx := propStr(props, "sx", ""); sx != "" {
		style += sx
	}

	styleAttr := ""
	if style != "" {
		styleAttr = fmt.Sprintf(` style="%s"`, style)
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<%s class="cs-box"%s%s>%s</%s>`,
		esc(component), styleAttr, idAttr, children, esc(component)), nil
}

// ── Grid (v2) ───────────────────────────────────────────────────────────
// MUI: Grid — responsive 12-column grid system.
// Props: container (bool), item (bool), spacing (0-12), columns (default 12),
//        xs/sm/md/lg/xl (column span or "auto"), direction (row/column)
//
// ["grid", { "container": true, "spacing": 3 },
//   ["grid", { "item": true, "xs": 12, "md": 6 }, "Left column"],
//   ["grid", { "item": true, "xs": 12, "md": 6 }, "Right column"]
// ]

func renderGrid(props map[string]interface{}, children string, e *Engine) (string, error) {
	container, _ := props["container"].(bool)
	item, _ := props["item"].(bool)
	spacing := propInt(props, "spacing", 0)
	direction := propStr(props, "direction", "row")
	dataID := propStr(props, "data-id", "")

	classes := "cs-grid"
	style := ""

	if container {
		classes += " cs-grid--container"
		style += fmt.Sprintf("flex-direction:%s;", direction)
		if spacing > 0 {
			style += fmt.Sprintf("gap:var(--spacing-%d);", spacing)
		}
	}

	if item {
		classes += " cs-grid--item"
		for _, bp := range []string{"xs", "sm", "md", "lg", "xl"} {
			if v, ok := props[bp]; ok {
				cols := toInt(v)
				if cols > 0 {
					classes += fmt.Sprintf(" cs-grid--%s-%d", bp, cols)
				}
			}
		}
	}

	styleAttr := ""
	if style != "" {
		styleAttr = fmt.Sprintf(` style="%s"`, style)
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<div class="%s"%s%s>%s</div>`, classes, styleAttr, idAttr, children), nil
}

// ── Stack ───────────────────────────────────────────────────────────────
// MUI: Stack — flex container with spacing between children.
// Props: direction (row/column), spacing (token scale), alignItems, justifyContent, divider (bool)
//
// ["stack", { "direction": "column", "spacing": 3 },
//   ["text", {}, "First"],
//   ["text", {}, "Second"],
//   ["text", {}, "Third"]
// ]

func renderStack(props map[string]interface{}, children string, e *Engine) (string, error) {
	direction := propStr(props, "direction", "column")
	spacing := propInt(props, "spacing", 0)
	alignItems := propStr(props, "alignItems", "")
	justifyContent := propStr(props, "justifyContent", "")
	dataID := propStr(props, "data-id", "")

	style := fmt.Sprintf("display:flex;flex-direction:%s;", direction)
	if spacing > 0 {
		style += fmt.Sprintf("gap:var(--spacing-%d);", spacing)
	}
	if alignItems != "" {
		style += fmt.Sprintf("align-items:%s;", alignItems)
	}
	if justifyContent != "" {
		style += fmt.Sprintf("justify-content:%s;", justifyContent)
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<div class="cs-stack" style="%s"%s>%s</div>`, style, idAttr, children), nil
}

// ── Image List ──────────────────────────────────────────────────────────
// MUI: ImageList — grid of images with optional titles/subtitles.
// Props: cols (default 3), gap, variant (standard/quilted/woven/masonry)
// Children: image-list-item atoms
//
// ["image-list", { "cols": 3, "gap": 2 },
//   ["image-list-item", { "src": "/img/1.jpg", "title": "Beach" }],
//   ["image-list-item", { "src": "/img/2.jpg", "title": "Mountain" }]
// ]

func renderImageList(props map[string]interface{}, children string, e *Engine) (string, error) {
	cols := propInt(props, "cols", 0)
	if cols == 0 {
		cols = 3
	}
	gap := propInt(props, "gap", 0)
	if gap == 0 {
		gap = 1
	}
	variant := propStr(props, "variant", "standard")
	dataID := propStr(props, "data-id", "")

	classes := fmt.Sprintf("cs-image-list cs-image-list--%s", esc(variant))

	style := fmt.Sprintf("grid-template-columns:repeat(%d,1fr);gap:var(--spacing-%d);", cols, gap)

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<div class="%s" style="%s"%s>%s</div>`, classes, style, idAttr, children), nil
}

func renderImageListItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	src := propStr(props, "src", "")
	title := propStr(props, "title", "")
	alt := propStr(props, "alt", title)
	rows := propInt(props, "rows", 0)
	cols := propInt(props, "cols", 0)

	style := ""
	if rows > 1 {
		style += fmt.Sprintf("grid-row:span %d;", rows)
	}
	if cols > 1 {
		style += fmt.Sprintf("grid-column:span %d;", cols)
	}
	styleAttr := ""
	if style != "" {
		styleAttr = fmt.Sprintf(` style="%s"`, style)
	}

	img := ""
	if src != "" {
		img = fmt.Sprintf(`<img class="cs-image-list__img" src="%s" alt="%s" loading="lazy">`, esc(src), esc(alt))
	}

	titleBar := ""
	if title != "" {
		titleBar = fmt.Sprintf(`<div class="cs-image-list__title-bar"><span>%s</span></div>`, esc(title))
	}

	return fmt.Sprintf(`<div class="cs-image-list__item"%s>%s%s%s</div>`, styleAttr, img, titleBar, children), nil
}

// ── Masonry ─────────────────────────────────────────────────────────────
// MUI Lab: Masonry — Pinterest-style layout using CSS columns.
// Props: columns (default 3), spacing
//
// ["masonry", { "columns": 3, "spacing": 2 },
//   ["paper", { "elevation": 1 }, "Item 1"],
//   ["paper", { "elevation": 1 }, "Item 2"],
//   ...
// ]

func renderMasonry(props map[string]interface{}, children string, e *Engine) (string, error) {
	columns := propInt(props, "columns", 0)
	if columns == 0 {
		columns = 3
	}
	spacing := propInt(props, "spacing", 0)
	if spacing == 0 {
		spacing = 2
	}
	dataID := propStr(props, "data-id", "")

	style := fmt.Sprintf("column-count:%d;column-gap:var(--spacing-%d);", columns, spacing)

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<div class="cs-masonry" style="%s"%s>%s</div>`, style, idAttr, children), nil
}
