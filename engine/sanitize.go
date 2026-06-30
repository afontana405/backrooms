package engine

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxFieldLength = 1024

// SanitizeMiddleware reads POST/PUT/PATCH JSON bodies, cleans all string
// fields, then replaces the request body so handlers see clean input.
func SanitizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			ct := r.Header.Get("Content-Type")
			if strings.Contains(ct, "application/json") {
				raw, err := io.ReadAll(io.LimitReader(r.Body, 50<<20)) // 50 MB max (file uploads)
				r.Body.Close()
				if err == nil && len(raw) > 0 {
					cleaned := sanitizeJSON(raw)
					r.Body = io.NopCloser(bytes.NewReader(cleaned))
					r.ContentLength = int64(len(cleaned))
				} else {
					r.Body = io.NopCloser(bytes.NewReader(raw))
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// sanitizeJSON walks a decoded JSON value and cleans all strings.
func sanitizeJSON(raw []byte) []byte {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return raw // not valid JSON — pass through, handler will reject
	}
	cleaned := cleanValue(v)
	out, err := json.Marshal(cleaned)
	if err != nil {
		return raw
	}
	return out
}

func cleanValue(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return cleanString(val)
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, v2 := range val {
			if k == "file" {
				out[k] = v2 // binary data — pass through unsanitized
			} else {
				out[k] = cleanValue(v2)
			}
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(val))
		for i, v2 := range val {
			out[i] = cleanValue(v2)
		}
		return out
	default:
		return val
	}
}

// cleanString applies all sanitization rules to a single string value.
func cleanString(s string) string {
	// Reject / strip null bytes
	s = strings.ReplaceAll(s, "\x00", "")

	// Strip HTML tags
	s = stripHTML(s)

	// Normalize unicode (remove non-printable, non-graphic runes)
	s = strings.Map(func(r rune) rune {
		if r == utf8.RuneError {
			return -1
		}
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)

	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)

	// Enforce max length
	if len(s) > maxFieldLength {
		s = s[:maxFieldLength]
	}

	return s
}

// stripHTML removes all HTML/XML tags from a string.
func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return b.String()
}
