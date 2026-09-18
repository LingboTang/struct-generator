package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doRequest(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rr := httptest.NewRecorder()
	NewHandler().ServeHTTP(rr, req)
	return rr
}

func TestGenerateHandler_Success(t *testing.T) {
	body := `{"payload": "{\"name\": \"Ada\", \"age\": 36}", "package": "models"}`
	rr := doRequest(t, http.MethodPost, "/generate", body)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if ct := rr.Header().Get("Content-Type"); ct != "text/x-go; charset=utf-8" {
		t.Errorf("unexpected Content-Type: %q", ct)
	}
	if cd := rr.Header().Get("Content-Disposition"); cd != `attachment; filename="schema.go"` {
		t.Errorf("unexpected Content-Disposition: %q", cd)
	}

	result := rr.Body.String()
	for _, want := range []string{"package models", "type Root struct", "Name string", "Age", "float64"} {
		if !strings.Contains(result, want) {
			t.Errorf("expected result to contain %q, got:\n%s", want, result)
		}
	}
}

func TestGenerateHandler_DefaultsPackageToMain(t *testing.T) {
	body := `{"payload": "{\"ok\": true}"}`
	rr := doRequest(t, http.MethodPost, "/generate", body)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "package main") {
		t.Errorf("expected default package main, got:\n%s", rr.Body.String())
	}
}

func TestGenerateHandler_StripsStrayControlCharacters(t *testing.T) {
	// Raw, unescaped control bytes (\r, \b) inside a string literal make the
	// JSON invalid; sanitizePayload should strip them so decoding still succeeds.
	inner := "{\"name\": \"ada\rlovelace\", \"note\": \"line1\bline2\"}"
	reqBody, err := json.Marshal(map[string]string{"payload": inner})
	if err != nil {
		t.Fatalf("marshaling request: %v", err)
	}

	rr := doRequest(t, http.MethodPost, "/generate", string(reqBody))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "Name string") {
		t.Errorf("expected result to contain %q, got:\n%s", "Name string", rr.Body.String())
	}
}

func TestSanitizePayload(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"trims whitespace", "  {\"a\":1}  \n", "{\"a\":1}"},
		{"strips carriage return", "{\"a\":\"x\ry\"}", "{\"a\":\"xy\"}"},
		{"strips backspace", "{\"a\":\"x\by\"}", "{\"a\":\"xy\"}"},
		{"strips form feed and vertical tab", "{\"a\":\"x\f\vy\"}", "{\"a\":\"xy\"}"},
		{"keeps tabs and newlines", "{\n\t\"a\": 1\n}", "{\n\t\"a\": 1\n}"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := sanitizePayload(c.in); got != c.want {
				t.Errorf("sanitizePayload(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestGenerateHandler_EmptyPayload(t *testing.T) {
	rr := doRequest(t, http.MethodPost, "/generate", `{"payload": ""}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestGenerateHandler_MalformedRequestBody(t *testing.T) {
	rr := doRequest(t, http.MethodPost, "/generate", `not json`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestGenerateHandler_InvalidPayloadJSON(t *testing.T) {
	rr := doRequest(t, http.MethodPost, "/generate", `{"payload": "not json"}`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if resp.Error == "" {
		t.Error("expected a non-empty error message")
	}
}

func TestGenerateHandler_WrongMethod(t *testing.T) {
	rr := doRequest(t, http.MethodGet, "/generate", "")
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestHealthHandler(t *testing.T) {
	rr := doRequest(t, http.MethodGet, "/healthz", "")
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
