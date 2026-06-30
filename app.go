package main

import (
	"net/http"

	"chefscript/engine"
)

// RegisterApp registers all components, pages, and actions for this application.
func RegisterApp(e *engine.Engine) {
	// ── Components ──
	e.Register("backrooms", engine.BackroomsComponent())

	// ── Pages ──
	engine.RegisterPage("start", func(r *http.Request) *engine.PageContext {
		return engine.NewPageContext()
	})
	engine.RegisterPage("backrooms", func(r *http.Request) *engine.PageContext {
		return engine.NewPageContext()
	})

	// ── Framework Actions ──
	engine.RegisterFlightRecorder()
	engine.RegisterScreenshot()
}
