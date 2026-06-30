package engine

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

const sessionCookieName = "cs_session"

// Session is stored in SQLite
type Session struct {
	ID        string
	Data      map[string]interface{}
	ExpiresAt time.Time
}

func newSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// EnsureSessionTTL creates the sessions table if it doesn't exist.
func EnsureSessionTTL() {
	if SqlDB == nil {
		return
	}
	SqlDB.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		data TEXT,
		expires_at INTEGER
	)`)
}

// GetSession reads the session from SQLite using the request cookie
func GetSession(r *http.Request) (*Session, error) {
	if SqlDB == nil {
		return nil, nil
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, nil
	}

	var dataJSON string
	var expiresUnix int64
	err = SqlDB.QueryRow(
		`SELECT data, expires_at FROM sessions WHERE id = ? AND expires_at > ?`,
		cookie.Value, time.Now().Unix(),
	).Scan(&dataJSON, &expiresUnix)
	if err != nil {
		return nil, nil
	}

	sess := &Session{
		ID:        cookie.Value,
		ExpiresAt: time.Unix(expiresUnix, 0),
	}
	json.Unmarshal([]byte(dataJSON), &sess.Data)
	return sess, nil
}

// CreateSession creates a new session in SQLite and sets the cookie
func CreateSession(w http.ResponseWriter, data map[string]interface{}) (*Session, error) {
	sess := &Session{
		ID:        newSessionID(),
		Data:      data,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	dataJSON, _ := json.Marshal(data)
	_, err := SqlDB.Exec(
		`INSERT INTO sessions (id, data, expires_at) VALUES (?, ?, ?)`,
		sess.ID, string(dataJSON), sess.ExpiresAt.Unix(),
	)
	if err != nil {
		return nil, err
	}
	GlobalFlight.Record("server", "session", DiagInfo, "session:create", sess.ID, sess.ID)

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sess.ID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400,
	})
	return sess, nil
}

// UpdateSession saves updated data back to an existing session
func UpdateSession(sessID string, data map[string]interface{}) error {
	dataJSON, _ := json.Marshal(data)
	_, err := SqlDB.Exec(
		`UPDATE sessions SET data = ? WHERE id = ?`,
		string(dataJSON), sessID,
	)
	if err == nil {
		GlobalFlight.Record("server", "session", DiagInfo, "session:update", sessID, sessID)
	}
	return err
}

// DestroySession deletes the session and clears the cookie
func DestroySession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return
	}

	SqlDB.Exec(`DELETE FROM sessions WHERE id = ?`, cookie.Value)
	GlobalFlight.Record("server", "session", DiagInfo, "session:destroy", cookie.Value, cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:   sessionCookieName,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
}

// ClearAllSessions removes all sessions from SQLite.
func ClearAllSessions() {
	if SqlDB == nil {
		return
	}
	SqlDB.Exec(`DELETE FROM sessions`)
}
