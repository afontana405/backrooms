package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// ── StatCard ────────────────────────────────────────────────────────────

func renderStatCard(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	value := propStr(props, "value", "0")
	trend := propStr(props, "trend", "")

	trendHTML := ""
	if trend == "up" {
		trendHTML = `<div class="trend-up">↑</div>`
	} else if trend == "down" {
		trendHTML = `<div class="trend-down">↓</div>`
	}

	return fmt.Sprintf(`<div%s>
  <div class="stat-label">%s</div>
  <div class="stat-value">%s</div>%s
</div>`, userAttrs(props, "cs-stat-card"), label, value, trendHTML), nil
}

// ── Table ───────────────────────────────────────────────────────────────

func renderTable(props map[string]interface{}, children string, e *Engine) (string, error) {
	striped := propBool(props, "striped", true)
	hoverable := propBool(props, "hoverable", true)

	cls := "table"
	if striped {
		cls += " table-striped"
	}
	if hoverable {
		cls += " table-hover"
	}

	var headerHTML strings.Builder
	if cols, ok := props["columns"]; ok {
		if colList, ok := cols.([]interface{}); ok {
			headerHTML.WriteString("<thead><tr>")
			for _, col := range colList {
				headerHTML.WriteString(fmt.Sprintf("<th>%v</th>", col))
			}
			headerHTML.WriteString("</tr></thead>")
		}
	}

	var bodyHTML strings.Builder
	bodyHTML.WriteString("<tbody>")
	if rows, ok := props["rows"]; ok {
		if rowList, ok := rows.([]interface{}); ok {
			for _, row := range rowList {
				bodyHTML.WriteString("<tr>")
				if cells, ok := row.([]interface{}); ok {
					for _, cell := range cells {
						bodyHTML.WriteString(fmt.Sprintf("<td>%v</td>", cell))
					}
				}
				bodyHTML.WriteString("</tr>")
			}
		}
	}
	if children != "" {
		bodyHTML.WriteString(children)
	}
	bodyHTML.WriteString("</tbody>")

	return fmt.Sprintf(`<div%s>
  <table class="%s">
    %s
    %s
  </table>
</div>`, userAttrs(props, "table-responsive"), cls, headerHTML.String(), bodyHTML.String()), nil
}

// ── List ────────────────────────────────────────────────────────────────

func renderList(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<div%s>%s</div>`, userAttrs(props, "cs-list"), children), nil
}

func renderListItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	href := propStr(props, "href", "")
	action := propStr(props, "on:click", "")

	onclick := ""
	if action != "" {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, action)
	}

	content := children
	if href != "" {
		return fmt.Sprintf(`<a class="cs-list-item cs-list__link" href="%s"%s>%s</a>`, href, onclick, content), nil
	}

	return fmt.Sprintf(`<div%s%s>%s</div>`, userAttrs(props, "cs-list-item"), onclick, content), nil
}

// ── KvList ──────────────────────────────────────────────────────────────

func renderKvList(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "kv-list")
	divided := propBool(props, "divided", true)

	cls := "cs-kv-list"
	if divided {
		cls += " cs-kv-list--divided"
	}

	return fmt.Sprintf(`<dl class="%s" data-id="%s">%s</dl>`, cls, dataID, children), nil
}

func renderKvItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	key := propStr(props, "key", "")
	value := propStr(props, "value", children)
	href := propStr(props, "href", "")
	valueVariant := propStr(props, "value-variant", "")
	dataID := propStr(props, "data-id", "kv-item")

	valueHTML := value
	if href != "" {
		valueHTML = fmt.Sprintf(`<a class="cs-kv-item__link" href="%s">%s</a>`, href, value)
	}

	variantClass := ""
	if valueVariant != "" {
		variantClass = fmt.Sprintf(` cs-kv-item__value--%s`, valueVariant)
	}

	return fmt.Sprintf(`<div class="cs-kv-item" data-id="%s">
  <dt class="cs-kv-item__key">%s</dt>
  <dd class="cs-kv-item__value%s">%s</dd>
</div>`, dataID, key, variantClass, valueHTML), nil
}

// ── DataGrid ────────────────────────────────────────────────────────────

func renderDataGrid(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "data-grid")
	sortable := propBool(props, "sortable", true)
	filterable := propBool(props, "filterable", false)
	emptyMsg := propStr(props, "empty", "No results")

	filterHTML := ""
	if filterable {
		filterHTML = fmt.Sprintf(`<div class="cs-data-grid__filter">
  <input class="cs-data-grid__filter-input" type="search" placeholder="Filter..."
    data-id="%s--filter" oninput="csDataGridFilter(this)" />
</div>`, dataID)
	}

	var headerHTML strings.Builder
	if cols, ok := props["columns"]; ok {
		if colList, ok := cols.([]interface{}); ok {
			headerHTML.WriteString("<thead><tr>")
			for i, col := range colList {
				if sortable {
					headerHTML.WriteString(fmt.Sprintf(
						`<th class="cs-data-grid__th" data-col-idx="%d" onclick="csDataGridSort(this)">%v <span class="cs-data-grid__sort-icon">&#8597;</span></th>`,
						i, col))
				} else {
					headerHTML.WriteString(fmt.Sprintf("<th>%v</th>", col))
				}
			}
			headerHTML.WriteString("</tr></thead>")
		}
	}

	var bodyHTML strings.Builder
	bodyHTML.WriteString("<tbody>")
	if rows, ok := props["rows"]; ok {
		if rowList, ok := rows.([]interface{}); ok {
			for _, row := range rowList {
				bodyHTML.WriteString("<tr>")
				if cells, ok := row.([]interface{}); ok {
					for _, cell := range cells {
						bodyHTML.WriteString(fmt.Sprintf("<td>%v</td>", cell))
					}
				}
				bodyHTML.WriteString("</tr>")
			}
		}
	}
	if children != "" {
		bodyHTML.WriteString(children)
	}
	bodyHTML.WriteString("</tbody>")

	emptyHTML := fmt.Sprintf(`<div class="cs-data-grid__empty" data-id="%s--empty">%s</div>`, dataID, emptyMsg)

	return fmt.Sprintf(`<div class="cs-data-grid" data-id="%s">
  %s
  <div class="cs-data-grid__wrap">
    <table class="cs-table cs-data-grid__table">
      %s
      %s
    </table>
  </div>
  %s
</div>`, dataID, filterHTML, headerHTML.String(), bodyHTML.String(), emptyHTML), nil
}

// ── Tree ────────────────────────────────────────────────────────────────

func renderTree(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "tree")
	return fmt.Sprintf(`<ul class="cs-tree" data-id="%s">%s</ul>`, dataID, children), nil
}

func renderTreeItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	icon := propStr(props, "icon", "")
	active := propBool(props, "active", false)
	open := propBool(props, "open", false)
	dataID := propStr(props, "data-id", "tree-item")
	hasChildren := children != ""

	cls := "cs-tree-item"
	if active {
		cls += " cs-tree-item--active"
	}
	if open && hasChildren {
		cls += " cs-tree-item--open"
	}

	iconHTML := ""
	if icon != "" {
		iconHTML, _ = renderIcon(map[string]interface{}{"name": icon, "size": float64(14)}, "", e)
	}

	chevron := ""
	rowOnclick := ""
	if hasChildren {
		chevron = `<span class="cs-tree-item__chevron">&#9660;</span>`
		rowOnclick = ` onclick="csTreeToggle(this)"`
	}

	nested := ""
	if hasChildren {
		display := "none"
		if open {
			display = ""
		}
		nested = fmt.Sprintf(`<ul class="cs-tree-item__children" style="display:%s">%s</ul>`, display, children)
	}

	return fmt.Sprintf(`<li class="%s" data-id="%s">
  <div class="cs-tree-item__row"%s>%s%s<span class="cs-tree-item__label">%s</span></div>
  %s
</li>`, cls, dataID, rowOnclick, chevron, iconHTML, label, nested), nil
}

// ── VirtualList ─────────────────────────────────────────────────────────

func renderVirtualList(props map[string]interface{}, children string, e *Engine) (string, error) {
	dataID := propStr(props, "data-id", "virtual-list")
	height := int(propFloat(props, "height", 400))
	rowHeight := int(propFloat(props, "row-height", 40))

	colsJSON := "[]"
	if cols, ok := props["columns"]; ok {
		if b, err := json.Marshal(cols); err == nil {
			colsJSON = string(b)
		}
	}

	rowsJSON := "[]"
	if rows, ok := props["rows"]; ok {
		if b, err := json.Marshal(rows); err == nil {
			rowsJSON = string(b)
		}
	}

	return fmt.Sprintf(`<div class="cs-virtual-list" data-id="%s"
  style="height:%dpx;overflow-y:auto;position:relative;"
  onscroll="csVirtualListScroll(this)">
  <div class="cs-virtual-list__inner" style="position:relative;"></div>
</div>
<script>csVirtualListInit('%s',%s,%s,%d);</script>`,
		dataID, height, dataID, colsJSON, rowsJSON, rowHeight), nil
}

// ── Chart ───────────────────────────────────────────────────────────────

var chartPalette = []string{
	"var(--accent)", "var(--color-success)", "var(--color-warning)",
	"var(--color-danger)", "var(--color-info)", "#8b5cf6", "#f97316", "#06b6d4",
}

type chartPoint struct {
	Label string
	Value float64
}

func parseChartData(props map[string]interface{}) []chartPoint {
	var points []chartPoint
	raw, ok := props["data"]
	if !ok {
		return points
	}
	switch v := raw.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				p := chartPoint{}
				if l, ok := m["label"].(string); ok {
					p.Label = l
				}
				if val, ok := m["value"].(float64); ok {
					p.Value = val
				}
				points = append(points, p)
			}
		}
	case string:
		var parsed []map[string]interface{}
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			for _, m := range parsed {
				p := chartPoint{}
				if l, ok := m["label"].(string); ok {
					p.Label = l
				}
				if val, ok := m["value"].(float64); ok {
					p.Value = val
				}
				points = append(points, p)
			}
		}
	}
	return points
}

func chartMax(points []chartPoint) float64 {
	max := 0.0
	for _, p := range points {
		if p.Value > max {
			max = p.Value
		}
	}
	return max
}

func formatChartNum(v float64) string {
	if v >= 1000000 {
		return fmt.Sprintf("%.1fM", v/1000000)
	}
	if v >= 1000 {
		return fmt.Sprintf("%.1fk", v/1000)
	}
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return fmt.Sprintf("%.1f", v)
}

func renderChart(props map[string]interface{}, children string, e *Engine) (string, error) {
	chartType := propStr(props, "type", "bar")
	height := int(propFloat(props, "height", 180))
	dataID := propStr(props, "data-id", "chart")
	title := propStr(props, "title", "")

	points := parseChartData(props)
	if len(points) == 0 {
		return fmt.Sprintf(`<div class="cs-chart cs-chart--empty" data-id="%s"><span>No data</span></div>`, dataID), nil
	}

	titleHTML := ""
	if title != "" {
		titleHTML = fmt.Sprintf(`<div class="cs-chart__title">%s</div>`, title)
	}

	var svg string
	var err error
	switch chartType {
	case "line":
		svg, err = renderLineChart(points, height)
	case "pie":
		svg, err = renderPieChart(points, height)
	default:
		svg, err = renderBarChart(points, height)
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`<div class="cs-chart" data-id="%s">%s%s</div>`, dataID, titleHTML, svg), nil
}

func renderBarChart(points []chartPoint, height int) (string, error) {
	n := len(points)
	padL, padR, padT, padB := 48.0, 16.0, 16.0, 40.0
	W := 500.0
	chartW := W - padL - padR
	chartH := float64(height)
	H := chartH + padT + padB

	maxVal := chartMax(points)
	if maxVal == 0 {
		maxVal = 1
	}

	barSlot := chartW / float64(n)
	barW := barSlot * 0.55

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg class="cs-chart__svg" viewBox="0 0 %.0f %.0f">`, W, H))

	for i := 0; i <= 4; i++ {
		yVal := maxVal * float64(i) / 4.0
		y := padT + chartH - (yVal/maxVal)*chartH
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" class="cs-chart__grid"/>`,
			padL, y, W-padR, y))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" class="cs-chart__axis-label" text-anchor="end" dominant-baseline="middle">%s</text>`,
			padL-6, y, formatChartNum(yVal)))
	}

	for i, p := range points {
		bH := (p.Value / maxVal) * chartH
		x := padL + float64(i)*barSlot + (barSlot-barW)/2
		y := padT + chartH - bH
		color := chartPalette[i%len(chartPalette)]
		sb.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" rx="3" class="cs-chart__bar"/>`,
			x, y, barW, bH, color))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" class="cs-chart__value" text-anchor="middle">%s</text>`,
			x+barW/2, y-5, formatChartNum(p.Value)))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" class="cs-chart__label" text-anchor="middle">%s</text>`,
			x+barW/2, padT+chartH+24, p.Label))
	}

	sb.WriteString(`</svg>`)
	return sb.String(), nil
}

func renderLineChart(points []chartPoint, height int) (string, error) {
	n := len(points)
	padL, padR, padT, padB := 48.0, 16.0, 16.0, 40.0
	W := 500.0
	chartW := W - padL - padR
	chartH := float64(height)
	H := chartH + padT + padB

	maxVal := chartMax(points)
	if maxVal == 0 {
		maxVal = 1
	}

	xStep := chartW
	if n > 1 {
		xStep = chartW / float64(n-1)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg class="cs-chart__svg" viewBox="0 0 %.0f %.0f">`, W, H))

	for i := 0; i <= 4; i++ {
		yVal := maxVal * float64(i) / 4.0
		y := padT + chartH - (yVal/maxVal)*chartH
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" class="cs-chart__grid"/>`,
			padL, y, W-padR, y))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" class="cs-chart__axis-label" text-anchor="end" dominant-baseline="middle">%s</text>`,
			padL-6, y, formatChartNum(yVal)))
	}

	var area, line strings.Builder
	for i, p := range points {
		x := padL + float64(i)*xStep
		y := padT + chartH - (p.Value/maxVal)*chartH
		if i == 0 {
			area.WriteString(fmt.Sprintf("M %.1f %.1f", x, y))
			line.WriteString(fmt.Sprintf("M %.1f %.1f", x, y))
		} else {
			area.WriteString(fmt.Sprintf(" L %.1f %.1f", x, y))
			line.WriteString(fmt.Sprintf(" L %.1f %.1f", x, y))
		}
	}
	lastX := padL + float64(n-1)*xStep
	area.WriteString(fmt.Sprintf(" L %.1f %.1f L %.1f %.1f Z", lastX, padT+chartH, padL, padT+chartH))

	sb.WriteString(fmt.Sprintf(`<path d="%s" class="cs-chart__area"/>`, area.String()))
	sb.WriteString(fmt.Sprintf(`<path d="%s" class="cs-chart__line" fill="none"/>`, line.String()))

	for i, p := range points {
		x := padL + float64(i)*xStep
		y := padT + chartH - (p.Value/maxVal)*chartH
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="4" class="cs-chart__dot"/>`, x, y))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" class="cs-chart__label" text-anchor="middle">%s</text>`,
			x, padT+chartH+24, p.Label))
	}

	sb.WriteString(`</svg>`)
	return sb.String(), nil
}

func renderPieChart(points []chartPoint, height int) (string, error) {
	cx, cy, r := 100.0, 100.0, 82.0
	W := 320.0
	H := math.Max(float64(height+20), 220)

	total := 0.0
	for _, p := range points {
		total += p.Value
	}
	if total == 0 {
		total = 1
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg class="cs-chart__svg" viewBox="0 0 %.0f %.0f">`, W, H))

	angle := -math.Pi / 2
	for i, p := range points {
		slice := (p.Value / total) * 2 * math.Pi
		end := angle + slice
		color := chartPalette[i%len(chartPalette)]

		x1 := cx + r*math.Cos(angle)
		y1 := cy + r*math.Sin(angle)
		x2 := cx + r*math.Cos(end)
		y2 := cy + r*math.Sin(end)
		largeArc := 0
		if slice > math.Pi {
			largeArc = 1
		}

		sb.WriteString(fmt.Sprintf(
			`<path d="M %.1f %.1f L %.1f %.1f A %.1f %.1f 0 %d 1 %.1f %.1f Z" fill="%s" class="cs-chart__slice"/>`,
			cx, cy, x1, y1, r, r, largeArc, x2, y2, color))

		angle = end
	}

	for i, p := range points {
		color := chartPalette[i%len(chartPalette)]
		ly := 20.0 + float64(i)*22
		pct := (p.Value / total) * 100
		sb.WriteString(fmt.Sprintf(`<rect x="200" y="%.1f" width="10" height="10" fill="%s" rx="2"/>`, ly, color))
		sb.WriteString(fmt.Sprintf(`<text x="215" y="%.1f" class="cs-chart__legend-label">%s (%.0f%%)</text>`,
			ly+9, p.Label, pct))
	}

	sb.WriteString(`</svg>`)
	return sb.String(), nil
}

// ── Calendar ────────────────────────────────────────────────────────────

func renderCalendar(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "cal")
	name := propStr(props, "name", "date")
	value := propStr(props, "value", "")
	action := propStr(props, "on:select", "")
	dataID := propStr(props, "data-id", "calendar--"+id)

	return fmt.Sprintf(`<div class="cs-calendar" id="%s" data-action="%s" data-id="%s">
  <div class="cs-calendar__header">
    <button class="cs-calendar__nav" type="button" onclick="csCalendarNav('%s',-1)" data-id="%s--prev">&#8249;</button>
    <span class="cs-calendar__label" data-cal-label="%s"></span>
    <button class="cs-calendar__nav" type="button" onclick="csCalendarNav('%s',1)" data-id="%s--next">&#8250;</button>
  </div>
  <div class="cs-calendar__grid" data-cal-grid="%s"></div>
  <input type="hidden" name="%s" value="%s" data-cal-value="%s" />
</div>
<script>csCalendarInit('%s','%s');</script>`,
		id, action, dataID,
		id, dataID,
		id,
		id, dataID,
		id,
		name, value, id,
		id, value), nil
}
