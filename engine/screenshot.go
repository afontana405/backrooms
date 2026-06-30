package engine

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

// WindowHandle stores the HWND of the main webview window.
// Set from main.go after webview creation.
var WindowHandle uintptr

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	gdi32           = syscall.NewLazyDLL("gdi32.dll")
	pGetClientRect  = user32.NewProc("GetClientRect")
	pGetDC          = user32.NewProc("GetDC")
	pReleaseDC      = user32.NewProc("ReleaseDC")
	pPrintWindow    = user32.NewProc("PrintWindow")
	pCreateCompatDC = gdi32.NewProc("CreateCompatibleDC")
	pCreateCompatBM = gdi32.NewProc("CreateCompatibleBitmap")
	pSelectObject   = gdi32.NewProc("SelectObject")
	pDeleteObject   = gdi32.NewProc("DeleteObject")
	pDeleteDC       = gdi32.NewProc("DeleteDC")
	pGetDIBits      = gdi32.NewProc("GetDIBits")
)

type rect struct {
	Left, Top, Right, Bottom int32
}

type bitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

// RegisterScreenshot registers the screenshot/save action.
func RegisterScreenshot() {
	RegisterAction("screenshot/save", screenshotSaveHandler())
}

func screenshotSaveHandler() APIHandler {
	return func(w http.ResponseWriter, r *http.Request) ActionResult {
		outPath, err := CaptureWindow()
		if err != nil {
			return ActionResult{Error: "Screenshot failed: " + err.Error()}
		}
		return ActionResult{
			Toast: "Screenshot saved: " + outPath,
		}
	}
}

// CaptureWindow captures the webview window as a PNG and saves it to the tmp folder.
func CaptureWindow() (string, error) {
	hwnd := WindowHandle
	if hwnd == 0 {
		return "", fmt.Errorf("window handle not set")
	}

	// Get client area dimensions
	var rc rect
	ret, _, _ := pGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rc)))
	if ret == 0 {
		return "", fmt.Errorf("GetClientRect failed")
	}
	w := int(rc.Right - rc.Left)
	h := int(rc.Bottom - rc.Top)
	if w == 0 || h == 0 {
		return "", fmt.Errorf("window has zero size")
	}

	// Get window DC
	hdcWindow, _, _ := pGetDC.Call(hwnd)
	if hdcWindow == 0 {
		return "", fmt.Errorf("GetDC failed")
	}
	defer pReleaseDC.Call(hwnd, hdcWindow)

	// Create compatible DC and bitmap
	hdcMem, _, _ := pCreateCompatDC.Call(hdcWindow)
	if hdcMem == 0 {
		return "", fmt.Errorf("CreateCompatibleDC failed")
	}
	defer pDeleteDC.Call(hdcMem)

	hBitmap, _, _ := pCreateCompatBM.Call(hdcWindow, uintptr(w), uintptr(h))
	if hBitmap == 0 {
		return "", fmt.Errorf("CreateCompatibleBitmap failed")
	}
	defer pDeleteObject.Call(hBitmap)

	// Select bitmap into memory DC
	pSelectObject.Call(hdcMem, hBitmap)

	// PrintWindow with PW_RENDERFULLCONTENT (flag 2) — captures hardware-accelerated content
	ret, _, _ = pPrintWindow.Call(hwnd, hdcMem, 2)
	if ret == 0 {
		return "", fmt.Errorf("PrintWindow failed")
	}

	// Read bitmap pixels via GetDIBits
	bmi := bitmapInfoHeader{
		BiSize:     uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		BiWidth:    int32(w),
		BiHeight:   -int32(h), // negative = top-down
		BiPlanes:   1,
		BiBitCount: 32,
	}

	pixels := make([]byte, w*h*4)
	ret, _, _ = pGetDIBits.Call(hdcMem, hBitmap, 0, uintptr(h), uintptr(unsafe.Pointer(&pixels[0])), uintptr(unsafe.Pointer(&bmi)), 0)
	if ret == 0 {
		return "", fmt.Errorf("GetDIBits failed")
	}

	// Convert BGRA → RGBA
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < len(pixels); i += 4 {
		img.Pix[i+0] = pixels[i+2] // R ← B
		img.Pix[i+1] = pixels[i+1] // G
		img.Pix[i+2] = pixels[i+0] // B ← R
		img.Pix[i+3] = 255         // A (fully opaque)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("PNG encode failed: %v", err)
	}

	// Save to tmp folder in project root (go up from exe dir which may be tmp/)
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)
	dir := exeDir
	if filepath.Base(exeDir) == "tmp" {
		dir = filepath.Dir(exeDir)
	}
	dir = filepath.Join(dir, "tmp")
	os.MkdirAll(dir, 0755)
	filename := fmt.Sprintf("screenshot_%s.png", time.Now().Format("20060102_150405"))
	outPath := filepath.Join(dir, filename)
	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("write failed: %v", err)
	}

	return outPath, nil
}

