package engine

import "fmt"

// ── ChatWidget ──────────────────────────────────────────────────────────
// ["chat-widget", { "webhook": "https://...", "route": "general", "title": "Chat" }]

func renderChatWidget(props map[string]interface{}, children string, e *Engine) (string, error) {
	webhook := propStr(props, "webhook", "")
	route := propStr(props, "route", "general")
	title := propStr(props, "title", "Chat")
	dataID := propStr(props, "data-id", "chat-widget")

	if webhook == "" {
		return "", fmt.Errorf("chat-widget requires a webhook prop")
	}

	return fmt.Sprintf(`
<button
  id="%s--bubble"
  class="cs-chat-bubble"
  aria-label="Open chat"
  onclick="csChatOpen('%s')">💬</button>

<div
  id="%s--container"
  class="cs-chat-container"
  data-webhook="%s"
  data-route="%s"
  data-id="%s">

  <div class="cs-chat-header">
    <span>%s</span>
    <button class="cs-chat-close" onclick="csChatClose('%s')" aria-label="Close chat">✕</button>
  </div>

  <div id="%s--body" class="cs-chat-body">
    <div class="cs-chat-msg cs-chat-msg--bot">
      <strong>Hi 👋</strong> — how can I help?
    </div>
  </div>

  <div class="cs-chat-footer">
    <input
      id="%s--input"
      class="cs-chat-input"
      type="text"
      placeholder="Type a message…"
      onkeydown="if(event.key==='Enter'){event.preventDefault();csChatSend('%s')}"
    />
    <button
      id="%s--send"
      class="cs-chat-send"
      onclick="csChatSend('%s')">Send</button>
  </div>
</div>
`, dataID, dataID,
		dataID, webhook, route, dataID,
		title,
		dataID,
		dataID,
		dataID, dataID,
		dataID, dataID), nil
}

// ── DataChat ────────────────────────────────────────────────────────────
// ["data-chat", { "schema": "schemas/tickets.json", "placeholder": "Ask about tickets..." }]

func renderDataChat(props map[string]interface{}, children string, e *Engine) (string, error) {
	schema := propStr(props, "schema", "")
	placeholder := propStr(props, "placeholder", "Ask a question about the data...")
	id := propStr(props, "id", "data-chat")

	return fmt.Sprintf(`<div%s>
  <div class="cs-card">
    <div style="display:flex;align-items:center;gap:8px;padding-bottom:8px;border-bottom:1px solid var(--border-subtle);margin-bottom:8px">
      <i class="bi bi-chat-dots"></i> <strong>Data Query</strong>
    </div>
    <div id="%s-messages" style="height:300px;overflow-y:auto;font-size:0.9rem">
      <div style="color:var(--text-secondary);font-size:12px">Ask questions in plain English. Answers come from the data, not AI guesswork.</div>
    </div>
    <div style="display:flex;gap:8px;padding-top:8px;border-top:1px solid var(--border-subtle);margin-top:8px">
      <input type="text" class="cs-input__field" id="%s-input" placeholder="%s" style="flex:1"
        onkeydown="if(event.key==='Enter'){event.preventDefault();csDataChat('%s','%s')}" />
      <button class="cs-button cs-button--outline cs-button--md" onclick="csDataChat('%s','%s')" type="button">
        <i class="bi bi-send"></i>
      </button>
    </div>
  </div>
</div>`, userAttrs(props, ""), id, id, placeholder, id, schema, id, schema), nil
}
