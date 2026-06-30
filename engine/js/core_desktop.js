// ── F-key shortcuts ─────────────────────────────────────────────────────────
document.addEventListener('keydown', function(e) {
  if (e.key === 'F8') {
    e.preventDefault();
    var b = document.getElementById('diag-btn');
    var p = document.getElementById('diag-panel');
    if (!b || !p) return;
    if (b.classList.contains('diag-btn--hidden')) {
      b.classList.remove('diag-btn--hidden');
      p.classList.remove('diag-panel--closed');
    } else if (!p.classList.contains('diag-panel--closed')) {
      if (window.__diagClose) window.__diagClose();
    } else {
      p.classList.remove('diag-panel--closed');
    }
  }
  if (e.key === 'F9') {
    e.preventDefault();
    csSnackbar.show('Capturing screenshot...', 'info');
    csPost('screenshot/save', {});
  }
  if (e.key === 'F10') {
    e.preventDefault();
    tgtToggle();
  }
});

// ─── Desktop OS ──────────────────────────────────────────────────────────────
var openWindows = {};

var openAppLock = {};
window.openApp = function(page, title, w, h) {
  // If already open, focus it
  if (openWindows[page] && !openWindows[page].closed) {
    openWindows[page].focus();
    return;
  }

  // Debounce rapid clicks — block for 1 second
  if (openAppLock[page]) return;
  openAppLock[page] = true;
  setTimeout(function() { delete openAppLock[page]; }, 1000);

  var left = Math.max(0, (screen.width - w) / 2 + Math.random() * 60 - 30);
  var top = Math.max(0, (screen.height - h) / 2 + Math.random() * 40 - 20);
  var features = 'width=' + w + ',height=' + h + ',left=' + Math.round(left) + ',top=' + Math.round(top) + ',resizable=yes,scrollbars=yes';

  var win = window.open('/page/' + page, 'app_' + page, features);
  openWindows[page] = win;

  // Add to taskbar
  var bar = document.getElementById('taskbar-apps');
  if (bar) {
    var existing = document.getElementById('tb-' + page);
    if (!existing) {
      var btn = document.createElement('div');
      btn.id = 'tb-' + page;
      btn.className = 'taskbar-btn';
      btn.textContent = title;
      btn.onclick = function() {
        if (openWindows[page] && !openWindows[page].closed) {
          openWindows[page].focus();
        }
      };
      bar.appendChild(btn);
    }
  }

  // Clean up taskbar when window closes
  var checkClosed = setInterval(function() {
    if (!openWindows[page] || openWindows[page].closed) {
      clearInterval(checkClosed);
      delete openWindows[page];
      var tbBtn = document.getElementById('tb-' + page);
      if (tbBtn) tbBtn.remove();
    }
  }, 1000);
};

// Taskbar clock
function updateClock() {
  var el = document.getElementById('taskbar-clock');
  if (!el) return;
  var now = new Date();
  var h = now.getHours();
  var m = now.getMinutes();
  var ampm = h >= 12 ? 'PM' : 'AM';
  h = h % 12 || 12;
  el.textContent = h + ':' + (m < 10 ? '0' : '') + m + ' ' + ampm;
}
updateClock();
setInterval(updateClock, 30000);

// ─── Data Chat (natural language query widget) ───────────────────────────────
window.csDataChat = function(chatId, schema) {
  var input = document.getElementById(chatId + '-input');
  var messages = document.getElementById(chatId + '-messages');
  if (!input || !messages) return;

  var question = input.value.trim();
  if (!question) return;
  input.value = '';

  // Show user message
  var userMsg = document.createElement('div');
  userMsg.className = 'mb-2 text-end';
  userMsg.innerHTML = '<span class="badge bg-primary text-wrap text-start" style="max-width:80%;font-size:0.85rem">' + question + '</span>';
  messages.appendChild(userMsg);

  // Show loading
  var loading = document.createElement('div');
  loading.className = 'mb-2 text-body-secondary small';
  loading.textContent = 'Thinking...';
  messages.appendChild(loading);
  messages.scrollTop = messages.scrollHeight;

  // Call API
  fetch('/api/data/ask', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ question: question, schema: schema })
  })
  .then(function(r) { return r.json(); })
  .then(function(result) {
    loading.remove();
    var answer = (result.data && result.data.answer) || result.error || 'No response.';
    var answerMsg = document.createElement('div');
    answerMsg.className = 'mb-2';
    answerMsg.innerHTML = '<div class="card card-body bg-body-tertiary" style="font-size:0.85rem;white-space:pre-wrap">' + answer + '</div>';
    messages.appendChild(answerMsg);
    messages.scrollTop = messages.scrollHeight;
  })
  .catch(function(err) {
    loading.remove();
    var errMsg = document.createElement('div');
    errMsg.className = 'mb-2 text-danger small';
    errMsg.textContent = 'Error: ' + err.message;
    messages.appendChild(errMsg);
  });
};

// ─── Binary Data Fetch ───────────────────────────────────────────────────────
// Generic: fetch a binary collection and call back with decoded records.
window.csBinaryFetch = function(collection, callback) {
  fetch('/api/binary/list', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ collection: collection })
  })
  .then(function(r) { return r.json(); })
  .then(function(result) {
    if (result.data && result.data.records) {
      callback(result.data.records);
    }
  });
};

// ─── Page Data Renderers ─────────────────────────────────────────────────────
// Each page with data-page-type auto-loads its binary data on DOMContentLoaded.

document.addEventListener('DOMContentLoaded', function() {
  var pageType = document.body.getAttribute('data-page-type');
  if (!pageType) return;

  var renderers = {
    'email': function() {
      csBinaryFetch('emails', function(records) {
        var container = document.getElementById('email-list');
        if (!container) return;
        var html = '';
        records.forEach(function(e) {
          var priTag = '';
          if (e.priority && e.priority !== '') {
            var priClass = e.priority === 'P1' ? 'text-bg-danger' : e.priority === 'P2' ? 'text-bg-warning' : 'text-bg-info';
            priTag = '<span class="badge ' + priClass + ' me-2">' + e.priority + '</span>';
          }
          html += '<div class="card mb-3"><div class="card-body">' +
            '<div class="d-flex align-items-start gap-2">' + priTag +
            '<div class="flex-grow-1">' +
            '<h6 class="mb-1">' + (e.subject || '') + '</h6>' +
            '<small class="text-body-secondary">From: ' + (e.fromName || '') + ' (' + (e.fromEmail || '') + ')</small>' +
            '</div></div><hr class="my-2">' +
            '<p class="mb-0" style="font-size:0.9rem">' + (e.body || '') + '</p>' +
            '</div></div>';
        });
        container.innerHTML = html;
      });
    },

    'tickets': function() {
      csBinaryFetch('tickets', function(records) {
        var container = document.getElementById('ticket-list');
        if (!container) return;
        var html = '';
        records.forEach(function(t) {
          var priClass = t.priority === 'P1' ? 'danger' : t.priority === 'P2' ? 'warning' : 'info';
          var borderStyle = t.priority === 'P1' ? 'border-left:4px solid #dc3545' : '';
          html += '<div class="card mb-3" style="' + borderStyle + '"><div class="card-body">' +
            '<div class="d-flex gap-3 align-items-start">' +
            '<div><strong class="fs-5">TK-' + String(t.id).padStart(5,'0') + '</strong><br>' +
            '<span class="badge text-bg-' + priClass + '">' + t.priority + '</span></div>' +
            '<div class="flex-grow-1"><h6 class="mb-1">' + (t.subject || '') + '</h6>' +
            '<small class="text-body-secondary">Requester: ' + (t.requester || '') +
            ' | Assigned: ' + (t.assignee || 'Unassigned') + '</small></div>' +
            '<div class="text-end"><span class="badge text-bg-' + (t.status==='Open'?'secondary':'primary') + '">' + t.status + '</span>' +
            (t.sla === 'BREACHED' ? '<br><small class="text-danger">SLA BREACHED</small>' : '') +
            '</div></div>' +
            '<hr class="my-2"><p class="mb-0" style="font-size:0.9rem">' + (t.description || '') + '</p>' +
            '</div></div>';
        });
        container.innerHTML = html;
      });
    },

    'siem': function() {
      csBinaryFetch('alerts', function(alerts) {
        var container = document.getElementById('alert-list');
        if (!container) return;
        var html = '';
        alerts.forEach(function(a) {
          var sevClass = a.severity === 'CRITICAL' || a.severity === 'HIGH' ? 'danger' :
                         a.severity === 'MEDIUM' ? 'warning' : 'info';
          html += '<div class="card mb-3" style="border-left:4px solid var(--bs-' + sevClass + ')"><div class="card-body">' +
            '<div class="d-flex gap-3 align-items-start">' +
            '<div><strong>' + (a.alertId || '') + '</strong><br>' +
            '<span class="badge text-bg-' + sevClass + '">' + a.severity + '</span></div>' +
            '<div class="flex-grow-1"><h6 class="mb-1">' + (a.title || '') + '</h6>' +
            '<small class="text-body-secondary">Source: ' + (a.source || '') + ' → ' + (a.destination || '') + ' | ' + (a.protocol || '') + '</small>' +
            '</div></div></div></div>';
        });
        container.innerHTML = html;
      });

      csBinaryFetch('events', function(events) {
        var container = document.getElementById('event-table-body');
        if (!container) return;
        var html = '';
        events.forEach(function(e) {
          var sevClass = e.severity === 'HIGH' || e.severity === 'CRITICAL' ? 'table-danger' :
                         e.severity === 'MEDIUM' ? 'table-warning' : '';
          html += '<tr class="' + sevClass + '">' +
            '<td>' + (e.time || '') + '</td>' +
            '<td>' + (e.source || '') + '</td>' +
            '<td>' + (e.event || '') + '</td>' +
            '<td>' + (e.severity || '') + '</td>' +
            '<td>' + (e.status || '') + '</td></tr>';
        });
        container.innerHTML = html;
      });
    },

    'chat': function() {
      csBinaryFetch('messages', function(records) {
        var container = document.getElementById('chat-messages');
        if (!container) return;
        // Filter to active conversation (jake by default)
        var msgs = records.filter(function(m) { return m.conversationId === 'conv-jake'; });
        var html = '';
        msgs.forEach(function(m) {
          var ts = m.timestamp || '';
          if (ts.indexOf('T') > -1) ts = ts.split('T')[1].substring(0,5);
          var fromName = m.from === 'chen' ? 'Jake Chen' : m.from === 'player' ? 'You' : m.from;
          html += '<div class="card mb-2" style="background:var(--bs-body-bg)">' +
            '<div class="card-body py-2 px-3">' +
            '<small class="text-success fw-bold">' + fromName + ' — ' + ts + '</small>' +
            '<p class="mb-0">' + (m.text || '') + '</p>' +
            '</div></div>';
        });
        container.innerHTML = html;
      });
    },

    'directory': function() {
      csBinaryFetch('staff', function(records) {
        var container = document.getElementById('staff-list');
        if (!container) return;
        var depts = {};
        records.forEach(function(s) {
          var dept = s.department || 'Other';
          if (!depts[dept]) depts[dept] = [];
          depts[dept].push(s);
        });
        var html = '';
        Object.keys(depts).forEach(function(dept) {
          html += '<h4 class="mt-4">' + dept + '</h4><div class="row g-3 mb-3">';
          depts[dept].forEach(function(s) {
            var statusClass = s.status === 'available' ? 'success' :
                              s.status === 'away' ? 'warning' :
                              s.status === 'on-leave' ? 'danger' : 'secondary';
            var statusLabel = s.status === 'on-leave' ? 'On Leave' :
                              s.status.charAt(0).toUpperCase() + s.status.slice(1);
            var initials = (s.name || '').split(' ').map(function(w) { return w[0]; }).join('').substring(0,2);
            html += '<div class="col"><div class="card"><div class="card-body">' +
              '<div class="rounded-circle bg-secondary bg-opacity-25 d-flex align-items-center justify-content-center text-body-secondary fw-bold mb-2" style="width:56px;height:56px;font-size:1.3rem">' + initials + '</div>' +
              '<h6 class="mb-1">' + (s.name || '') + '</h6>' +
              '<small class="text-body-secondary d-block">' + (s.title || '') + '</small>' +
              '<span class="badge text-bg-' + statusClass + ' mt-1">' + statusLabel + '</span>' +
              '<small class="text-body-secondary d-block mt-2">' + (s.email || '') + '</small>' +
              '<small class="text-body-secondary">Ext: ' + (s.ext || '') + ' | Building ' + (s.building || '') + ', Floor ' + (s.floor || '') + '</small>' +
              '</div></div></div>';
          });
          html += '</div>';
        });
        container.innerHTML = html;
      });
    },

    'group-inbox': function() {
      csBinaryFetch('group_emails', function(records) {
        var container = document.getElementById('group-list');
        if (!container) return;
        var html = '';
        records.forEach(function(e) {
          var claimed = e.claimedBy && e.claimedBy !== '';
          var tagClass = claimed ? 'info' : 'warning';
          var tagText = claimed ? 'Claimed — ' + e.claimedBy : 'Unclaimed';
          html += '<div class="card mb-3"><div class="card-body">' +
            '<div class="d-flex align-items-start gap-2">' +
            '<span class="badge text-bg-' + tagClass + '">' + tagText + '</span>' +
            '<div class="flex-grow-1">' +
            '<h6 class="mb-1">' + (e.subject || '') + '</h6>' +
            '<small class="text-body-secondary">From: ' + (e.fromName || '') + '</small>' +
            '</div></div><hr class="my-2">' +
            '<p class="mb-0" style="font-size:0.9rem">' + (e.body || '') + '</p>' +
            '</div></div>';
        });
        container.innerHTML = html;
      });
    },

    'firewall': function() {
      function makeTable(cols, records, fields) {
        var html = '<div class="table-responsive"><table class="table table-striped table-hover"><thead><tr>';
        cols.forEach(function(c) { html += '<th>' + c + '</th>'; });
        html += '</tr></thead><tbody>';
        records.forEach(function(r) {
          html += '<tr>';
          fields.forEach(function(f) { html += '<td>' + (r[f] || '') + '</td>'; });
          html += '</tr>';
        });
        html += '</tbody></table></div>';
        return html;
      }

      csBinaryFetch('fw_interfaces', function(r) {
        var el = document.getElementById('fw-interfaces');
        if (el) el.innerHTML = makeTable(['Port','Name','IP Address','Status','Speed','In/Out'], r, ['port','name','ip','status','speed','traffic']);
      });
      csBinaryFetch('fw_policies', function(r) {
        var el = document.getElementById('fw-policies');
        if (el) el.innerHTML = makeTable(['ID','Name','Source','Destination','Service','Action','Hits'], r, ['policyId','name','source','destination','service','action','hits']);
      });
      csBinaryFetch('fw_threats', function(r) {
        var el = document.getElementById('fw-threats');
        if (el) el.innerHTML = makeTable(['Time','Source','Destination','Threat','Action','Severity'], r, ['time','source','destination','threat','action','severity']);
      });
      csBinaryFetch('fw_vpn', function(r) {
        var el = document.getElementById('fw-vpn');
        if (el) el.innerHTML = makeTable(['Name','Remote Gateway','Status','Type','Incoming','Outgoing'], r, ['name','remoteGw','status','vpnType','incoming','outgoing']);
      });
    },

    'switch': function() {
      csBinaryFetch('sw_ports', function(records) {
        var el = document.getElementById('sw-ports');
        if (!el) return;
        var colors = { active: '#4ade80', warning: '#facc15', flagged: '#f87171', trunk: '#60a5fa', unused: '#333' };
        var textColors = { active: '#000', warning: '#000', flagged: '#000', trunk: '#000', unused: '#666' };
        var html = '<div style="display:grid;grid-template-columns:repeat(24,1fr);gap:4px;margin-bottom:16px">';
        records.forEach(function(p) {
          var bg = colors[p.status] || '#333';
          var fg = textColors[p.status] || '#666';
          var title = 'Gi0/' + p.portNum + (p.vlan ? ' — VLAN ' + p.vlan : '') + (p.device ? ' — ' + p.device : '');
          html += '<div style="background:' + bg + ';border-radius:3px;padding:4px;text-align:center;font-size:0.6rem;color:' + fg + ';font-weight:bold;cursor:default" title="' + title + '">' + p.portNum + '</div>';
        });
        html += '</div>';
        el.innerHTML = html;
      });

      function makeTable(cols, records, fields) {
        var html = '<div class="table-responsive"><table class="table table-striped table-hover"><thead><tr>';
        cols.forEach(function(c) { html += '<th>' + c + '</th>'; });
        html += '</tr></thead><tbody>';
        records.forEach(function(r) {
          html += '<tr>';
          fields.forEach(function(f) { html += '<td>' + (r[f] || '') + '</td>'; });
          html += '</tr>';
        });
        html += '</tbody></table></div>';
        return html;
      }

      csBinaryFetch('sw_vlans', function(r) {
        var el = document.getElementById('sw-vlans');
        if (el) el.innerHTML = makeTable(['VLAN ID','Name','Subnet','Gateway','Ports','Status'], r, ['vlanId','name','subnet','gateway','ports','status']);
      });
      csBinaryFetch('sw_mac', function(r) {
        var el = document.getElementById('sw-mac');
        if (el) el.innerHTML = makeTable(['MAC Address','Port','VLAN','Device','Last Seen'], r, ['mac','port','vlan','device','lastSeen']);
      });
    },

    'desktop': function() {
      // Matrix rain
      var bg = document.getElementById('matrix-bg');
      if (!bg) return;
      var chars = 'アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン0123456789ABCDEF';
      var cols = Math.floor(window.innerWidth / 14);
      var drops = [];
      for (var i = 0; i < cols; i++) drops[i] = Math.random() * -100;

      var canvas = document.createElement('canvas');
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
      canvas.style.cssText = 'width:100%;height:100%';
      bg.appendChild(canvas);
      var ctx = canvas.getContext('2d');

      function drawMatrix() {
        ctx.fillStyle = 'rgba(0,0,0,0.05)';
        ctx.fillRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = '#0f0';
        ctx.font = '14px monospace';
        for (var i = 0; i < drops.length; i++) {
          var ch = chars[Math.floor(Math.random() * chars.length)];
          ctx.fillText(ch, i * 14, drops[i] * 14);
          if (drops[i] * 14 > canvas.height && Math.random() > 0.975) {
            drops[i] = 0;
          }
          drops[i]++;
        }
      }
      setInterval(drawMatrix, 50);

      // Clock
      function updateClock() {
        var el = document.getElementById('taskbar-clock');
        if (el) {
          var now = new Date();
          el.textContent = now.toLocaleTimeString([], {hour:'2-digit',minute:'2-digit',second:'2-digit'});
        }
      }
      updateClock();
      setInterval(updateClock, 1000);

      // Hover effects
      document.querySelectorAll('.desk-icon').forEach(function(el) {
        el.addEventListener('mouseenter', function() {
          if (el.style.opacity !== '0.3') {
            el.style.background = 'rgba(0,255,0,0.08)';
            el.style.boxShadow = '0 0 20px rgba(0,255,0,0.15)';
          }
        });
        el.addEventListener('mouseleave', function() {
          el.style.background = '';
          el.style.boxShadow = '';
        });
      });
    }
  };

  if (renderers[pageType]) {
    renderers[pageType]();
  }
});
