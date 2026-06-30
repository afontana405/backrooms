package engine

import (
	_ "embed"
	"strings"
)

// JS source files — each is a self-contained section of the runtime.
// Core infrastructure (split into 6 concerns):

//go:embed js/core_helpers.js
var jsCoreHelpers string

//go:embed js/core_components.js
var jsCoreComponents string

//go:embed js/core_forms.js
var jsCoreForms string

//go:embed js/core_nav.js
var jsCoreNav string

//go:embed js/core_data.js
var jsCoreData string

//go:embed js/core_desktop.js
var jsCoreDesktop string

// State management:

//go:embed js/state.js
var jsState string

// Targeting inspector:

//go:embed js/targeting.js
var jsTargeting string

// csRuntime is the JavaScript runtime injected into every page.
// It powers all interactive component behaviors without any imports.
// Built from engine/js/*.js files via go:embed at compile time.
// Add new app-specific JS files here as your app grows.
var csRuntime = buildRuntime()

func buildRuntime() string {
	var b strings.Builder
	b.WriteString("<script>\n(function(){\n'use strict';\n")
	// Core infrastructure (order matters — helpers first)
	b.WriteString(jsCoreHelpers)
	b.WriteString("\n")
	b.WriteString(jsCoreComponents)
	b.WriteString("\n")
	b.WriteString(jsCoreForms)
	b.WriteString("\n")
	b.WriteString(jsCoreNav)
	b.WriteString("\n")
	b.WriteString(jsCoreData)
	b.WriteString("\n")
	b.WriteString(jsCoreDesktop)
	b.WriteString("\n")
	// State
	b.WriteString(jsState)
	b.WriteString("\n")
	// Targeting inspector
	b.WriteString(jsTargeting)
	b.WriteString("\n")
	// App-specific JS goes here:
	// b.WriteString(jsMyGame)
	// b.WriteString("\n")
	b.WriteString("\n})();\n</script>\n")
	return b.String()
}
