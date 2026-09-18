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

// schemaFilename is the filename the generated struct source is served
// under on success.
const schemaFilename = "schema.go"

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

	req.Payload = sanitizePayload(req.Payload)
	if req.Payload == "" {
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

	w.Header().Set("Content-Type", "text/x-go; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+schemaFilename+`"`)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(src))
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

// sanitizePayload trims surrounding whitespace and strips stray C0 control
// characters (e.g. raw \r, \b, \f, \v bytes) from the payload. JSON only
// permits such characters inside string literals when properly escaped
// (e.g. "\\r"); an unescaped raw control byte makes the document invalid and
// would otherwise cause json decoding to fail even though the payload is
// "morally" the same document. \n and \t are left alone since they commonly
// appear as insignificant whitespace between tokens in pretty-printed JSON.
func sanitizePayload(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r < 0x20 && r != '\n' && r != '\t' {
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}
