package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ── Flight Event ────────────────────────────────────────────────────────────

type FlightEvent struct {
	ID        uint64    `json:"id"`
	Time      time.Time `json:"time"`
	Source    string    `json:"src"`              // "server" or "client"
	Category string    `json:"cat"`              // action, command, flag, state, patch, session, page, error, csPost, domPatch, setAttr, event, nav
	Level    DiagLevel `json:"level"`
	Summary  string    `json:"msg"`
	Detail   string    `json:"detail,omitempty"`
	SessionID string   `json:"sid,omitempty"`
}

// ── Ring Buffer ─────────────────────────────────────────────────────────────

type FlightRecorder struct {
	mu     sync.Mutex
	buf    []FlightEvent
	head   int
	count  int
	cap    int
	nextID uint64
}

var GlobalFlight = NewFlightRecorder(500)

func NewFlightRecorder(capacity int) *FlightRecorder {
	return &FlightRecorder{
		buf: make([]FlightEvent, capacity),
		cap: capacity,
	}
}

func (fr *FlightRecorder) Record(src, category string, level DiagLevel, summary, detail, sessionID string) {
	if fr == nil {
		return
	}
	fr.mu.Lock()
	defer fr.mu.Unlock()

	fr.nextID++
	idx := fr.head
	fr.buf[idx] = FlightEvent{
		ID:        fr.nextID,
		Time:      time.Now(),
		Source:    src,
		Category: category,
		Level:    level,
		Summary:  summary,
		Detail:   detail,
		SessionID: sessionID,
	}
	fr.head = (fr.head + 1) % fr.cap
	if fr.count < fr.cap {
		fr.count++
	}
}

func (fr *FlightRecorder) Snapshot() []FlightEvent {
	if fr == nil {
		return nil
	}
	fr.mu.Lock()
	defer fr.mu.Unlock()

	result := make([]FlightEvent, fr.count)
	start := fr.head - fr.count
	if start < 0 {
		start += fr.cap
	}
	for i := 0; i < fr.count; i++ {
		result[i] = fr.buf[(start+i)%fr.cap]
	}
	return result
}

func (fr *FlightRecorder) SnapshotSince(afterID uint64) []FlightEvent {
	all := fr.Snapshot()
	for i, e := range all {
		if e.ID > afterID {
			return all[i:]
		}
	}
	return nil
}

func (fr *FlightRecorder) SnapshotJSON() []byte {
	b, err := json.Marshal(fr.Snapshot())
	if err != nil {
		return []byte("[]")
	}
	return b
}

// ── Session ID helper ───────────────────────────────────────────────────────

func sessionIDFromRequest(r *http.Request) string {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// ── API Endpoints ───────────────────────────────────────────────────────────

func RegisterFlightRecorder() {
	RegisterAction("flight/snapshot", flightSnapshotHandler())
	RegisterAction("flight/push", flightPushHandler())
}

func flightSnapshotHandler() APIHandler {
	return func(w http.ResponseWriter, r *http.Request) ActionResult {
		var body struct {
			AfterID uint64 `json:"afterId"`
		}
		DecodeBody(r, &body)

		var events []FlightEvent
		if body.AfterID > 0 {
			events = GlobalFlight.SnapshotSince(body.AfterID)
		} else {
			events = GlobalFlight.Snapshot()
		}
		if events == nil {
			events = []FlightEvent{}
		}

		return ActionResult{
			Data: map[string]interface{}{
				"events": events,
			},
		}
	}
}

func flightPushHandler() APIHandler {
	return func(w http.ResponseWriter, r *http.Request) ActionResult {
		var body struct {
			Events []struct {
				Category string `json:"cat"`
				Level    int    `json:"level"`
				Summary  string `json:"msg"`
				Detail   string `json:"detail"`
			} `json:"events"`
		}
		if err := DecodeBody(r, &body); err != nil {
			return ActionResult{Error: "invalid request"}
		}

		sid := sessionIDFromRequest(r)
		for _, evt := range body.Events {
			GlobalFlight.Record("client", evt.Category, DiagLevel(evt.Level), evt.Summary, evt.Detail, sid)
		}

		return ActionResult{
			Toast: fmt.Sprintf("Received %d client events", len(body.Events)),
		}
	}
}
