package engine

import (
	"context"
	"net/http"
)

type contextKey string

const (
	ctxSession contextKey = "session"
	ctxUser    contextKey = "user"
)

// User represents an authenticated user
type User struct {
	ID         string `bson:"_id,omitempty"   json:"id"`
	HackerName string `bson:"hackerName"      json:"hackerName"`
	FirstName  string `bson:"firstName"       json:"firstName"`
	LastName   string `bson:"lastName"        json:"lastName"`
	Role       string `bson:"role"            json:"role"` // student | instructor
	Password   string `bson:"password,omitempty" json:"-"`
}

// SessionMiddleware loads the session (and user) for every request
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := GetSession(r)
		rctx := r.Context()

		if sess != nil {
			rctx = context.WithValue(rctx, ctxSession, sess)

			// Reconstruct a lightweight User from session data
			user := userFromSession(sess)
			if user != nil {
				rctx = context.WithValue(rctx, ctxUser, user)
			}
		}

		next.ServeHTTP(w, r.WithContext(rctx))
	})
}

func userFromSession(sess *Session) *User {
	id, _ := sess.Data["userId"].(string)
	if id == "" {
		return nil
	}
	u := &User{ID: id}
	u.HackerName, _ = sess.Data["hackerName"].(string)
	u.FirstName, _ = sess.Data["firstName"].(string)
	u.LastName, _ = sess.Data["lastName"].(string)
	u.Role, _ = sess.Data["role"].(string)
	return u
}

// RequireAuth redirects unauthenticated users to the login page
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if GetSessionFromCtx(r) == nil {
			http.Redirect(w, r, "/page/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// RequireInstructor redirects non-instructors to login
func RequireInstructor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromCtx(r)
		if user == nil || user.Role != "instructor" {
			http.Redirect(w, r, "/page/login", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// GetSessionFromCtx retrieves the session from request context
func GetSessionFromCtx(r *http.Request) *Session {
	sess, _ := r.Context().Value(ctxSession).(*Session)
	return sess
}

// GetUserFromCtx retrieves the user from request context
func GetUserFromCtx(r *http.Request) *User {
	user, _ := r.Context().Value(ctxUser).(*User)
	return user
}
