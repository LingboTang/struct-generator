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

	var resp generateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	for _, want := range []string{"package models", "type Root struct", "Name string", "Age", "float64"} {
		if !strings.Contains(resp.Result, want) {
			t.Errorf("expected result to contain %q, got:\n%s", want, resp.Result)
		}
	}
}

func TestGenerateHandler_DefaultsPackageToMain(t *testing.T) {
	body := `{"payload": "{\"ok\": true}"}`
	rr := doRequest(t, http.MethodPost, "/generate", body)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp generateResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if !strings.Contains(resp.Result, "package main") {
		t.Errorf("expected default package main, got:\n%s", resp.Result)
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
