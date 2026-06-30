package engine

import "fmt"

// ── Video ───────────────────────────────────────────────────────────────
// ["video", { "src": "/public/intro.mp4", "poster": "/public/thumb.jpg", "controls": true }]

func renderVideo(props map[string]interface{}, children string, e *Engine) (string, error) {
	src := propStr(props, "src", "")
	poster := propStr(props, "poster", "")
	controls := propBool(props, "controls", true)
	autoplay := propBool(props, "autoplay", false)
	muted := propBool(props, "muted", false)
	loop := propBool(props, "loop", false)
	dataID := propStr(props, "data-id", "video")

	attrs := ""
	if controls {
		attrs += " controls"
	}
	if autoplay {
		attrs += " autoplay"
	}
	if muted {
		attrs += " muted"
	}
	if loop {
		attrs += " loop"
	}

	posterAttr := ""
	if poster != "" {
		posterAttr = fmt.Sprintf(` poster="%s"`, poster)
	}

	return fmt.Sprintf(`<video class="cs-video" src="%s"%s%s data-id="%s"></video>`,
		src, posterAttr, attrs, dataID), nil
}

// ── Audio ───────────────────────────────────────────────────────────────
// ["audio", { "src": "/public/track.mp3", "controls": true }]

func renderAudio(props map[string]interface{}, children string, e *Engine) (string, error) {
	src := propStr(props, "src", "")
	controls := propBool(props, "controls", true)
	autoplay := propBool(props, "autoplay", false)
	loop := propBool(props, "loop", false)
	dataID := propStr(props, "data-id", "audio")

	attrs := ""
	if controls {
		attrs += " controls"
	}
	if autoplay {
		attrs += " autoplay"
	}
	if loop {
		attrs += " loop"
	}

	return fmt.Sprintf(`<audio class="cs-audio" src="%s"%s data-id="%s"></audio>`,
		src, attrs, dataID), nil
}

// ── Iframe ──────────────────────────────────────────────────────────────
// ["iframe", { "src": "https://...", "height": "400", "title": "Map" }]

func renderIframe(props map[string]interface{}, children string, e *Engine) (string, error) {
	src := propStr(props, "src", "")
	height := propStr(props, "height", "400")
	title := propStr(props, "title", "")
	allow := propStr(props, "allow", "")
	sandbox := propStr(props, "sandbox", "")
	dataID := propStr(props, "data-id", "iframe")

	allowAttr := ""
	if allow != "" {
		allowAttr = fmt.Sprintf(` allow="%s"`, allow)
	}
	sandboxAttr := ""
	if sandbox != "" {
		sandboxAttr = fmt.Sprintf(` sandbox="%s"`, sandbox)
	}

	return fmt.Sprintf(`<iframe class="cs-iframe" src="%s" height="%s" title="%s"%s%s data-id="%s" frameborder="0" loading="lazy"></iframe>`,
		src, height, title, allowAttr, sandboxAttr, dataID), nil
}

// ── AspectRatio ─────────────────────────────────────────────────────────
// ["aspect-ratio", { "ratio": "16/9" }, ["video", { "src": "..." }]]

func renderAspectRatio(props map[string]interface{}, children string, e *Engine) (string, error) {
	ratio := propStr(props, "ratio", "16/9")
	dataID := propStr(props, "data-id", "aspect-ratio")

	return fmt.Sprintf(`<div class="cs-aspect-ratio" style="aspect-ratio:%s" data-id="%s">%s</div>`,
		ratio, dataID, children), nil
}

// ── Carousel ────────────────────────────────────────────────────────────
// ["carousel", { "id": "hero" }, ...slides]

func renderCarousel(props map[string]interface{}, children string, e *Engine) (string, error) {
	id := propStr(props, "id", "carousel")
	dataID := propStr(props, "data-id", "carousel--"+id)
	trackID := id + "--track"

	prevOnclick := fmt.Sprintf(
		`document.getElementById('%s').scrollBy({left:-document.getElementById('%s').offsetWidth,behavior:'smooth'})`,
		trackID, trackID)
	nextOnclick := fmt.Sprintf(
		`document.getElementById('%s').scrollBy({left:document.getElementById('%s').offsetWidth,behavior:'smooth'})`,
		trackID, trackID)

	return fmt.Sprintf(`<div class="cs-carousel" id="%s" data-id="%s">
  <div class="cs-carousel__track" id="%s">%s</div>
  <button class="cs-carousel__btn cs-carousel__btn--prev" type="button" data-id="%s--prev"
    onclick="%s">&#8249;</button>
  <button class="cs-carousel__btn cs-carousel__btn--next" type="button" data-id="%s--next"
    onclick="%s">&#8250;</button>
</div>`, id, dataID, trackID, children, dataID, prevOnclick, dataID, nextOnclick), nil
}

// ── RichText ────────────────────────────────────────────────────────────
// ["rich-text", { "content": "<p>Pre-rendered HTML</p>" }]

func renderRichText(props map[string]interface{}, children string, e *Engine) (string, error) {
	content := propStr(props, "content", children)
	dataID := propStr(props, "data-id", "rich-text")

	return fmt.Sprintf(`<div class="cs-rich-text" data-id="%s">%s</div>`, dataID, content), nil
}
