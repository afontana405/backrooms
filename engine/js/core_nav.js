// ─── data-navigate — wire select change to partial navigation ────────────────
function csWireNavigateSelect() {
  var els = document.querySelectorAll('select[data-navigate]');
  for (var i = 0; i < els.length; i++) {
    (function(el) {
      if (el._csNavWired) return;
      el._csNavWired = true;
      el.addEventListener('change', function() {
        var pattern = el.dataset.navigate;
        var url = pattern.replace('{value}', encodeURIComponent(el.value));
        if (url.indexOf('/page/') === 0) window.csNavigate(url);
        else window.location.href = url;
      });
    })(els[i]);
  }
}
document.addEventListener('DOMContentLoaded', csWireNavigateSelect);
new MutationObserver(csWireNavigateSelect).observe(document.body || document.documentElement, { childList: true, subtree: true });

// ─── data-linenums — sync line numbers with textarea ─────────────────────────
function csWireLineNums() {
  var areas = document.querySelectorAll('textarea[data-id]');
  for (var i = 0; i < areas.length; i++) {
    (function(ta) {
      var numId = ta.dataset.id.replace('code','linenums');
      var nums = document.querySelector('[data-id="'+numId+'"]');
      if (!nums || ta._csLinesWired) return;
      ta._csLinesWired = true;
      function sync() {
        var count = ta.value.split('\n').length;
        var s = '';
        for (var n = 1; n <= count; n++) { s += (n > 1 ? '\n' : '') + n; }
        nums.textContent = s;
      }
      ta.addEventListener('input', sync);
      sync();
    })(areas[i]);
  }
}
document.addEventListener('DOMContentLoaded', csWireLineNums);
new MutationObserver(csWireLineNums).observe(document.body || document.documentElement, { childList: true, subtree: true });

// ─── on:tab — wire Tab key on inputs for completion ──────────────────────────
function csWireOnTab() {
  var els = document.querySelectorAll('[data-on-tab]');
  for (var i = 0; i < els.length; i++) {
    (function(el) {
      if (el._csTabWired) return;
      el._csTabWired = true;
      var tabIndex = -1;
      var tabMatches = [];
      var tabPrefix = '';
      el.addEventListener('keydown', function(e) {
        if (e.key !== 'Tab') { tabIndex = -1; tabMatches = []; return; }
        e.preventDefault();
        var val = el.value;
        if (!val || val.length < 1) return;
        var src = el.dataset.completions;
        if (!src) return;
        var items;
        try { items = JSON.parse(src); } catch(x) { return; }
        // Split input into words — complete the last word
        var parts = val.split(/\s+/);
        var lastWord = parts[parts.length - 1].toLowerCase();
        if (lastWord.length < 1) return;
        // If prefix changed, reset matches
        if (lastWord !== tabPrefix) {
          tabPrefix = lastWord;
          tabIndex = -1;
          tabMatches = [];
          for (var j = 0; j < items.length; j++) {
            if (items[j].toLowerCase().indexOf(tabPrefix) === 0) {
              tabMatches.push(items[j]);
            }
          }
        }
        if (tabMatches.length === 0) return;
        // Cycle through matches
        tabIndex = (tabIndex + 1) % tabMatches.length;
        parts[parts.length - 1] = tabMatches[tabIndex];
        el.value = parts.join(' ');
      });
    })(els[i]);
  }
}
document.addEventListener('DOMContentLoaded', csWireOnTab);
new MutationObserver(csWireOnTab).observe(document.body || document.documentElement, { childList: true, subtree: true });

// ─── DirtByte panel + file functions (survive navigation) ────────────────────
window.dbTogglePanel = function(name) {
  var root = document.querySelector('.db');
  if (!root) return;
  var current = root.dataset.activePanel;
  root.dataset.activePanel = (current === name) ? '' : name;
  if (document.activeElement) document.activeElement.blur();
};
window.dbCloseFile = function() {
  var root = document.querySelector('.db');
  if (root) root.dataset.fileOpen = 'false';
  var fp = document.querySelector('[data-id="db-file-panel"]');
  if (fp) fp.innerHTML = '';
};

// ─── ESC closes panels globally ──────────────────────────────────────────────
document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    if (typeof dbTogglePanel === 'function') dbTogglePanel('');
    if (document.activeElement) document.activeElement.blur();
  }
});

// ─── Notepad autosave + terminal auto-scroll (re-wires on DOM change) ────────
function csWireDirtbyte() {
  // Notepad autosave
  var np = document.querySelector('[data-id=db-notepad-field]');
  if (np && !np._csAutosaveWired) {
    np._csAutosaveWired = true;
    var timer;
    np.addEventListener('input', function() {
      clearTimeout(timer);
      timer = setTimeout(function() {
        window.csPost('dirtbyte/notepad', { mission: np.dataset.mission, value: np.value });
      }, 600);
    });
  }
  // Terminal auto-scroll to bottom
  var term = document.querySelector('[data-id="db-terminal-output"]');
  if (term) term.scrollTop = term.scrollHeight;

  // Command history — arrow up/down cycling
  var inp = document.querySelector('[data-id="db-terminal-input"]');
  if (inp && !inp._csHistoryWired) {
    inp._csHistoryWired = true;
    var histIdx = -1;
    var savedInput = '';
    inp.addEventListener('keydown', function(e) {
      if (e.key !== 'ArrowUp' && e.key !== 'ArrowDown') {
        histIdx = -1;
        return;
      }
      e.preventDefault();
      var hist;
      try { hist = JSON.parse(inp.dataset.history || '[]'); } catch(x) { return; }
      if (!hist.length) return;
      if (histIdx === -1) savedInput = inp.value;
      if (e.key === 'ArrowUp') {
        histIdx = Math.min(histIdx + 1, hist.length - 1);
      } else {
        histIdx = histIdx - 1;
      }
      if (histIdx < 0) {
        histIdx = -1;
        inp.value = savedInput;
      } else {
        inp.value = hist[hist.length - 1 - histIdx];
      }
    });
  }
}
document.addEventListener('DOMContentLoaded', csWireDirtbyte);
new MutationObserver(csWireDirtbyte).observe(document.body || document.documentElement, { childList: true, subtree: true });

// ─── Client-side navigation — fetch partial, swap body ───────────────────────
window.csNavigate = function(url) {
  flightRecord('nav',0,'navigate:start',url);
  var partialUrl = url.replace('/page/', '/partial/');
  fetch(partialUrl, { credentials: 'same-origin' })
    .then(function(res) { return res.json(); })
    .then(function(data) {
      if (data.redirect) {
        // Server wants a redirect — follow it via partial too
        window.csNavigate(data.redirect);
        return;
      }
      // Swap body content
      document.body.innerHTML = data.html;
      // Execute inline scripts in new content
      document.body.querySelectorAll('script').forEach(function(old) {
        var s = document.createElement('script');
        s.textContent = old.textContent;
        old.parentNode.replaceChild(s, old);
      });
      // Update title
      if (data.title) document.title = data.title;
      // Update URL without reload
      if (data.url) history.pushState({ url: data.url }, data.title, data.url);
      // Scroll to top
      window.scrollTo(0, 0);
      flightRecord('nav',0,'navigate:done',data.url||url);
    })
    .catch(function() {
      flightRecord('nav',2,'navigate:fallback',url);
      // Fallback to full page load on error
      window.location.href = url;
    });
};

// Handle back/forward buttons
window.addEventListener('popstate', function(e) {
  if (e.state && e.state.url) {
    window.csNavigate(e.state.url);
  }
});

// Intercept link clicks — use partial navigation for /page/ links
document.addEventListener('click', function(e) {
  var link = e.target.closest('a[href^="/page/"]');
  if (!link) return;
  e.preventDefault();
  window.csNavigate(link.getAttribute('href'));
});
