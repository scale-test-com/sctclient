package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"scale-test/cli/internal/api"
)

func captureOutput(t *testing.T, fn func() error) (stdout string, stderr string, err error) {
	t.Helper()

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stderr: %v", err)
	}

	os.Stdout = wOut
	os.Stderr = wErr

	runErr := fn()

	_ = wOut.Close()
	_ = wErr.Close()

	outBytes, _ := io.ReadAll(rOut)
	errBytes, _ := io.ReadAll(rErr)
	_ = rOut.Close()
	_ = rErr.Close()

	return string(outBytes), string(errBytes), runErr
}

func resetGlobals(t *testing.T) {
	t.Helper()
	createScenarioID = 0
	createFile = ""
	createWait = false
	createPollInterval = 2 * time.Second
	apiKey = ""
	baseURL = ""
	client = nil
}

func TestPrintJSON(t *testing.T) {
	resetGlobals(t)
	out, _, err := captureOutput(t, func() error {
		return printJSON(map[string]any{"ok": true})
	})
	if err != nil {
		t.Fatalf("printJSON error: %v", err)
	}
	if !strings.Contains(out, "\"ok\": true") {
		t.Fatalf("unexpected JSON output: %s", out)
	}
}

func TestRunCreateRequiresInput(t *testing.T) {
	resetGlobals(t)
	err := runCreateCmd.RunE(runCreateCmd, nil)
	if err == nil || !strings.Contains(err.Error(), "provide at least one") {
		t.Fatalf("expected input validation error, got: %v", err)
	}
}

func TestRunCreateYAMLReadError(t *testing.T) {
	resetGlobals(t)
	createFile = filepath.Join(t.TempDir(), "missing.yaml")
	err := runCreateCmd.RunE(runCreateCmd, nil)
	if err == nil || !strings.Contains(err.Error(), "read file") {
		t.Fatalf("expected read file error, got: %v", err)
	}
}

func TestRunCreateYAMLParseError(t *testing.T) {
	resetGlobals(t)
	tmp := t.TempDir()
	file := filepath.Join(tmp, "bad.yaml")
	if writeErr := os.WriteFile(file, []byte("operations: [\n"), 0o644); writeErr != nil {
		t.Fatalf("write yaml: %v", writeErr)
	}
	createFile = file
	err := runCreateCmd.RunE(runCreateCmd, nil)
	if err == nil || !strings.Contains(err.Error(), "parse YAML") {
		t.Fatalf("expected yaml parse error, got: %v", err)
	}
}

func TestRunCreateWithScenarioID(t *testing.T) {
	resetGlobals(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/run/new" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"r1","state":"pending","created_at":"2026-01-01T00:00:00Z","scenario_id":42,"message":"ok"}`))
	}))
	defer server.Close()

	client = api.NewClient(server.URL, "k")
	createScenarioID = 42

	out, _, err := captureOutput(t, func() error {
		return runCreateCmd.RunE(runCreateCmd, nil)
	})
	if err != nil {
		t.Fatalf("run create failed: %v", err)
	}
	if !strings.Contains(out, "\"id\": \"r1\"") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRunCreateWithWait(t *testing.T) {
	resetGlobals(t)
	var getCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/run/new":
			_, _ = w.Write([]byte(`{"id":"r2","state":"pending","created_at":"2026-01-01T00:00:00Z","message":"ok"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/runs/r2":
			getCount++
			state := "running"
			if getCount >= 2 {
				state = "completed"
			}
			_, _ = w.Write([]byte(fmt.Sprintf(`{"data":{"id":"r2","created_at":"2026-01-01T00:00:00Z","state":"%s"}}`, state)))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client = api.NewClient(server.URL, "k")
	createScenarioID = 1
	createWait = true
	createPollInterval = 0

	out, errOut, err := captureOutput(t, func() error {
		return runCreateCmd.RunE(runCreateCmd, nil)
	})
	if err != nil {
		t.Fatalf("run create --wait failed: %v", err)
	}
	if !strings.Contains(errOut, "Waiting for completion") {
		t.Fatalf("expected progress on stderr, got: %s", errOut)
	}
	if !strings.Contains(out, "\"state\": \"completed\"") {
		t.Fatalf("expected completed run json, got: %s", out)
	}
}

func TestRunCreatePollError(t *testing.T) {
	resetGlobals(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/run/new":
			_, _ = w.Write([]byte(`{"id":"r3","state":"pending","created_at":"2026-01-01T00:00:00Z","message":"ok"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/runs/r3":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"fail"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client = api.NewClient(server.URL, "k")
	createScenarioID = 1
	createWait = true
	createPollInterval = 0

	err := runCreateCmd.RunE(runCreateCmd, nil)
	if err == nil || !strings.Contains(err.Error(), "poll run status") {
		t.Fatalf("expected poll error, got: %v", err)
	}
}

func TestRunGetAndDelete(t *testing.T) {
	resetGlobals(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/runs/r1":
			_, _ = w.Write([]byte(`{"data":{"id":"r1","created_at":"2026-01-01T00:00:00Z","state":"completed"}}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/runs/r1":
			_, _ = w.Write([]byte(`{"message":"deleted"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client = api.NewClient(server.URL, "k")

	getOut, _, getErr := captureOutput(t, func() error {
		return runGetCmd.RunE(runGetCmd, []string{"r1"})
	})
	if getErr != nil {
		t.Fatalf("run get failed: %v", getErr)
	}
	if !strings.Contains(getOut, "\"id\": \"r1\"") {
		t.Fatalf("unexpected get output: %s", getOut)
	}

	delOut, _, delErr := captureOutput(t, func() error {
		return runDeleteCmd.RunE(runDeleteCmd, []string{"r1"})
	})
	if delErr != nil {
		t.Fatalf("run delete failed: %v", delErr)
	}
	if !strings.Contains(delOut, "deleted") {
		t.Fatalf("unexpected delete output: %s", delOut)
	}
}

func TestRunGetAndDeleteErrorWrap(t *testing.T) {
	resetGlobals(t)
	client = api.NewClient("http://127.0.0.1:1", "k")

	if err := runGetCmd.RunE(runGetCmd, []string{"x"}); err == nil || !strings.Contains(err.Error(), "get run") {
		t.Fatalf("expected wrapped get error, got: %v", err)
	}
	if err := runDeleteCmd.RunE(runDeleteCmd, []string{"x"}); err == nil || !strings.Contains(err.Error(), "delete run") {
		t.Fatalf("expected wrapped delete error, got: %v", err)
	}
}

func TestRootPreRunRequiresAPIKey(t *testing.T) {
	resetGlobals(t)
	old := os.Getenv("SCALE_TEST_API_KEY")
	_ = os.Unsetenv("SCALE_TEST_API_KEY")
	defer func() { _ = os.Setenv("SCALE_TEST_API_KEY", old) }()

	err := rootCmd.PersistentPreRunE(rootCmd, nil)
	if err == nil || !strings.Contains(err.Error(), "API key required") {
		t.Fatalf("expected missing api key error, got: %v", err)
	}
}

func TestRootPreRunReadsEnv(t *testing.T) {
	resetGlobals(t)
	oldKey := os.Getenv("SCALE_TEST_API_KEY")
	oldBase := os.Getenv("SCALE_TEST_BASE_URL")
	defer func() {
		_ = os.Setenv("SCALE_TEST_API_KEY", oldKey)
		_ = os.Setenv("SCALE_TEST_BASE_URL", oldBase)
	}()

	_ = os.Setenv("SCALE_TEST_API_KEY", "env-key")
	_ = os.Setenv("SCALE_TEST_BASE_URL", "http://example.test/api/v1")

	err := rootCmd.PersistentPreRunE(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected pre-run error: %v", err)
	}
	if client == nil {
		t.Fatal("expected initialized client")
	}
	if apiKey != "env-key" {
		t.Fatalf("expected apiKey from env, got: %s", apiKey)
	}
	if baseURL != "http://example.test/api/v1" {
		t.Fatalf("expected baseURL from env, got: %s", baseURL)
	}
}

func TestRootPreRunDefaultBaseURL(t *testing.T) {
	resetGlobals(t)
	apiKey = "flag-key"

	oldBase := os.Getenv("SCALE_TEST_BASE_URL")
	_ = os.Unsetenv("SCALE_TEST_BASE_URL")
	defer func() { _ = os.Setenv("SCALE_TEST_BASE_URL", oldBase) }()

	err := rootCmd.PersistentPreRunE(rootCmd, nil)
	if err != nil {
		t.Fatalf("unexpected pre-run error: %v", err)
	}
	if baseURL != "https://scale-test.com/api/v1" {
		t.Fatalf("expected default base url, got: %s", baseURL)
	}
}

func TestRunCommandMetadata(t *testing.T) {
	if runCmd.Use == "" || runCmd.Short == "" {
		t.Fatalf("run command metadata should be set")
	}
}

func TestRootPersistentFlagsExist(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})
	_ = rootCmd.Execute()
	help := buf.String()
	if !strings.Contains(help, "--api-key") || !strings.Contains(help, "--base-url") {
		t.Fatalf("expected global flags in help output: %s", help)
	}
}
