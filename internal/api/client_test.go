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

func TestGetRunResultsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/runs/abc/results" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":{"state":"success","executed_requests":10,"duration_seconds":2,"completed_at":null,"summary":{"total_requests":10,"success_rate":null,"average_ms":null,"requests_per_second":null,"status_codes":[{"status_code":200,"count":10}]}}}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key")
	res, err := c.GetRunResults("abc")
	if err != nil {
		t.Fatalf("GetRunResults failed: %v", err)
	}
	if res.State != "success" || res.Summary.TotalRequests != 10 || res.Summary.SuccessRate != nil || len(res.Summary.StatusCodes) != 1 {
		t.Fatalf("unexpected results: %+v", res)
	}
}

func TestGetRunTimeseriesQuery(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/runs/abc/results/timeseries" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotQuery = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":{"requests_per_second":[{"executed_at":"2026-01-01 10:00:00","requests_per_second":10}]}}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "key")
	series, err := c.GetRunTimeseries("abc", "requests_per_second")
	if err != nil {
		t.Fatalf("GetRunTimeseries failed: %v", err)
	}
	if gotQuery != "metric=requests_per_second" || len(series.RequestsPerSecond) != 1 || series.StatusCodes != nil {
		t.Fatalf("unexpected query %q or series %+v", gotQuery, series)
	}
}
