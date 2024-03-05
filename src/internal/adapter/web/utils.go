package web

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ObjectMessage struct {
	Message     string       `json:"message"`
	ID          string       `json:"id,omitempty"`
	Version     string       `json:"version,omitempty"`
	Information string       `json:"information,omitempty"`
	Document    *DocumentDTO `json:"document,omitempty"`
}

// methodNotAllowed sends a method not allowed response.
func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Method %v is not allowed on %v", r.Method, r.URL.Path)
	http.Error(w, msg, http.StatusMethodNotAllowed)
}

// writeJson marshals and writes a JSON response
func writeJson(w http.ResponseWriter, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "inernal server error", http.StatusInternalServerError)
		slog.Error("handler.printjson", "error", err.Error())
	}
	w.Write(b)
}

// writeError sends an error response in JSON format.
func writeError(w http.ResponseWriter, code int, s string, msg *ObjectMessage) {
	msg.Message = s
	w.WriteHeader(code)
	writeJson(w, msg)
}
