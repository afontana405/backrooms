package engine

import (
	"encoding/json"
	"net/http"
)

// ActionResult is the JSON response from every API handler
type ActionResult struct {
	Redirect     string                 `json:"redirect,omitempty"`
	Error        string                 `json:"error,omitempty"`
	Fields       map[string]string      `json:"fields,omitempty"` // field-level errors: fieldName → message
	Data         interface{} `json:"data,omitempty"`
	Toast        string                 `json:"toast,omitempty"`
	ToastVariant string                 `json:"toastVariant,omitempty"` // info | success | warning | error
}

// APIHandler processes a POST request and returns an ActionResult
type APIHandler func(w http.ResponseWriter, r *http.Request) ActionResult

var apiHandlers = map[string]APIHandler{}

// RegisterAction registers a handler for a named API action.
// Actions are called via POST /api/{name} from the JS runtime.
func RegisterAction(name string, handler APIHandler) {
	apiHandlers[name] = handler
}

// GetAction returns the registered handler for an action, or nil
func GetAction(name string) APIHandler {
	return apiHandlers[name]
}

// DecodeBody reads and decodes the JSON request body into the target struct.
func DecodeBody(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
