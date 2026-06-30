package engine

import "fmt"

// DiagPanelHTML returns the floating diagnostics panel injected into every page.
func DiagPanelHTML(entries string, errs, warns, infos int) string {
	badgeColor := "#34d399"
	badgeText := fmt.Sprintf("%d", infos)
	if warns > 0 {
		badgeColor = "#fbbf24"
		badgeText = fmt.Sprintf("%d", warns)
	}
	if errs > 0 {
		badgeColor = "#f87171"
		badgeText = fmt.Sprintf("%d", errs)
	}
	if errs+warns+infos == 0 {
		badgeText = "0"
	}

	return fmt.Sprintf(`
<script>window.__diag = %s;</script>
<style>%s</style>
<div id="diag-btn" class="diag-btn diag-btn--hidden" onmousedown="window.__diagDragStart(event)">
  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 013 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
  <span class="diag-btn__badge" style="background:%s">%s</span>
</div>
<div id="diag-panel" class="diag-panel diag-panel--closed">
  <div class="diag-panel__header">
    <span class="diag-panel__title">Diagnostics</span>
    <span class="diag-panel__counts">
      <span style="color:#f87171">%d err</span>
      <span style="color:#fbbf24">%d warn</span>
      <span style="color:#34d399">%d ok</span>
    </span>
    <button class="diag-panel__hdr-btn" onclick="window.__diagCopy()" title="Copy to clipboard">&#9112;</button>
    <button class="diag-panel__hdr-btn" onclick="window.__diagExport()" title="Export filtered entries to JSON">&#8615;</button>
    <button class="diag-panel__hdr-btn" onclick="window.__diagMinimize()" title="Minimize (keep launcher visible)">&ndash;</button>
    <button class="diag-panel__hdr-btn" onclick="window.__diagClose()" title="Close (hide everything)">&times;</button>
  </div>
  <div class="diag-panel__filters">
    <button class="diag-filt diag-filt--active" data-filter="all" onclick="window.__diagFilter('all',this)">All</button>
    <button class="diag-filt" data-filter="error" onclick="window.__diagFilter('error',this)">Errors</button>
    <button class="diag-filt" data-filter="warn" onclick="window.__diagFilter('warn',this)">Warnings</button>
    <button class="diag-filt" data-filter="template" onclick="window.__diagFilter('template',this)">Template</button>
    <button class="diag-filt" data-filter="render" onclick="window.__diagFilter('render',this)">Render</button>
    <button class="diag-filt" data-filter="action" onclick="window.__diagFilter('action',this)">Actions</button>
    <button class="diag-filt" data-filter="rule" onclick="window.__diagFilter('rule',this)">Rules</button>
    <button class="diag-filt" data-filter="flight" onclick="window.__diagFilter('flight',this)">Flight</button>
    <button class="diag-filt" data-filter="targeting" onclick="window.__diagFilter('targeting',this)">Targeting</button>
  </div>
  <div id="diag-entries" class="diag-panel__entries"></div>
</div>
<script>%s</script>
`, entries, diagCSS(), badgeColor, badgeText, errs, warns, infos, diagJS())
}

func diagCSS() string {
	return `
.diag-btn {
  position:fixed; bottom:16px; left:16px; z-index:99999;
  width:44px; height:44px; border-radius:50%;
  background:#1a1a2e; border:1px solid #2a2a4a;
  display:flex; align-items:center; justify-content:center;
  cursor:grab; color:#6b7280; transition:border-color .2s, color .2s;
  box-shadow:0 4px 12px rgba(0,0,0,.4);
  user-select:none; -webkit-user-select:none;
}
.diag-btn:hover { color:#e2e8f0; border-color:#00e5ff; }
.diag-btn--hidden { display:none !important; }
.diag-btn--dragging { cursor:grabbing; opacity:.8; transition:none; }
.diag-btn__badge {
  position:absolute; top:-4px; right:-4px;
  min-width:18px; height:18px; border-radius:9px;
  font-size:10px; font-weight:700; color:#0a0a0f;
  display:flex; align-items:center; justify-content:center;
  padding:0 4px; pointer-events:none;
}

.diag-panel {
  position:fixed; top:0; left:0; bottom:0; width:420px; z-index:99998;
  background:#0d0d14; border-right:1px solid #1a1a2e;
  display:flex; flex-direction:column;
  font-family:'Segoe UI',system-ui,sans-serif;
  transition:transform .25s ease;
  box-shadow:4px 0 24px rgba(0,0,0,.5);
}
.diag-panel--closed { transform:translateX(-100%); pointer-events:none; }

.diag-panel__header {
  display:flex; align-items:center; gap:8px;
  padding:14px 16px; border-bottom:1px solid #1a1a2e;
}
.diag-panel__title { font-size:15px; font-weight:700; color:#e2e8f0; }
.diag-panel__counts { font-size:12px; display:flex; gap:10px; margin-left:auto; }
.diag-panel__hdr-btn {
  background:none; border:1px solid #1a1a2e; color:#6b7280; font-size:16px;
  cursor:pointer; padding:2px 8px; line-height:1; border-radius:var(--radius-sm);
}
.diag-panel__hdr-btn:hover { color:#e2e8f0; border-color:#333; }

.diag-panel__filters {
  display:flex; gap:4px; padding:8px 16px; border-bottom:1px solid #1a1a2e;
  flex-wrap:wrap;
}
.diag-filt {
  background:#111119; border:1px solid #1a1a2e; color:#6b7280;
  padding:3px 10px; border-radius:var(--radius-sm); font-size:11px; cursor:pointer;
}
.diag-filt:hover { color:#c8d0da; border-color:#333; }
.diag-filt--active { color:#00e5ff; border-color:#00e5ff33; background:#00e5ff0d; }

.diag-panel__entries {
  flex:1; overflow-y:auto; padding:8px 0;
}
.diag-entry {
  padding:8px 16px; border-bottom:1px solid #0a0a0f;
  font-size:13px; cursor:pointer; transition:background .1s;
}
.diag-entry:hover { background:#111119; }
.diag-entry__head { display:flex; align-items:center; gap:8px; }
.diag-entry__icon { font-size:12px; flex-shrink:0; }
.diag-entry__cat {
  font-size:10px; padding:1px 6px; border-radius:3px;
  background:#1a1a2e; color:#6b7280; text-transform:uppercase;
  letter-spacing:.05em; flex-shrink:0;
}
.diag-entry__msg { color:#c8d0da; flex:1; }
.diag-entry__time { font-size:10px; color:#3b3b5c; flex-shrink:0; }
.diag-entry__detail {
  display:none; margin-top:6px; padding:6px 10px;
  background:#0a0a0f; border-radius:var(--radius-sm);
  font-size:12px; color:#6b7280; white-space:pre-wrap; line-height:1.5;
}
.diag-entry--open .diag-entry__detail { display:block; }

.diag-entry--error .diag-entry__msg { color:#f87171; }
.diag-entry--warn .diag-entry__msg { color:#fbbf24; }
.diag-entry--info .diag-entry__msg { color:#6b7280; }
`
}

func diagJS() string {
	return `
(function(){
// Live DOM queries — elements are recreated on partial navigation
function $panel() { return document.getElementById('diag-panel'); }
function $btn() { return document.getElementById('diag-btn'); }
function $entries() { return document.getElementById('diag-entries'); }
var currentFilter = 'all';
var allEntries = window.__diag || [];
var levels = ['info','warn','error'];

// ── Panel open/close ─────────────────────────────────────────────────────
window.__diagToggle = function() {
  var p = $panel(); if (p) p.classList.toggle('diag-panel--closed');
};

// Minimize: close panel, keep launcher icon visible
window.__diagMinimize = function() {
  var p = $panel(); if (p) p.classList.add('diag-panel--closed');
  if(flightPollTimer){clearInterval(flightPollTimer);flightPollTimer=null;}
};

// Close: hide everything (as if Ctrl+Shift+D not pressed)
window.__diagClose = function() {
  var p = $panel(); if (p) p.classList.add('diag-panel--closed');
  var b = $btn(); if (b) b.classList.add('diag-btn--hidden');
  if(flightPollTimer){clearInterval(flightPollTimer);flightPollTimer=null;}
};

// ── Draggable launcher ───────────────────────────────────────────────────
var dragState = { active:false, startX:0, startY:0, origX:0, origY:0, moved:false };

window.__diagDragStart = function(e) {
  if (e.button !== 0) return;
  var b = $btn(); if (!b) return;
  dragState.active = true;
  dragState.moved = false;
  dragState.startX = e.clientX;
  dragState.startY = e.clientY;
  var rect = b.getBoundingClientRect();
  dragState.origX = rect.left;
  dragState.origY = rect.top;
  b.classList.add('diag-btn--dragging');
  e.preventDefault();
};

document.addEventListener('mousemove', function(e) {
  if (!dragState.active) return;
  var b = $btn(); if (!b) return;
  var dx = e.clientX - dragState.startX;
  var dy = e.clientY - dragState.startY;
  if (Math.abs(dx) > 3 || Math.abs(dy) > 3) dragState.moved = true;
  var newX = Math.max(0, Math.min(window.innerWidth - 48, dragState.origX + dx));
  var newY = Math.max(0, Math.min(window.innerHeight - 48, dragState.origY + dy));
  b.style.left = newX + 'px';
  b.style.top = newY + 'px';
  b.style.bottom = 'auto';
  b.style.right = 'auto';
});

document.addEventListener('mouseup', function() {
  if (!dragState.active) return;
  dragState.active = false;
  var b = $btn(); if (b) b.classList.remove('diag-btn--dragging');
  // If not dragged, treat as click → toggle panel
  if (!dragState.moved) window.__diagToggle();
});

// ── Filters ──────────────────────────────────────────────────────────────
window.__diagFilter = function(f, filterBtn) {
  currentFilter = f;
  document.querySelectorAll('.diag-filt').forEach(function(b){ b.classList.remove('diag-filt--active'); });
  if (filterBtn) filterBtn.classList.add('diag-filt--active');
  if (f === 'flight') {
    window.__diagFetchFlight();
    if (!flightPollTimer) {
      flightPollTimer = setInterval(function() {
        var p = $panel();
        if (p && !p.classList.contains('diag-panel--closed') && currentFilter === 'flight') {
          window.__diagFetchFlight();
        }
      }, 2500);
    }
  } else if (f === 'targeting') {
    if (flightPollTimer) { clearInterval(flightPollTimer); flightPollTimer = null; }
    if (window.__tgtRenderDiag) window.__tgtRenderDiag();
  } else {
    if (flightPollTimer) { clearInterval(flightPollTimer); flightPollTimer = null; }
    renderEntries();
  }
};

// ── Add entry (client-side) ──────────────────────────────────────────────
window.__diagAdd = function(entry) {
  allEntries.push(entry);
  renderEntries();
  updateBadge();
};

// ── Copy filtered entries to clipboard ───────────────────────────────────
window.__diagCopy = function() {
  var data;
  if (currentFilter === 'flight') {
    data = flightServerCache.concat((window.__flight||[]).slice());
    data.sort(function(a,b){return(a.time||'')<(b.time||'')?-1:(a.time||'')>(b.time||'')?1:0;});
  } else {
    data = getFilteredEntries();
  }
  var json = JSON.stringify(data, null, 2);
  navigator.clipboard.writeText(json).then(function(){
    var btn = document.querySelector('[onclick="window.__diagCopy()"]');
    if(btn){var orig=btn.textContent;btn.textContent='Copied';setTimeout(function(){btn.textContent=orig;},1500);}
  });
};

// ── Export filtered entries to JSON ──────────────────────────────────────
window.__diagExport = function() {
  var data;
  if (currentFilter === 'flight') {
    data = flightServerCache.concat((window.__flight||[]).slice());
    data.sort(function(a,b){return(a.time||'')<(b.time||'')?-1:(a.time||'')>(b.time||'')?1:0;});
  } else {
    data = getFilteredEntries();
  }
  var json = JSON.stringify(data, null, 2);
  var blob = new Blob([json], { type: 'application/json' });
  var url = URL.createObjectURL(blob);
  var a = document.createElement('a');
  a.href = url;
  a.download = 'diag-' + currentFilter + '-' + new Date().toISOString().slice(0,19).replace(/:/g,'-') + '.json';
  a.click();
  URL.revokeObjectURL(url);
};

function getFilteredEntries() {
  if (currentFilter === 'all') return allEntries;
  return allEntries.filter(function(e) {
    if (currentFilter === 'error') return e.level === 2;
    if (currentFilter === 'warn') return e.level === 1;
    return e.cat === currentFilter;
  });
}

// ── Render ────────────────────────────────────────────────────────────────
function renderEntries() {
  var filtered = getFilteredEntries();
  var html = '';
  filtered.forEach(function(e) {
    var lvl = levels[e.level] || 'info';
    var icon = e.level === 2 ? '&#10007;' : e.level === 1 ? '&#9888;' : '&#10003;';
    var detail = e.detail ? '<div class="diag-entry__detail">' + esc(e.detail) + '</div>' : '';
    html += '<div class="diag-entry diag-entry--' + lvl + '" onclick="this.classList.toggle(\'diag-entry--open\')">'
      + '<div class="diag-entry__head">'
      + '<span class="diag-entry__icon">' + icon + '</span>'
      + '<span class="diag-entry__cat">' + esc(e.cat) + '</span>'
      + '<span class="diag-entry__msg">' + esc(e.msg) + '</span>'
      + '<span class="diag-entry__time">' + esc(e.time || '') + '</span>'
      + '</div>' + detail + '</div>';
  });
  if (!html) html = '<div style="padding:24px;color:#3b3b5c;text-align:center">No entries match filter</div>';
  var el = $entries(); if (el) el.innerHTML = html;
}

// ─── Flight Recorder integration ──────────────────────────────────────────
var flightLastServerId = 0;
var flightServerCache = [];
var flightPollTimer = null;

window.__diagFetchFlight = function() {
  var clientEvents = (window.__flight || []).slice();
  fetch('/api/flight/snapshot', {
    method:'POST', headers:{'Content-Type':'application/json'},
    credentials:'same-origin',
    body:JSON.stringify({afterId:flightLastServerId})
  })
  .then(function(res){return res.json();})
  .then(function(result){
    var srvEvts = (result&&result.data&&result.data.events)||[];
    if(srvEvts.length>0) flightLastServerId=srvEvts[srvEvts.length-1].id;
    srvEvts=srvEvts.filter(function(e){return e.msg.indexOf('flight/snapshot')===-1&&e.msg.indexOf('flight/push')===-1;});
    flightServerCache=flightServerCache.concat(srvEvts);
    if(flightServerCache.length>500) flightServerCache=flightServerCache.slice(flightServerCache.length-500);
    var all=flightServerCache.concat(clientEvents);
    all.sort(function(a,b){return(a.time||'')<(b.time||'')?-1:(a.time||'')>(b.time||'')?1:0;});
    if(all.length>500) all=all.slice(all.length-500);
    renderFlightEntries(all);
  })
  .catch(function(){});
};

function renderFlightEntries(events) {
  var html='';
  events.forEach(function(e){
    var lvl=['info','warn','error'][e.level]||'info';
    var icon=e.level===2?'&#10007;':e.level===1?'&#9888;':'&#10003;';
    var src=e.src==='server'?'SRV':'CLI';
    var t=e.time||'';
    if(t.indexOf('T')>-1)t=t.split('T')[1].split('.')[0];
    var detail=e.detail?'<div class="diag-entry__detail">'+esc(e.detail)+'</div>':'';
    html+='<div class="diag-entry diag-entry--'+lvl+'" onclick="this.classList.toggle(\'diag-entry--open\')">'
      +'<div class="diag-entry__head">'
      +'<span class="diag-entry__icon">'+icon+'</span>'
      +'<span class="diag-entry__cat">'+src+'</span>'
      +'<span class="diag-entry__cat">'+esc(e.cat)+'</span>'
      +'<span class="diag-entry__msg">'+esc(e.msg)+'</span>'
      +'<span class="diag-entry__time">'+esc(t)+'</span>'
      +'</div>'+detail+'</div>';
  });
  if(!html) html='<div style="padding:24px;color:#3b3b5c;text-align:center">No flight events recorded</div>';
  var el=$entries();if(el)el.innerHTML=html;
}

function updateBadge() {
  var errs=0, warns=0;
  allEntries.forEach(function(e){ if(e.level===2)errs++; if(e.level===1)warns++; });
  var badge = document.querySelector('.diag-btn__badge');
  if (!badge) return;
  if (errs > 0) { badge.style.background='#f87171'; badge.textContent=errs; }
  else if (warns > 0) { badge.style.background='#fbbf24'; badge.textContent=warns; }
  else { badge.style.background='#34d399'; badge.textContent=allEntries.length; }
}

function esc(s) { if(!s)return''; return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }

// ── Intercept csPost for action tracing ──────────────────────────────────
if (window.csPost) {
  var origPost = window.csPost;
  window.csPost = function(action, data, opts) {
    var t = new Date().toLocaleTimeString();
    window.__diagAdd({level:0,cat:'action',msg:'POST /api/'+action, detail:'Payload: '+JSON.stringify(data), time:t});

    var origFetch = window.fetch;
    window.fetch = function(url, fetchOpts) {
      return origFetch.apply(this, arguments).then(function(res) {
        var clone = res.clone();
        clone.json().then(function(result) {
          if (result.error) {
            window.__diagAdd({level:2,cat:'action',msg:'Action error: '+action, detail:result.error, time:new Date().toLocaleTimeString()});
          } else if (result.redirect) {
            window.__diagAdd({level:0,cat:'action',msg:'Action redirect: '+result.redirect, time:new Date().toLocaleTimeString()});
          } else if (result.toast) {
            window.__diagAdd({level:0,cat:'action',msg:'Action toast: '+result.toast, time:new Date().toLocaleTimeString()});
          }
        }).catch(function(){});
        return res;
      });
    };
    origPost(action, data, opts);
    setTimeout(function(){ window.fetch = origFetch; }, 100);
  };
}

// ── Double-tap Ctrl: toggle diagnostics visibility ──────────────────────
var _diagLastCtrl = 0;
document.addEventListener('keyup', function(e) {
  if (e.key !== 'Control') return;
  var now = Date.now();
  if (now - _diagLastCtrl < 400) {
    _diagLastCtrl = 0;
    var b = $btn(); var p = $panel();
    if (!b || !p) return;
    if (b.classList.contains('diag-btn--hidden')) {
      b.classList.remove('diag-btn--hidden');
      p.classList.remove('diag-panel--closed');
    } else if (!p.classList.contains('diag-panel--closed')) {
      window.__diagClose();
    } else {
      p.classList.remove('diag-panel--closed');
    }
  } else {
    _diagLastCtrl = now;
  }
});

// Initial render
renderEntries();
})();
`
}
