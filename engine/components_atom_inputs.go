package engine

import (
	"fmt"
	"strings"
)

// ── Input ───────────────────────────────────────────────────────────────
// ["input", { "label": "Email", "type": "email", "name": "email" }]

func renderInput(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	inputType := propStr(props, "type", "text")
	name := propStr(props, "name", label)
	value := propAnyStr(props, "value")
	id := propStr(props, "id", "input-"+name)
	onEnter := propStr(props, "on:enter", "")
	onTab := propStr(props, "on:tab", "")
	clear := propBool(props, "clear", false)

	valueAttr := ""
	if value != "" {
		valueAttr = fmt.Sprintf(` value="%s"`, strings.ReplaceAll(value, `"`, "&quot;"))
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-input__label cs-input__label--float" for="%s">%s</label>`, id, label)
	}

	hint := propStr(props, "hint", "")
	hintHTML := ""
	if hint != "" {
		hintHTML = fmt.Sprintf(`<span class="cs-input__hint">%s</span>`, hint)
	}

	pipeAttrs := ""
	if onEnter != "" {
		pipeAttrs += fmt.Sprintf(` data-on-enter="%s"`, onEnter)
		if clear {
			pipeAttrs += ` data-clear`
		}
	}
	if onTab != "" {
		pipeAttrs += fmt.Sprintf(` data-on-tab="%s"`, onTab)
	}
	if onEnter != "" || onTab != "" {
		for k, v := range props {
			if strings.HasPrefix(k, "data-") && k != "data-id" && k != "data-on-enter" && k != "data-clear" && k != "data-on-tab" {
				if s, ok := v.(string); ok {
					pipeAttrs += fmt.Sprintf(` %s="%s"`, k, s)
				}
			}
		}
	}

	return fmt.Sprintf(`<div%s><div class="cs-input__wrap"><input type="%s" class="cs-input__field" id="%s" name="%s" placeholder=" "%s%s>%s</div>%s</div>`,
		userAttrs(props, "cs-input"), inputType, id, name, valueAttr, pipeAttrs, labelHTML, hintHTML), nil
}

// ── Textarea ────────────────────────────────────────────────────────────

func renderTextarea(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", label)
	rows := int(propFloat(props, "rows", 4))
	hint := propStr(props, "hint", "")
	dataID := propStr(props, "data-id", "textarea--"+name)
	disabled := propBool(props, "disabled", false)

	id := "cs-textarea-" + name

	disabledAttr := ""
	if disabled {
		disabledAttr = " disabled"
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-input__label" for="%s">%s</label>`, id, label)
	}

	hintHTML := ""
	if hint != "" {
		hintHTML = fmt.Sprintf(`<span class="cs-input__hint">%s</span>`, hint)
	}

	return fmt.Sprintf(`<div class="cs-input cs-textarea" data-id="%s">
  <div class="cs-input__wrap">
    <textarea class="cs-input__field cs-textarea__field" id="%s" name="%s" rows="%d" placeholder=" "%s>%s</textarea>
    %s
  </div>
  %s
</div>`, dataID, id, name, rows, disabledAttr, children, labelHTML, hintHTML), nil
}

// ── Select ──────────────────────────────────────────────────────────────

func renderSelect(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", label)
	dataID := propStr(props, "data-id", "select--"+name)
	placeholder := propStr(props, "placeholder", "Select...")

	var optionsHTML strings.Builder
	if opts, ok := props["options"]; ok {
		switch v := opts.(type) {
		case []interface{}:
			for _, o := range v {
				optStr := fmt.Sprintf("%v", o)
				optionsHTML.WriteString(fmt.Sprintf(`<div class="cs-select__option" data-select-option="%s">%s</div>`, optStr, optStr))
			}
		}
	}
	if children != "" {
		optionsHTML.WriteString(children)
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-select__label">%s</label>`, label)
	}

	return fmt.Sprintf(`<div class="cs-select" data-id="%s">
  %s
  <div class="cs-select__trigger" data-select-trigger>
    <span class="cs-select__value">%s</span>
    <span class="cs-select__arrow">&#9660;</span>
  </div>
  <div class="cs-select__dropdown" style="display:none">%s</div>
  <input type="hidden" name="%s" />
</div>`, dataID, labelHTML, placeholder, optionsHTML.String(), name), nil
}

// ── NativeSelect ────────────────────────────────────────────────────────

func renderNativeSelect(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", label)
	placeholder := propStr(props, "placeholder", "Select...")
	dataID := propStr(props, "data-id", "native-select--"+name)
	id := "ns-" + name

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-input__label" for="%s">%s</label>`, id, label)
	}

	var optionsHTML strings.Builder
	optionsHTML.WriteString(fmt.Sprintf(`<option value="" disabled selected>%s</option>`, placeholder))
	if opts, ok := props["options"]; ok {
		switch v := opts.(type) {
		case []interface{}:
			for _, o := range v {
				optStr := fmt.Sprintf("%v", o)
				optionsHTML.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, optStr, optStr))
			}
		}
	}

	return fmt.Sprintf(`<div class="cs-input" data-id="%s">%s<select class="cs-input__field" id="%s" name="%s">%s</select></div>`,
		dataID, labelHTML, id, name, optionsHTML.String()), nil
}

// ── Autocomplete ────────────────────────────────────────────────────────

func renderAutocomplete(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", label)
	placeholder := propStr(props, "placeholder", label)
	dataID := propStr(props, "data-id", "autocomplete--"+name)

	id := "cs-ac-" + name

	var optionsHTML strings.Builder
	if opts, ok := props["options"]; ok {
		switch v := opts.(type) {
		case []interface{}:
			for _, o := range v {
				optStr := fmt.Sprintf("%v", o)
				optionsHTML.WriteString(fmt.Sprintf(`<div class="cs-autocomplete__item" data-ac-item>%s</div>`, optStr))
			}
		}
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-input__label cs-input__label--float" for="%s">%s</label>`, id, label)
	}

	return fmt.Sprintf(`<div class="cs-autocomplete cs-input" data-id="%s">
  <div class="cs-input__wrap">
    <input class="cs-input__field" id="%s" name="%s" placeholder="%s" autocomplete="off" data-autocomplete />
    %s
  </div>
  <div class="cs-autocomplete__dropdown" style="display:none">%s</div>
</div>`, dataID, id, name, placeholder, labelHTML, optionsHTML.String()), nil
}

// ── FormField ───────────────────────────────────────────────────────────

func renderFormField(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	hint := propStr(props, "hint", "")
	errMsg := propStr(props, "error", "")
	required := propBool(props, "required", false)
	dataID := propStr(props, "data-id", "form-field")

	requiredMark := ""
	if required {
		requiredMark = ` <span class="cs-form-field__required">*</span>`
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-form-field__label">%s%s</label>`, label, requiredMark)
	}

	hintHTML := ""
	if errMsg != "" {
		hintHTML = fmt.Sprintf(`<span class="cs-form-field__hint cs-form-field__hint--error">%s</span>`, errMsg)
	} else if hint != "" {
		hintHTML = fmt.Sprintf(`<span class="cs-form-field__hint">%s</span>`, hint)
	}

	errClass := ""
	if errMsg != "" {
		errClass = " cs-form-field--error"
	}

	return fmt.Sprintf(`<div class="cs-form-field%s" data-id="%s">%s%s%s</div>`,
		errClass, dataID, labelHTML, children, hintHTML), nil
}

// ── MultiSelect ─────────────────────────────────────────────────────────

func renderMultiSelect(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", "multi")
	placeholder := propStr(props, "placeholder", "Select...")
	dataID := propStr(props, "data-id", "multi-select--"+name)

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-multi-select__label">%s</label>`, label)
	}

	var optionsHTML strings.Builder
	if opts, ok := props["options"]; ok {
		if optList, ok := opts.([]interface{}); ok {
			for _, o := range optList {
				v := fmt.Sprintf("%v", o)
				optionsHTML.WriteString(fmt.Sprintf(
					`<div class="cs-multi-select__option" data-ms-option="%s"
  onclick="csMultiSelectToggle(this.closest('[data-ms-wrap]'),'%s','%s')">%s</div>`,
					v, v, v, v))
			}
		}
	}

	return fmt.Sprintf(`<div class="cs-multi-select" data-ms-wrap data-id="%s">
  %s
  <div class="cs-multi-select__control" onclick="csMultiSelectOpen(this.closest('[data-ms-wrap]'))">
    <div class="cs-multi-select__tags" data-ms-tags>
      <span class="cs-multi-select__placeholder" data-ms-placeholder>%s</span>
    </div>
    <span class="cs-multi-select__arrow">&#9660;</span>
  </div>
  <div class="cs-multi-select__dropdown" data-ms-dropdown style="display:none">%s</div>
  <input type="hidden" name="%s" data-ms-value value="" />
</div>`, dataID, labelHTML, placeholder, optionsHTML.String(), name), nil
}

// ── Checkbox ────────────────────────────────────────────────────────────

func renderCheckbox(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	name := propStr(props, "name", "")
	checked := propBool(props, "checked", false)
	disabled := propBool(props, "disabled", false)
	dataID := propStr(props, "data-id", "checkbox--"+name)

	checkedAttr := ""
	if checked {
		checkedAttr = " checked"
	}
	disabledAttr := ""
	if disabled {
		disabledAttr = " disabled"
	}

	return fmt.Sprintf(`<label class="cs-checkbox" data-id="%s">
  <input class="cs-checkbox__input" type="checkbox" name="%s"%s%s />
  <span class="cs-checkbox__box"></span>
  <span class="cs-checkbox__label">%s</span>
</label>`, dataID, name, checkedAttr, disabledAttr, label), nil
}

// ── Radio ───────────────────────────────────────────────────────────────

func renderRadio(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	name := propStr(props, "name", "")
	value := propStr(props, "value", label)
	checked := propBool(props, "checked", false)
	dataID := propStr(props, "data-id", "radio--"+name+"--"+value)

	checkedAttr := ""
	if checked {
		checkedAttr = " checked"
	}

	return fmt.Sprintf(`<label class="cs-radio" data-id="%s">
  <input class="cs-radio__input" type="radio" name="%s" value="%s"%s />
  <span class="cs-radio__dot"></span>
  <span class="cs-radio__label">%s</span>
</label>`, dataID, name, value, checkedAttr, label), nil
}

// ── Switch ──────────────────────────────────────────────────────────────

func renderSwitch(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	name := propStr(props, "name", "")
	checked := propBool(props, "checked", false)
	dataID := propStr(props, "data-id", "switch--"+name)

	checkedAttr := ""
	if checked {
		checkedAttr = " checked"
	}

	return fmt.Sprintf(`<label class="cs-switch" data-id="%s">
  <input class="cs-switch__input" type="checkbox" name="%s"%s />
  <span class="cs-switch__track">
    <span class="cs-switch__thumb"></span>
  </span>
  <span class="cs-switch__label">%s</span>
</label>`, dataID, name, checkedAttr, label), nil
}

// ── Form ────────────────────────────────────────────────────────────────

func renderForm(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "")
	autosave := propStr(props, "data-autosave", "")
	dataID := propStr(props, "data-id", "form")

	idAttr := ""
	if id != "" {
		idAttr = fmt.Sprintf(` id="%s"`, id)
	}
	autosaveAttr := ""
	if autosave != "" {
		autosaveAttr = fmt.Sprintf(` data-autosave="%s"`, autosave)
	}

	return fmt.Sprintf(`<form class="cs-form"%s%s data-id="%s" onsubmit="event.preventDefault();var b=this.querySelector('[onclick]');if(b)b.click();">%s</form>`,
		idAttr, autosaveAttr, dataID, children), nil
}

// ── Slider ──────────────────────────────────────────────────────────────

func renderSlider(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", "slider")
	min := int(propFloat(props, "min", 0))
	max := int(propFloat(props, "max", 100))
	value := int(propFloat(props, "value", 50))
	step := int(propFloat(props, "step", 1))
	dataID := propStr(props, "data-id", "slider")

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<div class="cs-slider__header">
    <label class="cs-slider__label">%s</label>
    <span class="cs-slider__value" data-slider-value="%s">%d</span>
  </div>`, label, dataID, value)
	}

	return fmt.Sprintf(`<div class="cs-slider-wrap" data-id="%s">
  %s
  <input type="range" class="cs-slider" name="%s"
    min="%d" max="%d" value="%d" step="%d"
    data-slider-id="%s"
    oninput="csSliderUpdate(this)" />
</div>`, dataID, labelHTML, name, min, max, value, step, dataID), nil
}

// ── NumberInput ─────────────────────────────────────────────────────────

func renderNumberInput(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", "number")
	min := propStr(props, "min", "")
	max := propStr(props, "max", "")
	value := int(propFloat(props, "value", 0))
	step := int(propFloat(props, "step", 1))
	dataID := propStr(props, "data-id", "number-input")

	minAttr := ""
	if min != "" {
		minAttr = fmt.Sprintf(` min="%s"`, min)
	}
	maxAttr := ""
	if max != "" {
		maxAttr = fmt.Sprintf(` max="%s"`, max)
	}

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-number-input__label">%s</label>`, label)
	}

	return fmt.Sprintf(`<div class="cs-number-input-wrap" data-id="%s">
  %s
  <div class="cs-number-input">
    <button type="button" class="cs-number-input__btn" data-id="%s--dec" onclick="csNumberStep(this,-1)">−</button>
    <input type="number" class="cs-number-input__field" name="%s"
      value="%d" step="%d"%s%s data-id="%s--input" />
    <button type="button" class="cs-number-input__btn" data-id="%s--inc" onclick="csNumberStep(this,1)">+</button>
  </div>
</div>`, dataID, labelHTML, dataID, name, value, step, minAttr, maxAttr, dataID, dataID), nil
}

// ── FileUpload ──────────────────────────────────────────────────────────

func renderFileUpload(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "Drop files here or click to upload")
	hint := propStr(props, "hint", "")
	accept := propStr(props, "accept", "*")
	name := propStr(props, "name", "file")
	multiple := propBool(props, "multiple", false)
	dataID := propStr(props, "data-id", "file-upload")

	multipleAttr := ""
	if multiple {
		multipleAttr = " multiple"
	}

	hintHTML := ""
	if hint != "" {
		hintHTML = fmt.Sprintf(`<div class="cs-file-upload__hint">%s</div>`, hint)
	}

	inputID := fmt.Sprintf("fu-%s", dataID)

	return fmt.Sprintf(`<div class="cs-file-upload" data-id="%s" data-file-upload>
  <input type="file" class="cs-file-upload__input" id="%s"
    name="%s" accept="%s"%s
    onchange="csFileUploadChange(this)" />
  <label class="cs-file-upload__zone" for="%s"
    ondragover="csFileDragOver(event,this)" ondragleave="csFileDragLeave(this)" ondrop="csFileDrop(event,this,'%s')">
    <svg class="cs-file-upload__icon" viewBox="0 0 24 24" fill="currentColor" width="32" height="32">
      <path d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM14 13v4h-4v-4H7l5-5 5 5h-3z"/>
    </svg>
    <div class="cs-file-upload__text">%s</div>
    %s
  </label>
  <div class="cs-file-upload__list" data-file-list="%s"></div>
</div>`, dataID, inputID, name, accept, multipleAttr, inputID, dataID, label, hintHTML, dataID), nil
}

// ── TagInput ────────────────────────────────────────────────────────────

func renderTagInput(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	placeholder := propStr(props, "placeholder", "Add tag...")
	name := propStr(props, "name", "tags")
	tagsRaw := propStr(props, "tags", "")
	dataID := propStr(props, "data-id", "tag-input")

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-tag-input__label">%s</label>`, label)
	}

	initialTags := ""
	initialValues := ""
	if tagsRaw != "" {
		tags := strings.Split(tagsRaw, ",")
		for _, t := range tags {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			initialTags += fmt.Sprintf(`<span class="cs-tag-input__tag">%s<button type="button" class="cs-tag-input__remove" onclick="csTagRemove(this)" aria-label="Remove">×</button></span>`, t)
			if initialValues != "" {
				initialValues += ","
			}
			initialValues += t
		}
	}

	return fmt.Sprintf(`<div class="cs-tag-input-wrap" data-id="%s">
  %s
  <div class="cs-tag-input" data-tag-input="%s">
    %s
    <input type="text" class="cs-tag-input__field" placeholder="%s"
      data-id="%s--input"
      onkeydown="csTagKeydown(event,this)" />
  </div>
  <input type="hidden" name="%s" value="%s" data-tag-value="%s" />
</div>`, dataID, labelHTML, dataID, initialTags, placeholder, dataID, name, initialValues, dataID), nil
}

// ── Search ──────────────────────────────────────────────────────────────

func renderSearch(props map[string]interface{}, children string, e *Engine) (string, error) {
	placeholder := propStr(props, "placeholder", "Search...")

	return fmt.Sprintf(`<div%s><span class="cs-search__icon"><i class="bi bi-search"></i></span><input type="search" class="cs-search__input" placeholder="%s"></div>`,
		userAttrs(props, "cs-search"), placeholder), nil
}

// ── ColorInput ──────────────────────────────────────────────────────────

func renderColorInput(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", "color")
	value := propStr(props, "value", "#000000")
	dataID := propStr(props, "data-id", "color-input")

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-color-input__label">%s</label>`, label)
	}

	hexID := dataID + "--hex"

	return fmt.Sprintf(`<div class="cs-color-input-wrap" data-id="%s">
  %s
  <div class="cs-color-input">
    <input type="color" class="cs-color-input__field" name="%s" value="%s"
      data-id="%s--input"
      oninput="document.getElementById('%s').textContent=this.value" />
    <span class="cs-color-input__hex" id="%s">%s</span>
  </div>
</div>`, dataID, labelHTML, name, value, dataID, hexID, hexID, value), nil
}

// ── DateInput ───────────────────────────────────────────────────────────

func renderDateInput(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", "")
	name := propStr(props, "name", "date")
	value := propStr(props, "value", "")
	min := propStr(props, "min", "")
	max := propStr(props, "max", "")
	dataID := propStr(props, "data-id", "date-input")

	labelHTML := ""
	if label != "" {
		labelHTML = fmt.Sprintf(`<label class="cs-date-input__label" for="%s--field">%s</label>`, dataID, label)
	}

	minAttr := ""
	if min != "" {
		minAttr = fmt.Sprintf(` min="%s"`, min)
	}
	maxAttr := ""
	if max != "" {
		maxAttr = fmt.Sprintf(` max="%s"`, max)
	}
	valueAttr := ""
	if value != "" {
		valueAttr = fmt.Sprintf(` value="%s"`, value)
	}

	return fmt.Sprintf(`<div class="cs-date-input-wrap" data-id="%s">
  %s
  <input type="date" class="cs-date-input__field" id="%s--field"
    name="%s"%s%s%s data-id="%s--input" />
</div>`, dataID, labelHTML, dataID, name, valueAttr, minAttr, maxAttr, dataID), nil
}

// ── Floating Action Button ──────────────────────────────────────────────
// Props: label, icon, size (small/medium/large), color, variant (circular/extended),
//        on:click, href, disabled

func renderFab(props map[string]interface{}, children string, e *Engine) (string, error) {
	icon := propStr(props, "icon", "")
	label := propStr(props, "label", children)
	size := propStr(props, "size", "medium")
	color := propStr(props, "color", "primary")
	variant := propStr(props, "variant", "circular")
	action := propStr(props, "on:click", "")
	href := propStr(props, "href", "")
	disabled := propBool(props, "disabled", false)
	dataID := propStr(props, "data-id", "")

	classes := fmt.Sprintf("cs-fab cs-fab--%s cs-fab--%s cs-fab--%s", size, color, variant)

	disabledAttr := ""
	if disabled {
		disabledAttr = " disabled"
	}

	onclick := ""
	if action != "" {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, esc(action))
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	tag := "button"
	extra := ` type="button"`
	if href != "" {
		tag = "a"
		extra = fmt.Sprintf(` href="%s"`, esc(href))
	}

	iconHTML := bsIcon(icon, 24)
	content := iconHTML
	if variant == "extended" && label != "" {
		content = iconHTML + fmt.Sprintf(`<span class="cs-fab__label">%s</span>`, esc(label))
	}

	return fmt.Sprintf(`<%s class="%s"%s%s%s%s>%s</%s>`,
		tag, classes, extra, onclick, idAttr, disabledAttr, content, tag), nil
}

// ── Toggle Button Group ─────────────────────────────────────────────────
// Props: value, exclusive (bool), size, color, orientation (horizontal/vertical)
// Children: toggle-button atoms

func renderToggleGroup(props map[string]interface{}, children string, e *Engine) (string, error) {
	value := propStr(props, "value", "")
	size := propStr(props, "size", "medium")
	color := propStr(props, "color", "primary")
	orientation := propStr(props, "orientation", "horizontal")
	dataID := propStr(props, "data-id", "")

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<div class="cs-toggle-group cs-toggle-group--%s cs-toggle-group--size-%s cs-toggle-group--color-%s" role="group" data-value="%s"%s>%s</div>`,
		esc(orientation), esc(size), esc(color), esc(value), idAttr, children), nil
}

func renderToggleButton(props map[string]interface{}, children string, e *Engine) (string, error) {
	value := propStr(props, "value", "")
	icon := propStr(props, "icon", "")
	label := propStr(props, "label", children)
	selected := propBool(props, "selected", false)
	disabled := propBool(props, "disabled", false)

	selectedClass := ""
	if selected {
		selectedClass = " cs-toggle-btn--selected"
	}

	disabledAttr := ""
	if disabled {
		disabledAttr = " disabled"
	}

	content := ""
	if icon != "" {
		content += bsIcon(icon, 20)
	}
	if label != "" {
		content += fmt.Sprintf(`<span>%s</span>`, esc(label))
	}

	return fmt.Sprintf(`<button type="button" class="cs-toggle-btn%s" value="%s" onclick="csToggle(this)"%s>%s</button>`,
		selectedClass, esc(value), disabledAttr, content), nil
}

// ── Transfer List ───────────────────────────────────────────────────────
// Props: data-id, left-title, right-title, left-items, right-items

func renderTransferList(props map[string]interface{}, children string, e *Engine) (string, error) {
	leftTitle := propStr(props, "left-title", "Available")
	rightTitle := propStr(props, "right-title", "Selected")
	leftItems := toInterfaceSlice(props["left-items"])
	rightItems := toInterfaceSlice(props["right-items"])
	dataID := propStr(props, "data-id", "")

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	renderList := func(title string, items []interface{}, side string) string {
		var b strings.Builder
		b.WriteString(fmt.Sprintf(`<div class="cs-transfer__panel"><div class="cs-transfer__header">%s</div><div class="cs-transfer__list" data-side="%s">`, esc(title), side))
		for _, item := range items {
			m, _ := item.(map[string]interface{})
			if m == nil {
				continue
			}
			id := toString(m["id"])
			label := toString(m["label"])
			b.WriteString(fmt.Sprintf(`<label class="cs-transfer__item"><input type="checkbox" value="%s" class="cs-transfer__check"><span>%s</span></label>`, esc(id), esc(label)))
		}
		b.WriteString(`</div></div>`)
		return b.String()
	}

	controls := `<div class="cs-transfer__controls">` +
		`<button type="button" class="cs-transfer__btn" onclick="csTransfer(this,'right')" title="Move right">&#9654;</button>` +
		`<button type="button" class="cs-transfer__btn" onclick="csTransfer(this,'left')" title="Move left">&#9664;</button>` +
		`</div>`

	return fmt.Sprintf(`<div class="cs-transfer"%s>%s%s%s</div>`,
		idAttr,
		renderList(leftTitle, leftItems, "left"),
		controls,
		renderList(rightTitle, rightItems, "right"),
	), nil
}
