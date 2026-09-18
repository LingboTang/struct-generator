// Package api exposes struct-generator's functionality over HTTP.
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LingboTang/struct-generator/structgen"
)

// maxRequestBodyBytes caps the size of a request body to guard against
// unbounded reads from misbehaving or malicious clients.
const maxRequestBodyBytes = 10 << 20 // 10 MiB

// generateRequest is the expected shape of a POST /generate request body:
//
//	{"payload": "<the actual JSON to generate structs from>"}
type generateRequest struct {
	Payload string `json:"payload"`
	Package string `json:"package,omitempty"`
}

type generateResponse struct {
	Result string `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// NewMux builds the HTTP routes for the struct-generator API.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("POST /generate", generateHandler)
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "decoding request body: "+err.Error())
		return
	}

	if strings.TrimSpace(req.Payload) == "" {
		writeError(w, http.StatusBadRequest, "payload must not be empty")
		return
	}

	packageName := req.Package
	if packageName == "" {
		packageName = "main"
	}

	src, err := structgen.Generate(strings.NewReader(req.Payload), packageName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, generateResponse{Result: src})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
