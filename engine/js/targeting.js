// ── Element Targeting Inspector (F10) ────────────────────────────────────────
// Pesticide-style element inspector. F10 toggles on/off.
// Hover highlights elements with outline + info tooltip.
// Click captures element details to flight recorder + diagnostics Targeting tab.

var _tgtActive = false;
var _tgtOverlay = null;
var _tgtTooltip = null;
var _tgtCurrent = null;
var _tgtLog = [];

function tgtBuildSelector(el) {
  if (el.id) return '#' + el.id;
  var parts = [];
  while (el && el !== document.body && el !== document.documentElement) {
    var s = el.tagName.toLowerCase();
    if (el.id) { parts.unshift('#' + el.id); break; }
    if (el.getAttribute('data-id')) { parts.unshift('[data-id="' + el.getAttribute('data-id') + '"]'); break; }
    if (el.className && typeof el.className === 'string') {
      var cls = el.className.trim().split(/\s+/).filter(function(c) { return c.indexOf('tgt-') !== 0; }).slice(0, 3);
      if (cls.length) s += '.' + cls.join('.');
    }
    parts.unshift(s);
    el = el.parentElement;
  }
  return parts.join(' > ');
}

function tgtGetInfo(el) {
  var rect = el.getBoundingClientRect();
  var cs = window.getComputedStyle(el);
  return {
    tag: el.tagName.toLowerCase(),
    id: el.id || '',
    dataId: el.getAttribute('data-id') || '',
    classes: el.className && typeof el.className === 'string' ? el.className.trim().split(/\s+/).filter(function(c) { return c.indexOf('tgt-') !== 0; }) : [],
    selector: tgtBuildSelector(el),
    x: Math.round(rect.left),
    y: Math.round(rect.top),
    w: Math.round(rect.width),
    h: Math.round(rect.height),
    display: cs.display,
    position: cs.position,
    padding: cs.padding,
    margin: cs.margin,
    fontSize: cs.fontSize,
    color: cs.color,
    background: cs.backgroundColor,
    flex: cs.flex !== 'none' ? cs.flex : '',
    flexDirection: cs.flexDirection !== 'row' ? cs.flexDirection : '',
    alignItems: cs.alignItems !== 'normal' ? cs.alignItems : '',
    justifyContent: cs.justifyContent !== 'normal' ? cs.justifyContent : '',
    gap: cs.gap !== 'normal' ? cs.gap : '',
    overflow: cs.overflow !== 'visible' ? cs.overflow : ''
  };
}

function tgtCreateOverlay() {
  if (_tgtOverlay) return;
  _tgtOverlay = document.createElement('div');
  _tgtOverlay.className = 'tgt-overlay';
  document.body.appendChild(_tgtOverlay);

  _tgtTooltip = document.createElement('div');
  _tgtTooltip.className = 'tgt-tooltip';
  document.body.appendChild(_tgtTooltip);
}

function tgtRemoveOverlay() {
  if (_tgtOverlay) { _tgtOverlay.remove(); _tgtOverlay = null; }
  if (_tgtTooltip) { _tgtTooltip.remove(); _tgtTooltip = null; }
}

function tgtUpdateOverlay(el) {
  if (!_tgtOverlay || !el) return;
  var rect = el.getBoundingClientRect();
  _tgtOverlay.style.left = rect.left + 'px';
  _tgtOverlay.style.top = rect.top + 'px';
  _tgtOverlay.style.width = rect.width + 'px';
  _tgtOverlay.style.height = rect.height + 'px';
  _tgtOverlay.style.display = 'block';

  var tag = el.tagName.toLowerCase();
  var id = el.id ? '#' + el.id : '';
  var did = el.getAttribute('data-id') ? ' [' + el.getAttribute('data-id') + ']' : '';
  var cls = el.className && typeof el.className === 'string' ? '.' + el.className.trim().split(/\s+/).filter(function(c) { return c.indexOf('tgt-') !== 0; }).slice(0, 2).join('.') : '';
  var dims = Math.round(rect.width) + 'x' + Math.round(rect.height);
  _tgtTooltip.textContent = tag + id + cls + did + '  ' + dims;
  _tgtTooltip.style.display = 'block';

  var tx = rect.left;
  var ty = rect.top - 28;
  if (ty < 0) ty = rect.bottom + 4;
  if (tx + 300 > window.innerWidth) tx = window.innerWidth - 310;
  _tgtTooltip.style.left = Math.max(0, tx) + 'px';
  _tgtTooltip.style.top = ty + 'px';
}

function tgtHideOverlay() {
  if (_tgtOverlay) _tgtOverlay.style.display = 'none';
  if (_tgtTooltip) _tgtTooltip.style.display = 'none';
}

function tgtOnMove(e) {
  var el = document.elementFromPoint(e.clientX, e.clientY);
  if (!el || el === _tgtOverlay || el === _tgtTooltip) return;
  // Ignore diag panel elements
  var dp = document.getElementById('diag-panel');
  var db = document.getElementById('diag-btn');
  if (dp && dp.contains(el)) return;
  if (db && db.contains(el)) return;
  _tgtCurrent = el;
  tgtUpdateOverlay(el);
}

function tgtOnClick(e) {
  if (!_tgtActive || !_tgtCurrent) return;
  // Ignore diag panel clicks
  var dp = document.getElementById('diag-panel');
  var db = document.getElementById('diag-btn');
  if (dp && dp.contains(e.target)) return;
  if (db && db.contains(e.target)) return;
  e.preventDefault();
  e.stopPropagation();
  e.stopImmediatePropagation();

  var info = tgtGetInfo(_tgtCurrent);
  _tgtLog.unshift(info);
  if (_tgtLog.length > 50) _tgtLog.pop();

  flightRecord('targeting', 0, 'target:click', JSON.stringify(info));

  // Auto-render if targeting tab is active in diag panel
  if (window.__tgtRenderDiag) window.__tgtRenderDiag();
}

// ── Render targeting entries into diag panel ─────────────────────────────────
window.__tgtRenderDiag = function() {
  var el = document.getElementById('diag-entries');
  if (!el) return;
  var html = '';
  for (var i = 0; i < Math.min(_tgtLog.length, 20); i++) {
    var info = _tgtLog[i];
    var fresh = i === 0 ? ' tgt-entry--fresh' : '';
    html += '<div class="diag-entry tgt-entry' + fresh + '">';
    html += '<div class="tgt-entry-sel">' + escH(info.selector) + '</div>';
    if (info.dataId) html += '<div class="tgt-entry-did">data-id: ' + escH(info.dataId) + '</div>';
    html += '<div class="tgt-entry-dims">' + info.w + 'x' + info.h + ' @ ' + info.x + ',' + info.y + '</div>';
    html += '<div class="tgt-entry-props">';
    html += '<span>display: ' + info.display + '</span>';
    if (info.position !== 'static') html += '<span>position: ' + info.position + '</span>';
    html += '<span>padding: ' + info.padding + '</span>';
    html += '<span>margin: ' + info.margin + '</span>';
    if (info.flex) html += '<span>flex: ' + info.flex + '</span>';
    if (info.flexDirection) html += '<span>flex-direction: ' + info.flexDirection + '</span>';
    if (info.alignItems) html += '<span>align-items: ' + info.alignItems + '</span>';
    if (info.justifyContent) html += '<span>justify-content: ' + info.justifyContent + '</span>';
    if (info.gap) html += '<span>gap: ' + info.gap + '</span>';
    if (info.overflow) html += '<span>overflow: ' + info.overflow + '</span>';
    html += '<span>font-size: ' + info.fontSize + '</span>';
    html += '<span>bg: ' + info.background + '</span>';
    html += '</div>';
    html += '</div>';
  }
  if (!html) html = '<div style="padding:24px;color:#3b3b5c;text-align:center">F10 to activate targeting, then click elements to inspect</div>';
  el.innerHTML = html;
};

window.tgtToggle = function() {
  _tgtActive = !_tgtActive;
  if (_tgtActive) {
    tgtCreateOverlay();
    document.addEventListener('mousemove', tgtOnMove, true);
    document.addEventListener('click', tgtOnClick, true);
    document.body.classList.add('tgt-active');
    flightRecord('targeting', 0, 'target:on', '');
    // Open diag panel to targeting tab
    var b = document.getElementById('diag-btn');
    var p = document.getElementById('diag-panel');
    if (b) b.classList.remove('diag-btn--hidden');
    if (p) p.classList.remove('diag-panel--closed');
    var tBtn = document.querySelector('[data-filter="targeting"]');
    if (tBtn) window.__diagFilter('targeting', tBtn);
  } else {
    document.removeEventListener('mousemove', tgtOnMove, true);
    document.removeEventListener('click', tgtOnClick, true);
    tgtHideOverlay();
    tgtRemoveOverlay();
    document.body.classList.remove('tgt-active');
    _tgtCurrent = null;
    flightRecord('targeting', 0, 'target:off', '');
  }
};

// ── Targeting CSS (injected once) ────────────────────────────────────────────
(function() {
  var s = document.createElement('style');
  s.textContent = [
    '.tgt-overlay { position:fixed; z-index:99998; pointer-events:none; border:2px solid #06b6d4; background:rgba(6,182,212,0.08); display:none; transition:all 0.05s; }',
    '.tgt-tooltip { position:fixed; z-index:99999; pointer-events:none; background:#0a0e17; color:#06b6d4; font-family:monospace; font-size:12px; padding:3px 8px; border:1px solid #06b6d4; border-radius:3px; white-space:nowrap; display:none; }',
    '.tgt-entry { padding:8px 16px; }',
    '.tgt-entry--fresh { border-left:2px solid #06b6d4; }',
    '.tgt-entry-sel { color:#06b6d4; font-weight:700; margin-bottom:4px; word-break:break-all; font-size:12px; }',
    '.tgt-entry-did { color:#f59e0b; font-size:11px; margin-bottom:4px; }',
    '.tgt-entry-dims { color:#f59e0b; margin-bottom:4px; font-size:12px; }',
    '.tgt-entry-props { display:flex; flex-wrap:wrap; gap:2px 10px; font-size:11px; }',
    '.tgt-entry-props span { color:#94a3b8; }',
    '.tgt-active * { cursor:crosshair !important; }'
  ].join('\n');
  document.head.appendChild(s);
})();
