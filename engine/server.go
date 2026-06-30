package engine

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
)

// ServerPort is the local port the HTTP server listens on
const ServerPort = 7070

// DefaultPage is the page name the root "/" redirects to.
// If empty, the server scans the pages directory for the first available page.
var DefaultPage string

// LANIP is the detected LAN IP address of this machine.
var LANIP string

// GetLANIP returns the first non-loopback IPv4 address found.
func GetLANIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			return ip.String()
		}
	}
	return "127.0.0.1"
}

// StartServer starts the ChefScript HTTP server.
// pagesDir is the directory containing JSON page templates.
func StartServer(e *Engine, pagesDir string) error {
	// Detect LAN IP
	LANIP = GetLANIP()

	// Register server info API
	RegisterAction("server/info", APIHandler(func(w http.ResponseWriter, r *http.Request) ActionResult {
		return ActionResult{
			Data: map[string]interface{}{
				"ip":   LANIP,
				"port": ServerPort,
				"url":  fmt.Sprintf("http://%s:%d", LANIP, ServerPort),
			},
		}
	}))

	mux := http.NewServeMux()

	// Static files — serve from embedded FS if available, else disk
	if EmbeddedFS != nil {
		sub, _ := fs.Sub(EmbeddedFS, "public")
		mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.FS(sub))))
	} else {
		mux.Handle("/public/", http.StripPrefix("/public/",
			http.FileServer(http.Dir(e.PublicPath))))
	}

	// Page routes — GET /page/:name
	mux.Handle("/page/", SessionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/page/")
		// Strip any trailing path segments (e.g. /page/quiz/1 → use name=quiz, param=1)
		parts := strings.SplitN(name, "/", 2)
		pageName := parts[0]
		if pageName == "" {
			pageName = "index"
		}
		servePageHandler(e, pagesDir, pageName, w, r)
	})))

	// Partial routes — GET /partial/:name (body content only, no shell)
	mux.Handle("/partial/", SessionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/partial/")
		parts := strings.SplitN(name, "/", 2)
		pageName := parts[0]
		if pageName == "" {
			pageName = "index"
		}
		servePartialHandler(e, pagesDir, pageName, w, r)
	})))

	// API routes — POST /api/:action (session → sanitize → handler)
	mux.Handle("/api/", SessionMiddleware(SanitizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action := strings.TrimPrefix(r.URL.Path, "/api/")
		serveAPIHandler(w, r, action)
	}))))

	// Root redirect — use DefaultPage if set, otherwise first page found
	landing := DefaultPage
	if landing == "" {
		landing = discoverFirstPage(pagesDir)
	}
	if landing == "" {
		landing = "index"
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/page/"+landing, http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	addr := fmt.Sprintf(":%d", ServerPort)
	fmt.Printf("ChefScript server running at http://localhost%s\n", addr)
	fmt.Printf("LAN address: http://%s:%d\n", LANIP, ServerPort)
	return http.ListenAndServe(addr, mux)
}

// discoverFirstPage scans the pages directory and returns the first .json filename (without extension).
func discoverFirstPage(pagesDir string) string {
	entries, err := ReadEmbedDir(pagesDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(name, ".json") && name != "login.json" {
			return strings.TrimSuffix(name, ".json")
		}
	}
	return ""
}

func servePageHandler(e *Engine, pagesDir, name string, w http.ResponseWriter, r *http.Request) {
	sid := sessionIDFromRequest(r)
	GlobalFlight.Record("server", "page", DiagInfo, "page:load "+name, r.URL.String(), sid)

	path := pagesDir + "/" + name + ".json"
	raw, err := ReadEmbedFile(path)
	if err != nil {
		GlobalFlight.Record("server", "error", DiagError, "page:notfound "+name, "", sid)
		http.Error(w, "Page not found: "+name, http.StatusNotFound)
		return
	}

	// Build base context from session
	ctx := buildPageContext(r)

	// Run registered page loader to inject live data
	if loader := GetPageLoader(name); loader != nil {
		loaded := loader(r)
		if loaded != nil {
			if loaded.Redirect != "" {
				http.Redirect(w, r, loaded.Redirect, http.StatusFound)
				return
			}
			if loaded.Data != nil {
				ctx.Data = loaded.Data
			}
			if loaded.Flash != "" {
				ctx.Flash = loaded.Flash
			}
		}
	}

	// Create diagnostics collector for this request
	diag := NewDiagCollector()

	// Substitute template variables
	bound, err := ApplyContextDiag(raw, ctx, diag)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render JSON → HTML (with diagnostics)
	html, err := e.RenderDiag(bound, diag)
	if err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func servePartialHandler(e *Engine, pagesDir, name string, w http.ResponseWriter, r *http.Request) {
	sid := sessionIDFromRequest(r)
	GlobalFlight.Record("server", "page", DiagInfo, "partial:load "+name, r.URL.String(), sid)

	path := pagesDir + "/" + name + ".json"
	raw, err := ReadEmbedFile(path)
	if err != nil {
		GlobalFlight.Record("server", "error", DiagError, "partial:notfound "+name, "", sid)
		http.Error(w, "Page not found: "+name, http.StatusNotFound)
		return
	}

	ctx := buildPageContext(r)
	if loader := GetPageLoader(name); loader != nil {
		loaded := loader(r)
		if loaded != nil {
			if loaded.Redirect != "" {
				// Return redirect as JSON for client-side handling
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"redirect": loaded.Redirect})
				return
			}
			if loaded.Data != nil {
				ctx.Data = loaded.Data
			}
		}
	}

	bound, err := ApplyContext(raw, ctx)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	body, title, err := e.RenderPartial(bound)
	if err != nil {
		http.Error(w, "Render error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"html":  body,
		"title": title,
		"url":   "/page/" + name + "?" + r.URL.RawQuery,
	})
}

func serveAPIHandler(w http.ResponseWriter, r *http.Request, action string) {
	w.Header().Set("Content-Type", "application/json")
	sid := sessionIDFromRequest(r)

	handler := GetAction(action)
	if handler == nil {
		GlobalFlight.Record("server", "action", DiagError, "action:notfound "+action, "", sid)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ActionResult{Error: "unknown action: " + action})
		return
	}

	GlobalFlight.Record("server", "action", DiagInfo, "action:enter "+action, "", sid)
	result := handler(w, r)

	level := DiagInfo
	detail := ""
	if result.Error != "" {
		level = DiagError
		detail = result.Error
	} else if result.Redirect != "" {
		detail = "redirect:" + result.Redirect
	} else if result.Toast != "" {
		detail = "toast:" + result.Toast
	} else if result.Data != nil {
		detail = "data:" + stripBase64FromDetail(result.Data)
	}
	GlobalFlight.Record("server", "action", level, "action:exit "+action, detail, sid)

	json.NewEncoder(w).Encode(result)
}

func buildPageContext(r *http.Request) *PageContext {
	ctx := NewPageContext()

	sess := GetSessionFromCtx(r)
	if sess == nil {
		return ctx
	}

	ctx.Session = sess.Data

	if id, ok := sess.Data["userId"].(string); ok {
		ctx.User["id"] = id
	}
	if v, ok := sess.Data["hackerName"].(string); ok {
		ctx.User["hackerName"] = v
	}
	if v, ok := sess.Data["firstName"].(string); ok {
		ctx.User["firstName"] = v
	}
	if v, ok := sess.Data["lastName"].(string); ok {
		ctx.User["lastName"] = v
	}
	if v, ok := sess.Data["role"].(string); ok {
		ctx.User["role"] = v
	}

	return ctx
}

// stripBase64FromDetail redacts base64 file content from flight recorder details.
// Keeps filenames and metadata, replaces large binary strings with a length summary.
func stripBase64FromDetail(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		clean := make(map[string]interface{}, len(v))
		for k, val := range v {
			if k == "file" || k == "download" || k == "content" {
				if s, ok := val.(string); ok {
					clean[k] = fmt.Sprintf("[base64 %d chars]", len(s))
				} else {
					clean[k] = val
				}
			} else {
				clean[k] = val
			}
		}
		b, err := json.Marshal(clean)
		if err != nil {
			return "{}"
		}
		d := string(b)
		if len(d) > 500 {
			d = d[:500] + "..."
		}
		return d
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return "{}"
		}
		d := string(b)
		if len(d) > 500 {
			d = d[:500] + "..."
		}
		return d
	}
}
