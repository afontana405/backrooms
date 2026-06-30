// ─── Global error capture ─────────────────────────────────────────────────────
window.addEventListener('error', function(e) {
  flightRecord('error',2,'js:error',(e.message||'')+' at '+(e.filename||'')+':'+(e.lineno||0));
});
window.addEventListener('unhandledrejection', function(e) {
  flightRecord('error',2,'js:unhandledrejection',String(e.reason||''));
});

// ─── Chat Widget ──────────────────────────────────────────────────────────────

function csChatOpen(id) {
  document.getElementById(id + '--bubble').style.display = 'none';
  document.getElementById(id + '--container').classList.add('cs-chat--open');
  document.getElementById(id + '--input').focus();
}

function csChatClose(id) {
  document.getElementById(id + '--container').classList.remove('cs-chat--open');
  document.getElementById(id + '--bubble').style.display = '';
}

function csChatSend(id) {
  var input = document.getElementById(id + '--input');
  var body  = document.getElementById(id + '--body');
  var container = document.getElementById(id + '--container');
  var msg = input.value.trim();
  if (!msg) return;

  var webhook = container.dataset.webhook;
  var route   = container.dataset.route || 'general';

  // Append user message
  var userEl = document.createElement('div');
  userEl.className = 'cs-chat-msg cs-chat-msg--user';
  userEl.textContent = msg;
  body.appendChild(userEl);
  input.value = '';

  // Typing indicator
  var typing = document.createElement('div');
  typing.className = 'cs-chat-typing';
  typing.innerHTML = '<span></span><span></span><span></span>';
  body.appendChild(typing);
  body.scrollTop = body.scrollHeight;

  // Retrieve or create session chat ID
  var chatId = sessionStorage.getItem('cs_chat_id_' + id);
  if (!chatId) {
    chatId = 'chat_' + Math.random().toString(36).substr(2, 9);
    sessionStorage.setItem('cs_chat_id_' + id, chatId);
  }

  fetch(webhook, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ chatId: chatId, message: msg, route: route })
  })
  .then(function(res) { return res.json(); })
  .then(function(data) {
    typing.remove();
    var botEl = document.createElement('div');
    botEl.className = 'cs-chat-msg cs-chat-msg--bot';
    botEl.innerHTML = data.output || 'Sorry, I couldn\'t understand that.';
    body.appendChild(botEl);
    body.scrollTop = body.scrollHeight;
  })
  .catch(function() {
    typing.remove();
    var errEl = document.createElement('div');
    errEl.className = 'cs-chat-msg cs-chat-msg--bot';
    errEl.textContent = 'Network error — please try again.';
    body.appendChild(errEl);
    body.scrollTop = body.scrollHeight;
  });
}

// ─── DataGrid ─────────────────────────────────────────────────────────────────
function csDataGridSort(th) {
  var table = th.closest('table');
  var idx = Array.from(th.parentNode.children).indexOf(th);
  var asc = th.dataset.sortDir !== 'asc';
  th.parentNode.querySelectorAll('th').forEach(function(h) {
    h.dataset.sortDir = '';
    h.classList.remove('cs-data-grid__th--asc','cs-data-grid__th--desc');
  });
  th.dataset.sortDir = asc ? 'asc' : 'desc';
  th.classList.toggle('cs-data-grid__th--asc', asc);
  th.classList.toggle('cs-data-grid__th--desc', !asc);
  var tbody = table.querySelector('tbody');
  var rows = Array.from(tbody.querySelectorAll('tr'));
  rows.sort(function(a,b) {
    var av = (a.cells[idx]||{}).textContent||'';
    var bv = (b.cells[idx]||{}).textContent||'';
    var an = parseFloat(av), bn = parseFloat(bv);
    if (!isNaN(an) && !isNaN(bn)) return asc ? an-bn : bn-an;
    return asc ? av.localeCompare(bv) : bv.localeCompare(av);
  });
  rows.forEach(function(r) { tbody.appendChild(r); });
}
function csDataGridFilter(input) {
  var q = input.value.toLowerCase();
  var grid = input.closest('.cs-data-grid');
  var rows = grid.querySelectorAll('tbody tr');
  var visible = 0;
  rows.forEach(function(row) {
    var show = row.textContent.toLowerCase().includes(q);
    row.style.display = show ? '' : 'none';
    if (show) visible++;
  });
  var empty = grid.querySelector('.cs-data-grid__empty');
  if (empty) empty.style.display = visible === 0 ? 'block' : 'none';
}

// ─── Tree ─────────────────────────────────────────────────────────────────────
function csTreeToggle(row) {
  var item = row.closest('.cs-tree-item');
  var children = item.querySelector('.cs-tree-item__children');
  if (!children) return;
  item.classList.toggle('cs-tree-item--open');
  children.style.display = item.classList.contains('cs-tree-item--open') ? '' : 'none';
}

// ─── VirtualList ──────────────────────────────────────────────────────────────
var _vlData = {};
function csVirtualListInit(id, cols, rows, rowHeight) {
  var wrap = document.querySelector('[data-id="' + id + '"]');
  if (!wrap) return;
  _vlData[id] = { cols: cols, rows: rows, rowHeight: rowHeight };
  var inner = wrap.querySelector('.cs-virtual-list__inner');
  inner.style.height = (rows.length * rowHeight) + 'px';
  csVirtualListRender(id, wrap);
}
function csVirtualListScroll(el) {
  var id = el.dataset.id;
  csVirtualListRender(id, el);
}
function csVirtualListRender(id, wrap) {
  var d = _vlData[id];
  if (!d) return;
  var scrollTop = wrap.scrollTop;
  var viewH = wrap.offsetHeight;
  var start = Math.max(0, Math.floor(scrollTop / d.rowHeight) - 2);
  var end = Math.min(d.rows.length, Math.ceil((scrollTop + viewH) / d.rowHeight) + 2);
  var inner = wrap.querySelector('.cs-virtual-list__inner');
  inner.innerHTML = '';
  for (var i = start; i < end; i++) {
    var row = d.rows[i];
    var el = document.createElement('div');
    el.className = 'cs-virtual-list__row';
    el.style.cssText = 'position:absolute;top:' + (i * d.rowHeight) + 'px;left:0;right:0;height:' + d.rowHeight + 'px';
    d.cols.forEach(function(col) {
      var cell = document.createElement('div');
      cell.className = 'cs-virtual-list__cell';
      cell.textContent = (row && row[col] !== undefined) ? row[col] : '';
      el.appendChild(cell);
    });
    inner.appendChild(el);
  }
}

// ─── Notification ─────────────────────────────────────────────────────────────
function csNotificationOpen(id) {
  var panel = document.getElementById('notification-' + id);
  if (!panel) return;
  var isOpen = panel.style.display !== 'none';
  document.querySelectorAll('[id^="notification-"]').forEach(function(p) { p.style.display = 'none'; });
  if (!isOpen) panel.style.display = 'block';
}
function csNotificationMarkAll(id) {
  var panel = document.getElementById('notification-' + id);
  if (!panel) return;
  panel.querySelectorAll('.cs-notification-item--unread').forEach(function(i) {
    i.classList.remove('cs-notification-item--unread');
  });
  var badge = document.querySelector('[data-notification="' + id + '"] .cs-notification__badge');
  if (badge) badge.remove();
}
document.addEventListener('click', function(e) {
  if (!e.target.closest('[data-notification]')) {
    document.querySelectorAll('[id^="notification-"]').forEach(function(p) { p.style.display = 'none'; });
  }
});

// ─── Command ──────────────────────────────────────────────────────────────────
function csCommandOpen(id) {
  var el = document.getElementById('command-' + id);
  if (!el) return;
  el.style.display = 'flex';
  var input = el.querySelector('.cs-command__input');
  if (input) { input.value = ''; input.focus(); }
  csCommandFilter(input, id);
}
function csCommandClose(id) {
  var el = document.getElementById('command-' + id);
  if (el) el.style.display = 'none';
}
function csCommandFilter(input, id) {
  if (!input) return;
  var q = input.value.toLowerCase();
  var items = document.querySelectorAll('#command-' + id + ' [data-command-item]');
  var first = null;
  items.forEach(function(item) {
    var show = item.textContent.toLowerCase().includes(q);
    item.style.display = show ? '' : 'none';
    item.classList.remove('cs-command__item--active');
    if (show && !first) first = item;
  });
  if (first) first.classList.add('cs-command__item--active');
}
function csCommandKey(e, id) {
  var selector = '#command-' + id + ' [data-command-item]:not([style*="none"])';
  var items = Array.from(document.querySelectorAll(selector));
  var active = document.querySelector('#command-' + id + ' .cs-command__item--active');
  var idx = items.indexOf(active);
  if (e.key === 'Escape') { csCommandClose(id); }
  else if (e.key === 'ArrowDown') {
    e.preventDefault();
    if (active) active.classList.remove('cs-command__item--active');
    var next = items[Math.min(idx + 1, items.length - 1)];
    if (next) { next.classList.add('cs-command__item--active'); next.scrollIntoView({block:'nearest'}); }
  } else if (e.key === 'ArrowUp') {
    e.preventDefault();
    if (active) active.classList.remove('cs-command__item--active');
    var prev = items[Math.max(idx - 1, 0)];
    if (prev) { prev.classList.add('cs-command__item--active'); prev.scrollIntoView({block:'nearest'}); }
  } else if (e.key === 'Enter' && active) {
    e.preventDefault();
    var action = active.dataset.action;
    if (action) csAction(action, active);
    csCommandClose(id);
  }
}
document.addEventListener('keydown', function(e) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault();
    var cmd = document.querySelector('.cs-command');
    if (cmd) csCommandOpen(cmd.id.replace('command-', ''));
  }
});

// ─── ContextMenu ──────────────────────────────────────────────────────────────
document.addEventListener('contextmenu', function(e) {
  var target = e.target.closest('[data-context-menu]');
  if (!target) return;
  e.preventDefault();
  var menuId = target.dataset.contextMenu;
  var menu = document.getElementById('ctx-' + menuId);
  if (!menu) return;
  document.querySelectorAll('.cs-context-menu').forEach(function(m) { m.style.display = 'none'; });
  menu.style.display = 'block';
  menu.style.left = e.clientX + 'px';
  menu.style.top = e.clientY + 'px';
  var rect = menu.getBoundingClientRect();
  if (rect.right > window.innerWidth) menu.style.left = (e.clientX - rect.width) + 'px';
  if (rect.bottom > window.innerHeight) menu.style.top = (e.clientY - rect.height) + 'px';
});
document.addEventListener('click', function(e) {
  if (!e.target.closest('.cs-context-menu')) {
    document.querySelectorAll('.cs-context-menu').forEach(function(m) { m.style.display = 'none'; });
  }
});

// ─── SplitView ────────────────────────────────────────────────────────────────
var _splitDrag = null;
function csSplitStart(e, id) {
  e.preventDefault();
  var wrap = document.getElementById(id);
  var pane = wrap.querySelector('.cs-split-view__pane--first');
  var isV = wrap.classList.contains('cs-split-view--vertical');
  _splitDrag = { wrap: wrap, pane: pane, isV: isV,
    startPos: isV ? e.clientY : e.clientX,
    startSize: isV ? pane.offsetHeight : pane.offsetWidth };
}
document.addEventListener('mousemove', function(e) {
  if (!_splitDrag) return;
  var s = _splitDrag;
  var delta = (s.isV ? e.clientY : e.clientX) - s.startPos;
  var total = s.isV ? s.wrap.offsetHeight : s.wrap.offsetWidth;
  var newPct = Math.max(10, Math.min(90, ((s.startSize + delta) / total) * 100));
  s.pane.style.flexBasis = newPct + '%';
});
document.addEventListener('mouseup', function() { _splitDrag = null; });
document.addEventListener('DOMContentLoaded', function() {
  document.querySelectorAll('.cs-split-view').forEach(function(wrap) {
    var panes = wrap.querySelectorAll('.cs-split-view__pane');
    if (panes.length < 2) return;
    var def = parseInt(wrap.dataset.splitDefault) || 50;
    panes[0].classList.add('cs-split-view__pane--first');
    panes[0].style.flexBasis = def + '%';
    var divider = document.createElement('div');
    divider.className = 'cs-split-view__divider';
    divider.addEventListener('mousedown', function(e) { csSplitStart(e, wrap.id); });
    wrap.insertBefore(divider, panes[1]);
  });
});

// ─── Calendar ─────────────────────────────────────────────────────────────────
var _calState = {};
function csCalendarInit(id, value) {
  var today = new Date();
  var sel = value ? new Date(value + 'T12:00:00') : null;
  _calState[id] = {
    year: sel ? sel.getFullYear() : today.getFullYear(),
    month: sel ? sel.getMonth() : today.getMonth(),
    selected: value || null, today: today
  };
  csCalendarRender(id);
}
function csCalendarNav(id, dir) {
  var s = _calState[id]; if (!s) return;
  s.month += dir;
  if (s.month < 0)  { s.month = 11; s.year--; }
  if (s.month > 11) { s.month = 0;  s.year++; }
  csCalendarRender(id);
}
function csCalendarSelect(id, dateStr) {
  var s = _calState[id]; if (!s) return;
  s.selected = dateStr;
  var hidden = document.querySelector('[data-cal-value="' + id + '"]');
  if (hidden) hidden.value = dateStr;
  var wrap = document.getElementById(id);
  if (wrap && wrap.dataset.action) csAction(wrap.dataset.action, wrap);
  csCalendarRender(id);
}
function csCalendarRender(id) {
  var s = _calState[id]; if (!s) return;
  var months = ['January','February','March','April','May','June','July','August','September','October','November','December'];
  var lbl = document.querySelector('[data-cal-label="' + id + '"]');
  if (lbl) lbl.textContent = months[s.month] + ' ' + s.year;
  var grid = document.querySelector('[data-cal-grid="' + id + '"]');
  if (!grid) return;
  var days = ['Su','Mo','Tu','We','Th','Fr','Sa'];
  var html = '<div class="cs-calendar__days-header">' + days.map(function(d){return '<span>'+d+'</span>';}).join('') + '</div>';
  html += '<div class="cs-calendar__days">';
  var firstDay = new Date(s.year, s.month, 1).getDay();
  var total = new Date(s.year, s.month + 1, 0).getDate();
  for (var i = 0; i < firstDay; i++) html += '<span class="cs-calendar__day cs-calendar__day--empty"></span>';
  for (var d = 1; d <= total; d++) {
    var ds = s.year + '-' + String(s.month+1).padStart(2,'0') + '-' + String(d).padStart(2,'0');
    var cls = 'cs-calendar__day';
    if (s.selected === ds) cls += ' cs-calendar__day--selected';
    else if (s.today.getFullYear()===s.year && s.today.getMonth()===s.month && s.today.getDate()===d) cls += ' cs-calendar__day--today';
    html += '<span class="'+cls+'" onclick="csCalendarSelect(\''+id+'\',\''+ds+'\')">'+d+'</span>';
  }
  html += '</div>';
  grid.innerHTML = html;
}

// ─── MultiSelect ──────────────────────────────────────────────────────────────
function csMultiSelectOpen(wrap) {
  var dd = wrap.querySelector('[data-ms-dropdown]');
  var isOpen = dd.style.display !== 'none';
  document.querySelectorAll('[data-ms-dropdown]').forEach(function(d) { d.style.display = 'none'; });
  document.querySelectorAll('[data-ms-wrap]').forEach(function(w) { w.classList.remove('cs-multi-select--open'); });
  if (!isOpen) {
    dd.style.display = 'block';
    wrap.classList.add('cs-multi-select--open');
  }
}
function csMultiSelectToggle(wrap, value, label) {
  var hidden = wrap.querySelector('[data-ms-value]');
  var tags = wrap.querySelector('[data-ms-tags]');
  var placeholder = wrap.querySelector('[data-ms-placeholder]');
  var vals = hidden.value ? hidden.value.split(',') : [];
  var idx = vals.indexOf(value);
  if (idx === -1) {
    vals.push(value);
    var tag = document.createElement('span');
    tag.className = 'cs-multi-select__tag';
    tag.dataset.msTag = value;
    tag.innerHTML = label + '<button type="button" onclick="event.stopPropagation();csMultiSelectRemove(this.closest(\'[data-ms-wrap]\'),\'' + value + '\')">&#215;</button>';
    tags.insertBefore(tag, placeholder);
  } else {
    vals.splice(idx, 1);
    var existing = tags.querySelector('[data-ms-tag="' + value + '"]');
    if (existing) existing.remove();
  }
  hidden.value = vals.join(',');
  placeholder.style.display = vals.length ? 'none' : '';
  wrap.querySelectorAll('[data-ms-option]').forEach(function(opt) {
    opt.classList.toggle('cs-multi-select__option--active', vals.indexOf(opt.dataset.msOption) !== -1);
  });
}
function csMultiSelectRemove(wrap, value) {
  var hidden = wrap.querySelector('[data-ms-value]');
  var tags = wrap.querySelector('[data-ms-tags]');
  var placeholder = wrap.querySelector('[data-ms-placeholder]');
  var vals = hidden.value ? hidden.value.split(',') : [];
  var idx = vals.indexOf(value);
  if (idx !== -1) vals.splice(idx, 1);
  hidden.value = vals.join(',');
  var tag = tags.querySelector('[data-ms-tag="' + value + '"]');
  if (tag) tag.remove();
  placeholder.style.display = vals.length ? 'none' : '';
  var opt = wrap.querySelector('[data-ms-option="' + value + '"]');
  if (opt) opt.classList.remove('cs-multi-select__option--active');
}
document.addEventListener('click', function(e) {
  if (!e.target.closest('[data-ms-wrap]')) {
    document.querySelectorAll('[data-ms-dropdown]').forEach(function(d) { d.style.display = 'none'; });
    document.querySelectorAll('[data-ms-wrap]').forEach(function(w) { w.classList.remove('cs-multi-select--open'); });
  }
});
