// ─── Form utilities ───────────────────────────────────────────────────────────

// csSetButtonLoading — disable a button and show a spinner label
window.csSetButtonLoading = function(btn, loading) {
  if (!btn) return;
  if (loading) {
    btn.disabled = true;
    btn.dataset.origLabel = btn.dataset.origLabel || btn.textContent;
    btn.textContent = 'Loading…';
    btn.classList.add('cs-loading');
  } else {
    btn.disabled = false;
    if (btn.dataset.origLabel) btn.textContent = btn.dataset.origLabel;
    btn.classList.remove('cs-loading');
  }
};

// csClearFieldErrors — remove all field error messages from a form
window.csClearFieldErrors = function(formEl) {
  if (!formEl) return;
  formEl.querySelectorAll('.cs-field-error').forEach(function(el) { el.remove(); });
  formEl.querySelectorAll('.cs-input-error').forEach(function(el) {
    el.classList.remove('cs-input-error');
  });
};

// csShowFieldErrors — highlight inputs and inject error messages
window.csShowFieldErrors = function(formEl, fields) {
  if (!formEl || !fields) return;
  Object.keys(fields).forEach(function(name) {
    var input = formEl.querySelector('[name="' + name + '"]');
    if (!input) return;
    input.classList.add('cs-input-error');
    var msg = document.createElement('div');
    msg.className = 'cs-field-error';
    msg.textContent = fields[name];
    input.parentNode.insertBefore(msg, input.nextSibling);
  });
};

// ─── csPost — fetch-based API calls ──────────────────────────────────────────
window.csPost = function(action, data, opts) {
  opts = opts || {};
  flightRecord('action',0,'csPost:send /api/'+action,(function(d){if(!d)return'{}';var c={};for(var k in d){var v=d[k];c[k]=typeof v==='string'&&v.length>512?'[base64 stripped]':v}return JSON.stringify(c)})(data));
  var btn = opts.btn || null;
  var formEl = opts.form || null;

  if (formEl) window.csClearFieldErrors(formEl);

  fetch('/api/' + action, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'same-origin',
    body: JSON.stringify(data || {})
  })
  .then(function(res) { return res.json(); })
  .then(function(result) {
    flightRecord('action',result.error?2:0,result.error?'csPost:error '+action:result.redirect?'csPost:redirect '+action:'csPost:ok '+action,result.error||result.redirect||result.toast||'');
    if (result.error) {
      window.csSnackbar.show(result.error, 'error');
      if (result.fields && formEl) window.csShowFieldErrors(formEl, result.fields);
      if (btn) csSetButtonLoading(btn, false);
      if (opts.onError) opts.onError(result);
      return;
    }
    if (result.toast) window.csSnackbar.show(result.toast, result.toastVariant || 'success');
    // Clear autosave checkpoint on success
    if (formEl && formEl.dataset.autosave) {
      window.csAutosave.clear(formEl.dataset.autosave);
    }
    if (result.redirect) {
      window.location.href = result.redirect;
      return;
    }
    // ── Data targeting: route response data to DOM elements by data-id ──
    if (result.data) {
      var patches = Array.isArray(result.data) ? result.data : [result.data];
      for (var i = 0; i < patches.length; i++) {
        var p = patches[i];
        // setAttr patch — set attributes on an element by selector
        if (p.setAttr) {
          var attrEl = document.querySelector(p.setAttr);
          if (attrEl && p.attrs) {
            for (var k in p.attrs) attrEl.setAttribute(k, p.attrs[k]);
          }
          flightRecord('domPatch',0,'setAttr '+p.setAttr,JSON.stringify(p.attrs||{}));
          continue;
        }
        if (!p.target) continue;
        var el = document.querySelector('[data-id="' + p.target + '"]');
        if (el) {
          if (p.value !== undefined) { el.value = p.value; el.dispatchEvent(new Event('input',{bubbles:true})); flightRecord('domPatch',0,'value '+p.target,'len='+p.value.length); }
          else if (p.append) el.innerHTML += p.html;
          else el.innerHTML = p.html;
          if (p.scroll) el.scrollTop = el.scrollHeight;
          flightRecord('domPatch',0,'target '+p.target,'html='+((p.html||'').length)+' append='+!!p.append+' scroll='+!!p.scroll);
        }
      }
    }
    if (window.csState && window.csState.rescan) window.csState.rescan();
    if (btn) csSetButtonLoading(btn, false);
    if (opts.onSuccess) opts.onSuccess(result);
  })
  .catch(function() {
    flightRecord('action',2,'csPost:network-error '+action,'');
    if (btn) csSetButtonLoading(btn, false);
    window.csSnackbar.show('Network error — please try again', 'error');
  });
};

// csForm — collect form fields and post as a named action
window.csForm = function(formId, action, btn) {
  var formEl = typeof formId === 'string' ? document.getElementById(formId) : formId;
  if (!formEl) return;
  var data = {};
  new FormData(formEl).forEach(function(val, key) { data[key] = val; });
  window.csPost(action, data, { btn: btn || null, form: formEl });
};

// ─── csAction ────────────────────────────────────────────────────────────
window.csAction = function(action, triggerEl) {
  // post:apiAction:optionalFormId
  if (action.startsWith('post:')) {
    var parts = action.slice(5).split(':');
    var apiAction = parts[0];
    var formId = parts[1];
    var data = {};
    var formEl = null;
    if (formId) {
      formEl = document.getElementById(formId);
      if (formEl) new FormData(formEl).forEach(function(v, k) { data[k] = v; });
    }
    window.csPost(apiAction, data, { btn: triggerEl || null, form: formEl });
    return;
  }
  // navigate:url
  if (action.startsWith('navigate:')) { var dest = action.slice(9); if (dest.indexOf('/page/') === 0) window.csNavigate(dest); else window.location.href = dest; return; }
  if (action.startsWith('modal:')) { window.csModal.open(action.slice(6)); return; }
  if (action.startsWith('drawer:')) { window.csDrawer.open(action.slice(7)); return; }
  if (action.startsWith('snackbar:')) {
    var parts = action.slice(9).split(':');
    window.csSnackbar.show(parts[0], parts[1] || 'info');
    return;
  }
  if (action.startsWith('close-modal:')) { window.csModal.close(action.slice(12)); return; }
  if (action.startsWith('close-drawer:')) { window.csDrawer.close(action.slice(13)); return; }
};

// ─── on:enter — wire Enter key on inputs to csAction ────────────────────
function csWireOnEnter() {
  var els = document.querySelectorAll('[data-on-enter]');
  for (var i = 0; i < els.length; i++) {
    (function(el) {
      if (el._csEnterWired) return;
      el._csEnterWired = true;
      el.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
          el.dispatchEvent(new CustomEvent('cs:escape', { bubbles: true }));
          return;
        }
        if (e.key !== 'Enter') return;
        var val = el.value.trim();
        if (!val) return;
        var action = el.dataset.onEnter;
        // Collect data-* attributes as extra payload
        var extra = { value: val };
        for (var k in el.dataset) {
          if (k !== 'onEnter' && k !== 'clear') extra[k] = el.dataset[k];
        }
        // post:action — merge extra data into the POST body
        if (action.startsWith('post:')) {
          var parts = action.slice(5).split(':');
          var apiAction = parts[0];
          var formId = parts[1];
          var data = {};
          if (formId) {
            var formEl = document.getElementById(formId);
            if (formEl) new FormData(formEl).forEach(function(v, k) { data[k] = v; });
          }
          for (var k in extra) data[k] = extra[k];
          window.csPost(apiAction, data, {
            onSuccess: function() {
              // Push command to history for arrow key cycling
              if (val && el.dataset.history !== undefined) {
                try {
                  var h = JSON.parse(el.dataset.history || '[]');
                  h.push(val);
                  el.dataset.history = JSON.stringify(h);
                } catch(x) {}
              }
              if (el.dataset.clear !== undefined) el.value = '';
            }
          });
          return;
        }
        // Non-post actions pass through
        csAction(action, el);
        if (el.dataset.clear !== undefined) el.value = '';
      });
    })(els[i]);
  }
}
document.addEventListener('DOMContentLoaded', csWireOnEnter);
// Re-wire after dynamic patches (MutationObserver)
new MutationObserver(csWireOnEnter).observe(document.body || document.documentElement, { childList: true, subtree: true });
