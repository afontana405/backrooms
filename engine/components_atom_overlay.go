package engine

import "fmt"

// ── Menu ────────────────────────────────────────────────────────────────

func renderMenu(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "menu")
	label := propStr(props, "label", "Actions")
	variant := propVariant(props, "outline")
	size := propSize(props, "md")
	dataID := propStr(props, "data-id", "menu--"+id)

	cls := fmt.Sprintf("cs-button cs-button--%s cs-button--%s", variant, size)

	return fmt.Sprintf(`<div class="cs-menu" data-menu="%s" data-id="%s">
  <button class="%s" data-menu-trigger="%s" data-id="%s--trigger" type="button">
    %s <span class="cs-menu__arrow">▾</span>
  </button>
  <div class="cs-menu__dropdown" data-menu-dropdown="%s">%s</div>
</div>`, id, dataID, cls, id, dataID, label, id, children), nil
}

func renderMenuItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	action := propStr(props, "on:click", "")
	icon := propStr(props, "icon", "")
	color := propStr(props, "color", "")
	disabled := propBool(props, "disabled", false)
	dataID := propStr(props, "data-id", "menu-item")

	iconHTML := ""
	if icon != "" {
		iconHTML, _ = renderIcon(map[string]interface{}{"name": icon, "size": float64(14)}, "", e)
	}

	cls := "cs-menu__item"
	if color != "" {
		cls += fmt.Sprintf(" cs-menu__item--color-%s", color)
	}
	if disabled {
		cls += " cs-menu__item--disabled"
	}

	onclick := ""
	if action != "" && !disabled {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, action)
	}

	disabledAttr := ""
	if disabled {
		disabledAttr = " disabled"
	}

	return fmt.Sprintf(`<button class="%s" type="button" data-id="%s"%s%s>%s%s</button>`,
		cls, dataID, onclick, disabledAttr, iconHTML, label), nil
}

// ── Popover ─────────────────────────────────────────────────────────────

func renderPopover(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "pop")
	label := propStr(props, "label", "More info")
	variant := propVariant(props, "outline")
	size := propSize(props, "md")
	placement := propStr(props, "placement", "bottom")
	dataID := propStr(props, "data-id", "popover--"+id)

	cls := fmt.Sprintf("cs-button cs-button--%s cs-button--%s", variant, size)

	return fmt.Sprintf(`<div class="cs-popover" data-popover="%s" data-id="%s">
  <button class="%s" data-popover-trigger="%s" data-id="%s--trigger" type="button">%s</button>
  <div class="cs-popover__panel cs-popover__panel--%s" data-popover-panel="%s">
    <div class="cs-popover__inner">%s</div>
  </div>
</div>`, id, dataID, cls, id, dataID, label, placement, id, children), nil
}

// ── Drawer ──────────────────────────────────────────────────────────────
// ["drawer", { "id": "settings", "title": "Settings", "side": "right" }, ...children]

func renderDrawer(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "drawer")
	title := propStr(props, "title", "")
	side := propStr(props, "side", "right")
	dataID := propStr(props, "data-id", "drawer--"+id)

	titleHTML := ""
	if title != "" {
		titleHTML = fmt.Sprintf(`<div class="cs-drawer__header">
    <h2 class="cs-drawer__title">%s</h2>
    <button class="cs-drawer__close" data-drawer-close="%s" aria-label="Close">&#10005;</button>
  </div>`, title, id)
	}

	return fmt.Sprintf(`<div class="cs-drawer cs-drawer--%s" id="drawer-%s" data-id="%s">
  <div class="cs-drawer__backdrop"></div>
  <div class="cs-drawer__panel">
    %s
    <div class="cs-drawer__body">%s</div>
  </div>
</div>`, side, id, dataID, titleHTML, children), nil
}

// ── Snackbar ────────────────────────────────────────────────────────────
// ["snackbar", { "position": "bottom-right" }]

func renderSnackbar(props map[string]interface{}, children string, e *Engine) (string, error) {
	position := propStr(props, "position", "bottom-right")
	dataID := propStr(props, "data-id", "snackbar")

	return fmt.Sprintf(`<div id="cs-snackbar" class="cs-snackbar cs-snackbar--%s" data-id="%s" aria-live="polite"></div>`,
		position, dataID), nil
}

// ── Confirm ─────────────────────────────────────────────────────────────
// ["confirm", { "id": "del", "title": "Delete item?", "on:confirm": "items/delete" }]

func renderConfirm(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "confirm")
	title := propStr(props, "title", "Are you sure?")
	message := propStr(props, "message", "")
	confirmLabel := propStr(props, "confirm-label", "Confirm")
	cancelLabel := propStr(props, "cancel-label", "Cancel")
	action := propStr(props, "on:confirm", "")
	variant := propStr(props, "variant", "danger")
	dataID := propStr(props, "data-id", "confirm--"+id)

	messageHTML := ""
	if message != "" {
		messageHTML = fmt.Sprintf(`<p class="cs-confirm__message">%s</p>`, message)
	}

	confirmOnclick := fmt.Sprintf(`csModal('%s','close')`, id)
	if action != "" {
		confirmOnclick = fmt.Sprintf(`csAction('%s',this);csModal('%s','close')`, action, id)
	}

	confirmCls := fmt.Sprintf("cs-button cs-button--solid cs-button--md cs-confirm__btn--confirm cs-confirm__btn--%s", variant)
	cancelCls := "cs-button cs-button--ghost cs-button--md cs-confirm__btn--cancel"

	return fmt.Sprintf(`<div class="cs-modal" id="modal-%s" data-id="%s">
  <div class="cs-modal__backdrop"></div>
  <div class="cs-modal__dialog cs-modal__dialog--sm">
    <div class="cs-confirm__header">
      <h2 class="cs-modal__title">%s</h2>
    </div>
    <div class="cs-modal__body">
      %s
      <div class="cs-confirm__actions">
        <button class="%s" type="button" data-id="%s--cancel"
          onclick="csModal('%s','close')">%s</button>
        <button class="%s" type="button" data-id="%s--confirm"
          onclick="%s">%s</button>
      </div>
    </div>
  </div>
</div>`, id, dataID,
		title,
		messageHTML,
		cancelCls, dataID, id, cancelLabel,
		confirmCls, dataID, confirmOnclick, confirmLabel), nil
}

// ── Notification ────────────────────────────────────────────────────────

func renderNotification(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "notif")
	count := int(propFloat(props, "count", 0))
	dataID := propStr(props, "data-id", "notification--"+id)

	badgeHTML := ""
	if count > 0 {
		badgeHTML = fmt.Sprintf(`<span class="cs-notification__badge" data-id="%s--badge">%d</span>`, dataID, count)
	}

	return fmt.Sprintf(`<div class="cs-notification" data-notification="%s" data-id="%s">
  <button class="cs-notification__bell" type="button" onclick="csNotificationOpen('%s')" data-id="%s--bell"
    aria-label="Notifications">&#128276;%s</button>
  <div class="cs-notification__panel" id="notification-%s" style="display:none">
    <div class="cs-notification__header">
      <span class="cs-notification__title">Notifications</span>
      <button class="cs-notification__mark-all" type="button" onclick="csNotificationMarkAll('%s')" data-id="%s--mark-all">Mark all read</button>
    </div>
    <div class="cs-notification__list">%s</div>
  </div>
</div>`, id, dataID, id, dataID, badgeHTML, id, id, dataID, children), nil
}

func renderNotificationItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	title := propStr(props, "title", "")
	body := propStr(props, "body", children)
	time := propStr(props, "time", "")
	unread := propBool(props, "unread", false)
	action := propStr(props, "on:click", "")
	dataID := propStr(props, "data-id", "notification-item")

	cls := "cs-notification-item"
	if unread {
		cls += " cs-notification-item--unread"
	}

	onclick := ""
	if action != "" {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, action)
	}

	timeHTML := ""
	if time != "" {
		timeHTML = fmt.Sprintf(`<span class="cs-notification-item__time">%s</span>`, time)
	}

	return fmt.Sprintf(`<div class="%s" data-id="%s"%s>
  <span class="cs-notification-item__dot"></span>
  <div class="cs-notification-item__content">
    <div class="cs-notification-item__title">%s</div>
    <div class="cs-notification-item__body">%s</div>
    %s
  </div>
</div>`, cls, dataID, onclick, title, body, timeHTML), nil
}

// ── Command ─────────────────────────────────────────────────────────────

func renderCommand(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "command")
	placeholder := propStr(props, "placeholder", "Type a command or search...")
	dataID := propStr(props, "data-id", "command--"+id)

	var itemsHTML string
	if items, ok := props["items"]; ok {
		if itemList, ok := items.([]interface{}); ok {
			for _, item := range itemList {
				if m, ok := item.(map[string]interface{}); ok {
					label := fmt.Sprintf("%v", m["label"])
					action := fmt.Sprintf("%v", m["action"])
					desc := ""
					if d, ok := m["description"]; ok {
						desc = fmt.Sprintf(`<span class="cs-command__item-desc">%v</span>`, d)
					}
					iconHTML := ""
					if ic, ok := m["icon"]; ok {
						iconHTML, _ = renderIcon(map[string]interface{}{"name": fmt.Sprintf("%v", ic), "size": float64(14)}, "", e)
					}
					itemsHTML += fmt.Sprintf(
						`<div class="cs-command__item" data-command-item data-action="%s" data-id="%s--item"
  onclick="csAction('%s',this);csCommandClose('%s')">
  <span class="cs-command__item-icon">%s</span>
  <div class="cs-command__item-body"><span class="cs-command__item-label">%s</span>%s</div>
</div>`, action, dataID, action, id, iconHTML, label, desc)
				}
			}
		}
	}
	itemsHTML += children

	return fmt.Sprintf(`<div class="cs-command" id="command-%s" data-id="%s" style="display:none" role="dialog" aria-modal="true">
  <div class="cs-command__backdrop" onclick="csCommandClose('%s')"></div>
  <div class="cs-command__dialog">
    <div class="cs-command__search-wrap">
      <span class="cs-command__search-icon">&#8984;</span>
      <input class="cs-command__input" type="text" placeholder="%s"
        oninput="csCommandFilter(this,'%s')"
        onkeydown="csCommandKey(event,'%s')"
        data-id="%s--input" autocomplete="off" />
    </div>
    <div class="cs-command__results" data-command-results="%s">%s</div>
  </div>
</div>`, id, dataID, id, placeholder, id, id, dataID, id, itemsHTML), nil
}

func renderCommandItem(props map[string]interface{}, children string, e *Engine) (string, error) {
	label := propStr(props, "label", children)
	action := propStr(props, "on:click", "")
	desc := propStr(props, "description", "")
	icon := propStr(props, "icon", "")
	dataID := propStr(props, "data-id", "command-item")

	iconHTML := ""
	if icon != "" {
		iconHTML, _ = renderIcon(map[string]interface{}{"name": icon, "size": float64(14)}, "", e)
	}

	descHTML := ""
	if desc != "" {
		descHTML = fmt.Sprintf(`<span class="cs-command__item-desc">%s</span>`, desc)
	}

	onclick := fmt.Sprintf(`csAction('%s',this);this.closest('.cs-command').style.display='none'`, action)

	return fmt.Sprintf(`<div class="cs-command__item" data-command-item data-action="%s" data-id="%s"
  onclick="%s">
  <span class="cs-command__item-icon">%s</span>
  <div class="cs-command__item-body"><span class="cs-command__item-label">%s</span>%s</div>
</div>`, action, dataID, onclick, iconHTML, label, descHTML), nil
}

// ── ContextMenu ─────────────────────────────────────────────────────────

func renderContextMenu(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "ctx")
	dataID := propStr(props, "data-id", "context-menu--"+id)

	return fmt.Sprintf(`<div class="cs-context-menu" id="ctx-%s" data-id="%s" style="display:none;position:fixed;z-index:9000">
  <div class="cs-context-menu__inner">%s</div>
</div>`, id, dataID, children), nil
}

// ── HoverCard ───────────────────────────────────────────────────────────

func renderHoverCard(props map[string]interface{}, children string, e *Engine) (string, error) {
	placement := propStr(props, "placement", "top")
	dataID := propStr(props, "data-id", "hover-card")

	return fmt.Sprintf(`<span class="cs-hover-card cs-hover-card--%s" data-id="%s">%s</span>`,
		placement, dataID, children), nil
}

// ── Backdrop ────────────────────────────────────────────────────────────

func renderBackdrop(props map[string]interface{}, children string, e *Engine) (string, error) {
	open, _ := props["open"].(bool)
	invisible, _ := props["invisible"].(bool)
	action := propStr(props, "on:click", "")
	dataID := propStr(props, "data-id", "")

	classes := "cs-backdrop"
	if open {
		classes += " cs-backdrop--open"
	}
	if invisible {
		classes += " cs-backdrop--invisible"
	}

	onclick := ""
	if action != "" {
		onclick = fmt.Sprintf(` onclick="csAction('%s',this)"`, esc(action))
	}

	idAttr := ""
	if dataID != "" {
		idAttr = fmt.Sprintf(` data-id="%s"`, esc(dataID))
	}

	return fmt.Sprintf(`<div class="%s"%s%s>%s</div>`, classes, onclick, idAttr, children), nil
}

// ── Dialog ──────────────────────────────────────────────────────────────

func renderDialog(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "")
	fullScreen, _ := props["fullScreen"].(bool)
	maxWidth := propStr(props, "maxWidth", "sm")
	fullWidth, _ := props["fullWidth"].(bool)

	classes := "cs-dialog"
	if fullScreen {
		classes += " cs-dialog--fullscreen"
	}
	if fullWidth {
		classes += " cs-dialog--fullwidth"
	}

	paperClasses := fmt.Sprintf("cs-dialog__paper cs-dialog__paper--%s", esc(maxWidth))

	idAttr := ""
	if id != "" {
		idAttr = fmt.Sprintf(` id="%s" data-id="%s"`, esc(id), esc(id))
	}

	return fmt.Sprintf(`<div class="%s" role="dialog"%s>`+
		`<div class="cs-backdrop cs-dialog__backdrop" onclick="csAction('close-modal:%s',this)"></div>`+
		`<div class="%s">%s</div>`+
		`</div>`,
		classes, idAttr, esc(id), paperClasses, children), nil
}

func renderDialogTitle(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<div class="cs-dialog__title">%s</div>`, children), nil
}

func renderDialogContent(props map[string]interface{}, children string, e *Engine) (string, error) {
	dividers, _ := props["dividers"].(bool)
	classes := "cs-dialog__content"
	if dividers {
		classes += " cs-dialog__content--dividers"
	}
	return fmt.Sprintf(`<div class="%s">%s</div>`, classes, children), nil
}

func renderDialogActions(props map[string]interface{}, children string, e *Engine) (string, error) {
	return fmt.Sprintf(`<div class="cs-dialog__actions">%s</div>`, children), nil
}
