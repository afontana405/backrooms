// ── csState — client-side state system ──────────────────────────────
// ── Emoji splitter — keeps multi-codepoint sequences intact ─────────
window.splitEmoji = function(str) {
  if (!str) return [];
  if (typeof Intl !== 'undefined' && Intl.Segmenter) {
    var seg = new Intl.Segmenter('en', { granularity: 'grapheme' });
    return Array.from(seg.segment(str), function(s) { return s.segment; });
  }
  return Array.from(str);
};

window.csState = (function(){
  var _bindings = {};

  function scan(feature) {
    var prefix = feature + '.';
    _bindings[feature] = {};
    document.querySelectorAll('[data-state]').forEach(function(el) {
      var key = el.getAttribute('data-state');
      if (key.indexOf(prefix) === 0) {
        var field = key.substring(prefix.length);
        if (!_bindings[feature][field]) _bindings[feature][field] = [];
        _bindings[feature][field].push({el: el, type: 'text'});
      }
    });
    document.querySelectorAll('[data-state-attr]').forEach(function(el) {
      var val = el.getAttribute('data-state-attr');
      var parts = val.split(':');
      if (parts[0].indexOf(prefix) === 0) {
        var field = parts[0].substring(prefix.length);
        if (!_bindings[feature][field]) _bindings[feature][field] = [];
        _bindings[feature][field].push({el: el, type: 'attr', attr: parts[1]});
      }
    });
  }

  function updateBindings(feature, field) {
    var fb = _bindings[feature];
    if (!fb || !fb[field]) return;
    var value = window['__' + feature][field];
    var list = fb[field];
    for (var i = 0; i < list.length; i++) {
      var b = list[i];
      if (b.type === 'text') b.el.textContent = (value == null) ? '' : value;
      else if (b.type === 'attr') {
        if (value === false || value == null) b.el.removeAttribute(b.attr);
        else b.el.setAttribute(b.attr, value === true ? '' : value);
      }
    }
  }

  return {
    init: function(feature, data) {
      window['__' + feature] = data;
      scan(feature);
      var state = data;
      var keys = [];
      for (var k in state) {
        if (state.hasOwnProperty(k)) { updateBindings(feature, k); keys.push(k); }
      }
      flightRecord('state', 0, 'state:init ' + feature, keys.join(','));
    },
    set: function(feature, changes) {
      var state = window['__' + feature];
      if (!state) return;
      var keys = [];
      for (var k in changes) {
        if (changes.hasOwnProperty(k)) {
          state[k] = changes[k];
          updateBindings(feature, k);
          keys.push(k);
        }
      }
      flightRecord('state', 0, 'state:set ' + feature, keys.join(','));
    },
    get: function(feature) {
      return window['__' + feature];
    },
    rescan: function() {
      for (var feature in _bindings) {
        scan(feature);
        var state = window['__' + feature];
        if (state) {
          for (var k in state) {
            if (state.hasOwnProperty(k)) updateBindings(feature, k);
          }
        }
      }
    }
  };
})();
