// ─── Ripple ───────────────────────────────────────────────────────────────────
document.addEventListener('click', function(e) {
  var el = e.target.closest('[data-ripple]');
  if (!el) return;
  var r = document.createElement('span');
  var rect = el.getBoundingClientRect();
  var size = Math.max(rect.width, rect.height);
  r.style.cssText = 'position:absolute;border-radius:50%;background:rgba(255,255,255,0.25);pointer-events:none;transform:scale(0);animation:cs-ripple 0.5s linear;width:'+size+'px;height:'+size+'px;left:'+(e.clientX-rect.left-size/2)+'px;top:'+(e.clientY-rect.top-size/2)+'px;';
  el.style.position = el.style.position || 'relative';
  el.style.overflow = 'hidden';
  el.appendChild(r);
  setTimeout(function(){ r.remove(); }, 600);
});

// ─── Floating Label ───────────────────────────────────────────────────────────
function updateFloatingLabel(input) {
  var wrap = input.closest('.cs-input__wrap');
  if (!wrap) return;
  var label = wrap.querySelector('.cs-input__label');
  if (!label) return;
  if (input.value || document.activeElement === input || input.placeholder !== ' ') {
    label.classList.add('cs-input__label--float');
  } else {
    label.classList.remove('cs-input__label--float');
  }
}

document.addEventListener('focus', function(e) {
  if (e.target.classList.contains('cs-input__field')) updateFloatingLabel(e.target);
}, true);
document.addEventListener('blur', function(e) {
  if (e.target.classList.contains('cs-input__field')) updateFloatingLabel(e.target);
}, true);
document.addEventListener('input', function(e) {
  if (e.target.classList.contains('cs-input__field')) updateFloatingLabel(e.target);
});

// Init existing inputs
document.querySelectorAll('.cs-input__field').forEach(function(input) {
  updateFloatingLabel(input);
});

// ─── Tabs ─────────────────────────────────────────────────────────────────────
document.addEventListener('click', function(e) {
  var trigger = e.target.closest('[data-tab-trigger]');
  if (!trigger) return;
  var tabId = trigger.getAttribute('data-tab-trigger');
  var tabGroup = trigger.closest('.cs-tabs');
  if (!tabGroup) return;
  // Deactivate all
  tabGroup.querySelectorAll('[data-tab-trigger]').forEach(function(t) {
    t.classList.remove('cs-tab--active');
  });
  tabGroup.querySelectorAll('[data-tab-panel]').forEach(function(p) {
    p.classList.remove('cs-tab-panel--active');
  });
  // Activate selected
  trigger.classList.add('cs-tab--active');
  var panel = tabGroup.querySelector('[data-tab-panel="'+tabId+'"]');
  if (panel) panel.classList.add('cs-tab-panel--active');
});

// ─── Accordion ────────────────────────────────────────────────────────────────
document.addEventListener('click', function(e) {
  var trigger = e.target.closest('[data-accordion-trigger]');
  if (!trigger) return;
  var item = trigger.closest('.cs-accordion-item');
  if (!item) return;
  var body = item.querySelector('.cs-accordion-body');
  var isOpen = item.classList.contains('cs-accordion-item--open');
  item.classList.toggle('cs-accordion-item--open', !isOpen);
  if (body) {
    if (isOpen) {
      body.style.maxHeight = '0';
    } else {
      body.style.maxHeight = body.scrollHeight + 'px';
      // After transition, switch to none so dynamic content isn't clipped
      body.addEventListener('transitionend', function handler() {
        body.removeEventListener('transitionend', handler);
        if (item.classList.contains('cs-accordion-item--open')) {
          body.style.maxHeight = 'none';
        }
      });
    }
  }
});

// ─── Modal ────────────────────────────────────────────────────────────────────
window.csModal = {
  open: function(id) {
    var modal = document.getElementById('modal-'+id);
    if (!modal) return;
    modal.classList.add('cs-modal--open');
    document.body.style.overflow = 'hidden';
  },
  close: function(id) {
    var modal = id ? document.getElementById('modal-'+id) : document.querySelector('.cs-modal--open');
    if (!modal) return;
    modal.classList.remove('cs-modal--open');
    document.body.style.overflow = '';
  }
};

document.addEventListener('click', function(e) {
  if (e.target.classList.contains('cs-modal__backdrop')) {
    window.csModal.close();
  }
  var closeBtn = e.target.closest('[data-modal-close]');
  if (closeBtn) {
    var id = closeBtn.getAttribute('data-modal-close');
    window.csModal.close(id);
  }
});

// ─── Drawer ───────────────────────────────────────────────────────────────────
window.csDrawer = {
  open: function(id) {
    var drawer = document.getElementById('drawer-'+id);
    if (!drawer) return;
    drawer.classList.add('cs-drawer--open');
    document.body.style.overflow = 'hidden';
  },
  close: function(id) {
    var drawer = id ? document.getElementById('drawer-'+id) : document.querySelector('.cs-drawer--open');
    if (!drawer) return;
    drawer.classList.remove('cs-drawer--open');
    document.body.style.overflow = '';
  }
};

document.addEventListener('click', function(e) {
  if (e.target.classList.contains('cs-drawer__backdrop')) {
    window.csDrawer.close();
  }
  var closeBtn = e.target.closest('[data-drawer-close]');
  if (closeBtn) {
    var id = closeBtn.getAttribute('data-drawer-close');
    window.csDrawer.close(id);
  }
});

// ─── Snackbar ─────────────────────────────────────────────────────────────────
var snackbarQueue = [];
var snackbarActive = false;

window.csSnackbar = {
  show: function(message, variant) {
    snackbarQueue.push({ message: message, variant: variant || 'info' });
    if (!snackbarActive) processSnackbar();
  }
};

function processSnackbar() {
  if (!snackbarQueue.length) { snackbarActive = false; return; }
  snackbarActive = true;
  var item = snackbarQueue.shift();
  var host = document.getElementById('cs-snackbar-host');
  if (!host) {
    host = document.createElement('div');
    host.id = 'cs-snackbar-host';
    host.style.cssText = 'position:fixed;top:24px;left:50%;transform:translateX(-50%);z-index:9999;display:flex;flex-direction:column;gap:8px;align-items:center;pointer-events:none;';
    document.body.appendChild(host);
  }
  var el = document.createElement('div');
  el.className = 'cs-snackbar cs-snackbar--'+item.variant;
  el.textContent = item.message;
  el.style.cssText = 'pointer-events:all;';
  host.appendChild(el);
  requestAnimationFrame(function() { el.classList.add('cs-snackbar--show'); });
  setTimeout(function() {
    el.classList.remove('cs-snackbar--show');
    setTimeout(function() { el.remove(); processSnackbar(); }, 300);
  }, 3000);
}

// ─── Autocomplete ─────────────────────────────────────────────────────────────
document.addEventListener('input', function(e) {
  if (!e.target.hasAttribute('data-autocomplete')) return;
  var input = e.target;
  var wrap = input.closest('.cs-autocomplete');
  if (!wrap) return;
  var dropdown = wrap.querySelector('.cs-autocomplete__dropdown');
  if (!dropdown) return;
  var query = input.value.toLowerCase().trim();
  var items = dropdown.querySelectorAll('[data-ac-item]');
  var visible = 0;
  items.forEach(function(item) {
    var match = !query || item.textContent.toLowerCase().includes(query);
    item.style.display = match ? '' : 'none';
    if (match) visible++;
  });
  dropdown.style.display = (visible > 0 && query !== '') ? '' : 'none';
});

document.addEventListener('click', function(e) {
  var item = e.target.closest('[data-ac-item]');
  if (item) {
    var wrap = item.closest('.cs-autocomplete');
    var input = wrap ? wrap.querySelector('[data-autocomplete]') : null;
    if (input) {
      input.value = item.textContent;
      var dropdown = wrap.querySelector('.cs-autocomplete__dropdown');
      if (dropdown) dropdown.style.display = 'none';
    }
    return;
  }
  // Close dropdown when clicking outside
  document.querySelectorAll('.cs-autocomplete__dropdown').forEach(function(d) {
    if (!d.closest('.cs-autocomplete').contains(e.target)) {
      d.style.display = 'none';
    }
  });
});

document.addEventListener('keydown', function(e) {
  if (!e.target.hasAttribute('data-autocomplete')) return;
  var wrap = e.target.closest('.cs-autocomplete');
  if (!wrap) return;
  var dropdown = wrap.querySelector('.cs-autocomplete__dropdown');
  if (!dropdown || dropdown.style.display === 'none') return;
  var items = Array.from(dropdown.querySelectorAll('[data-ac-item]')).filter(function(i) { return i.style.display !== 'none'; });
  var active = dropdown.querySelector('[data-ac-item].cs-autocomplete__item--active');
  var idx = active ? items.indexOf(active) : -1;
  if (e.key === 'ArrowDown') {
    e.preventDefault();
    if (active) active.classList.remove('cs-autocomplete__item--active');
    idx = (idx + 1) % items.length;
    items[idx].classList.add('cs-autocomplete__item--active');
  } else if (e.key === 'ArrowUp') {
    e.preventDefault();
    if (active) active.classList.remove('cs-autocomplete__item--active');
    idx = (idx - 1 + items.length) % items.length;
    items[idx].classList.add('cs-autocomplete__item--active');
  } else if (e.key === 'Enter' && active) {
    e.target.value = active.textContent;
    dropdown.style.display = 'none';
  } else if (e.key === 'Escape') {
    dropdown.style.display = 'none';
  }
});

// ─── Select ───────────────────────────────────────────────────────────────────
document.addEventListener('click', function(e) {
  var trigger = e.target.closest('[data-select-trigger]');
  if (trigger) {
    var wrap = trigger.closest('.cs-select');
    var dropdown = wrap ? wrap.querySelector('.cs-select__dropdown') : null;
    if (dropdown) {
      var isOpen = dropdown.style.display !== 'none';
      document.querySelectorAll('.cs-select__dropdown').forEach(function(d) { d.style.display = 'none'; });
      dropdown.style.display = isOpen ? 'none' : '';
    }
    return;
  }
  var option = e.target.closest('[data-select-option]');
  if (option) {
    var wrap = option.closest('.cs-select');
    var trigger = wrap ? wrap.querySelector('[data-select-trigger]') : null;
    var valueEl = trigger ? trigger.querySelector('.cs-select__value') : null;
    var hidden = wrap ? wrap.querySelector('input[type=hidden]') : null;
    if (valueEl) valueEl.textContent = option.textContent;
    if (hidden) hidden.value = option.getAttribute('data-select-option');
    var dropdown = wrap ? wrap.querySelector('.cs-select__dropdown') : null;
    if (dropdown) dropdown.style.display = 'none';
    wrap.querySelectorAll('[data-select-option]').forEach(function(o) { o.classList.remove('cs-select__option--active'); });
    option.classList.add('cs-select__option--active');
    return;
  }
  // Close on outside click
  document.querySelectorAll('.cs-select__dropdown').forEach(function(d) {
    if (!d.closest('.cs-select').contains(e.target)) d.style.display = 'none';
  });
});

// ─── Rating ───────────────────────────────────────────────────────────────────
document.addEventListener('click', function(e) {
  var star = e.target.closest('[data-rating-value]');
  if (!star) return;
  var wrap = star.closest('.cs-rating');
  if (!wrap || wrap.hasAttribute('data-rating-readonly')) return;
  var val = parseInt(star.getAttribute('data-rating-value'), 10);
  wrap.setAttribute('data-rating', val);
  var input = wrap.querySelector('[data-rating-input]');
  if (input) input.value = val;
  wrap.querySelectorAll('[data-rating-value]').forEach(function(s) {
    var sv = parseInt(s.getAttribute('data-rating-value'), 10);
    s.classList.toggle('cs-rating__star--filled', sv <= val);
  });
});

// ─── Slider ───────────────────────────────────────────────────────────────────
window.csSliderUpdate = function(input) {
  var id = input.getAttribute('data-slider-id');
  var display = document.querySelector('[data-slider-value="'+id+'"]');
  if (display) display.textContent = input.value;
};

// ─── NumberInput ──────────────────────────────────────────────────────────────
window.csNumberStep = function(btn, dir) {
  var wrap = btn.closest('.cs-number-input');
  var field = wrap ? wrap.querySelector('.cs-number-input__field') : null;
  if (!field) return;
  var step = parseFloat(field.step) || 1;
  var val = parseFloat(field.value) || 0;
  var min = field.min !== '' ? parseFloat(field.min) : -Infinity;
  var max = field.max !== '' ? parseFloat(field.max) : Infinity;
  val = Math.min(max, Math.max(min, val + dir * step));
  field.value = val;
};

// ─── TagInput ─────────────────────────────────────────────────────────────────
window.csTagKeydown = function(e, input) {
  if (e.key !== 'Enter' && e.key !== ',') return;
  e.preventDefault();
  var val = input.value.trim().replace(/,$/, '');
  if (!val) return;
  var wrap = input.closest('[data-tag-input]');
  if (!wrap) return;
  var id = wrap.getAttribute('data-tag-input');
  var tag = document.createElement('span');
  tag.className = 'cs-tag-input__tag';
  tag.innerHTML = val + '<button type="button" class="cs-tag-input__remove" onclick="csTagRemove(this)" aria-label="Remove">\xd7</button>';
  wrap.insertBefore(tag, input);
  input.value = '';
  csTagSync(id);
};

window.csTagRemove = function(btn) {
  var tag = btn.closest('.cs-tag-input__tag');
  if (!tag) return;
  var wrap = tag.closest('[data-tag-input]');
  var id = wrap ? wrap.getAttribute('data-tag-input') : null;
  tag.remove();
  if (id) csTagSync(id);
};

function csTagSync(id) {
  var wrap = document.querySelector('[data-tag-input="'+id+'"]');
  var hidden = document.querySelector('[data-tag-value="'+id+'"]');
  if (!wrap || !hidden) return;
  var tags = Array.from(wrap.querySelectorAll('.cs-tag-input__tag')).map(function(t) {
    return t.firstChild.textContent.trim();
  });
  hidden.value = tags.join(',');
}

// ─── FileUpload ───────────────────────────────────────────────────────────────
window.csFileUploadChange = function(input) {
  var wrap = input.closest('[data-file-upload]');
  var id = wrap ? wrap.getAttribute('data-id') : null;
  csFileShowList(id, input.files);
};

window.csFileDragOver = function(e, zone) {
  e.preventDefault();
  zone.classList.add('cs-file-upload__zone--drag');
};

window.csFileDragLeave = function(zone) {
  zone.classList.remove('cs-file-upload__zone--drag');
};

window.csFileDrop = function(e, zone, id) {
  e.preventDefault();
  zone.classList.remove('cs-file-upload__zone--drag');
  csFileShowList(id, e.dataTransfer.files);
};

function csFileShowList(id, files) {
  var list = document.querySelector('[data-file-list="'+id+'"]');
  if (!list || !files) return;
  list.innerHTML = '';
  Array.from(files).forEach(function(f) {
    var item = document.createElement('div');
    item.className = 'cs-file-upload__file';
    item.innerHTML = '<span>'+f.name+'</span><span>'+(f.size > 1024*1024 ? (f.size/1024/1024).toFixed(1)+'MB' : (f.size/1024).toFixed(0)+'KB')+'</span>';
    list.appendChild(item);
  });
}

// ─── CodeBlock copy ───────────────────────────────────────────────────────────
window.csCopyCode = function(id) {
  var pre = document.getElementById(id);
  if (!pre) return;
  navigator.clipboard.writeText(pre.textContent).then(function() {
    var btn = document.querySelector('[onclick="csCopyCode(\''+id+'\')"]');
    if (btn) { btn.textContent = 'Copied!'; setTimeout(function(){ btn.textContent = 'Copy'; }, 1500); }
  });
};

// ─── Menu ─────────────────────────────────────────────────────────────────────
document.addEventListener('click', function(e) {
  var trigger = e.target.closest('[data-menu-trigger]');
  if (trigger) {
    var id = trigger.getAttribute('data-menu-trigger');
    var dropdown = document.querySelector('[data-menu-dropdown="'+id+'"]');
    if (dropdown) {
      var isOpen = dropdown.classList.contains('cs-menu--open');
      document.querySelectorAll('.cs-menu__dropdown').forEach(function(d) { d.classList.remove('cs-menu--open'); });
      if (!isOpen) dropdown.classList.add('cs-menu--open');
    }
    return;
  }
  document.querySelectorAll('.cs-menu__dropdown.cs-menu--open').forEach(function(d) {
    if (!d.closest('.cs-menu').contains(e.target)) d.classList.remove('cs-menu--open');
  });
});

// ─── Popover ──────────────────────────────────────────────────────────────────
document.addEventListener('click', function(e) {
  var trigger = e.target.closest('[data-popover-trigger]');
  if (trigger) {
    var id = trigger.getAttribute('data-popover-trigger');
    var panel = document.querySelector('[data-popover-panel="'+id+'"]');
    if (panel) {
      var isOpen = panel.classList.contains('cs-popover--open');
      document.querySelectorAll('.cs-popover__panel').forEach(function(p) { p.classList.remove('cs-popover--open'); });
      if (!isOpen) panel.classList.add('cs-popover--open');
    }
    return;
  }
  document.querySelectorAll('.cs-popover__panel.cs-popover--open').forEach(function(p) {
    if (!p.closest('.cs-popover').contains(e.target)) p.classList.remove('cs-popover--open');
  });
});

// ─── Stepper init ─────────────────────────────────────────────────────────────
document.addEventListener('DOMContentLoaded', function() {
  document.querySelectorAll('[data-stepper]').forEach(function(stepper) {
    var current = parseInt(stepper.getAttribute('data-current'), 10) || 1;
    stepper.querySelectorAll('.cs-stepper-step').forEach(function(step) {
      var n = parseInt(step.getAttribute('data-step'), 10);
      step.classList.toggle('cs-stepper-step--active', n === current);
      step.classList.toggle('cs-stepper-step--done', n < current);
    });
  });
});

// ─── Autosave checkpoints ─────────────────────────────────────────────────────
// Forms with data-autosave="key" are automatically saved to localStorage on
// input change (debounced 600ms) and restored on page load.
// Checkpoints are cleared on successful API response.

window.csAutosave = (function() {
  var timers = {};

  function save(key, formEl) {
    var data = {};
    new FormData(formEl).forEach(function(val, k) { data[k] = val; });
    try { localStorage.setItem('cs_autosave_' + key, JSON.stringify(data)); } catch(e) {}
  }

  function restore(key, formEl) {
    var raw;
    try { raw = localStorage.getItem('cs_autosave_' + key); } catch(e) { return; }
    if (!raw) return;
    var data;
    try { data = JSON.parse(raw); } catch(e) { return; }
    Object.keys(data).forEach(function(name) {
      var el = formEl.querySelector('[name="' + name + '"]');
      if (!el) return;
      if (el.type === 'checkbox' || el.type === 'radio') {
        el.checked = data[name] === 'on' || data[name] === el.value;
      } else {
        el.value = data[name];
        // Fire input event so any bound components (tag-input, slider etc.) update
        el.dispatchEvent(new Event('input', { bubbles: true }));
      }
    });
  }

  function clear(key) {
    try { localStorage.removeItem('cs_autosave_' + key); } catch(e) {}
  }

  function wire(formEl) {
    var key = formEl.dataset.autosave;
    if (!key) return;
    restore(key, formEl);
    formEl.addEventListener('input', function() {
      clearTimeout(timers[key]);
      timers[key] = setTimeout(function() { save(key, formEl); }, 600);
    });
    formEl.addEventListener('change', function() {
      clearTimeout(timers[key]);
      timers[key] = setTimeout(function() { save(key, formEl); }, 600);
    });
  }

  // Wire all autosave forms on load
  document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('[data-autosave]').forEach(wire);
  });

  return { save: save, restore: restore, clear: clear, wire: wire };
})();
