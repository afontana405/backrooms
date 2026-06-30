package engine

import "strings"

// BackroomsComponent renders the first-person Backrooms scene: a full-screen
// container plus the Three.js runtime, then the embedded backrooms.js (which
// builds the walkable structure from its ASCII map). Mirrors GameComponent's
// dynamic-load pattern so it also runs after partial navigation.
func BackroomsComponent() ComponentFunc {
	return func(props map[string]interface{}, children string, e *Engine) (string, error) {
		raw, _ := ReadEmbedFile("public/js/backrooms.js")
		var b strings.Builder
		b.WriteString(`<div id="br-root" class="br-root" data-id="backrooms"></div>`)
		b.WriteString(`<script>(function(){var _r=function(){` + string(raw) + `};if(window.THREE){_r();}else{var s=document.createElement('script');s.src='/public/js/three.min.js';s.onload=_r;document.head.appendChild(s);}})();</script>`)
		b.WriteString(`<style>.br-root{position:fixed;inset:0;background:#000;overflow:hidden}.br-root canvas{display:block}</style>`)
		return b.String(), nil
	}
}
