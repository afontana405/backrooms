package main

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"chefscript/engine"

	webview "github.com/jchv/go-webview2"
)

func hideConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsole := kernel32.NewProc("GetConsoleWindow")
	user32 := syscall.NewLazyDLL("user32.dll")
	showWindow := user32.NewProc("ShowWindow")
	hwnd, _, _ := getConsole.Call()
	if hwnd != 0 {
		showWindow.Call(hwnd, 0) // SW_HIDE = 0
	}
}

func main() {
	hideConsole()

	// Set embedded filesystem — all file reads go through it
	engine.SetEmbedFS(embeddedAssets)

	// Connect to SQLite
	if err := engine.ConnectSQLite("backrooms.db"); err != nil {
		log.Fatalf("SQLite connection failed: %v", err)
	}
	log.Println("SQLite connected")
	engine.EnsureSessionTTL()

	// Load binary schemas
	if err := engine.LoadBinarySchemas("schemas/binary"); err != nil {
		log.Printf("Warning: loading schemas: %v", err)
	}

	// Create engine
	e := engine.New()

	// Register app — add your components, pages, and actions here
	RegisterApp(e)

	// Start HTTP server
	go func() {
		if err := engine.StartServer(e, "pages"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	time.Sleep(150 * time.Millisecond)

	url := fmt.Sprintf("http://localhost:%d/page/start", engine.ServerPort)

	w := webview.New(false)
	if w == nil {
		log.Fatal("Failed to create webview")
	}
	defer w.Destroy()

	w.SetTitle("Backrooms")
	w.SetSize(1280, 800, webview.HintNone)
	engine.WindowHandle = uintptr(w.Window())
	w.Navigate(url)
	w.Run()

	// App window closed — clear all sessions
	engine.ClearAllSessions()
}
