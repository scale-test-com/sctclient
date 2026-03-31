package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"scale-test/cli/internal/model"
)

func TestNewClientTrimsTrailingSlash(t *testing.T) {
	c := NewClient("http://example.test/api/v1/", "k")
	if c.baseURL != "http://example.test/api/v1" {
		t.Fatalf("unexpected baseURL: %s", c.baseURL)
	}
}

func TestCreateRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/run/new" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer key-123" {
			t.Fatalf("unexpected auth header: %s", got)
		}
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
			t.Fatalf("unexpected content-type: %s", got)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"run-1","state":"pending","created_at":"2026-01-01T00:00:00Z","scenario_id":123,"message":"ok"}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key-123")
	sid := 123
	resp, err := c.CreateRun(model.CreateRunRequest{ScenarioID: &sid})
	if err != nil {
		t.Fatalf("CreateRun failed: %v", err)
	}
	if resp.ID != "run-1" || resp.State != "pending" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/runs/abc" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":{"id":"abc","created_at":"2026-01-01T00:00:00Z","state":"completed"}}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key")
	run, err := c.GetRun("abc")
	if err != nil {
		t.Fatalf("GetRun failed: %v", err)
	}
	if run.ID != "abc" || run.State != "completed" {
		t.Fatalf("unexpected run: %+v", run)
	}
}

func TestDeleteRunSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/runs/abc" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"message":"deleted"}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key")
	resp, err := c.DeleteRun("abc")
	if err != nil {
		t.Fatalf("DeleteRun failed: %v", err)
	}
	if resp.Message != "deleted" {
		t.Fatalf("unexpected message: %s", resp.Message)
	}
}

func TestDoReturnsJSONApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Invalid API key"}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key")
	_, err := c.GetRun("x")
	if err == nil || !strings.Contains(err.Error(), "Invalid API key") {
		t.Fatalf("expected API error with message, got: %v", err)
	}
}

func TestDoReturnsRawBodyOnNonJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key")
	_, err := c.GetRun("x")
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected raw-body API error, got: %v", err)
	}
}

func TestDecodeResponseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key")
	_, err := c.GetRun("x")
	if err == nil || !strings.Contains(err.Error(), "decode response") {
		t.Fatalf("expected decode response error, got: %v", err)
	}
}

func TestNewRequestErrorWithInvalidBaseURL(t *testing.T) {
	c := NewClient("http://[::1", "key")
	_, err := c.newRequest(http.MethodGet, "/runs/x", nil)
	if err == nil {
		t.Fatal("expected newRequest error")
	}
}
