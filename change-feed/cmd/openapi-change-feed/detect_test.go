package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/github/openapi-change-feed/internal/detector"
	"github.com/github/openapi-change-feed/internal/monitor"
	"github.com/github/openapi-change-feed/internal/output"
	"github.com/github/openapi-change-feed/internal/specfetch"
)

func TestNewDetectCommandHasFlags(t *testing.T) {
	cmd := newDetectCmd()
	if cmd.Use != "detect" {
		t.Fatalf("Use = %q, want detect", cmd.Use)
	}
	for _, f := range []string{"spec-url", "out", "cursor"} {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("missing --%s flag", f)
		}
	}
}

func TestDetectEndToEndWritesFeed(t *testing.T) {
	const base = `{"openapi":"3.0.3","info":{"title":"t","version":"1"},` +
		`"paths":{"/health":{"get":{"operationId":"health","responses":{"200":{"description":"ok"}}}}}}`
	const head = `{"openapi":"3.0.3","info":{"title":"t","version":"1"},` +
		`"paths":{"/health":{"get":{"operationId":"health","responses":{"200":{"description":"ok"}}}},` +
		`"/ping":{"get":{"operationId":"ping","responses":{"200":{"description":"ok"}}}}}}`

	var n int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n++
		if n == 1 {
			w.Header().Set("ETag", `"v1"`)
			_, _ = io.WriteString(w, base)
			return
		}
		if r.Header.Get("If-None-Match") == `"v2"` {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("ETag", `"v2"`)
		_, _ = io.WriteString(w, head)
	}))
	defer srv.Close()

	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	m := monitor.Monitor{
		Fetcher: specfetch.New(srv.Client()), Detector: detector.New(),
		Emitter: output.NewFileEmitter(out), CursorPath: filepath.Join(dir, "state", "cursor.json"),
		SpecURL: srv.URL,
	}
	if _, _, err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ran, err := m.RunOnce(context.Background()); err != nil || !ran {
		t.Fatalf("second run ran=%v err=%v, want ran=true", ran, err)
	}
	got, err := os.ReadFile(filepath.Join(out, "changes.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "/ping") || !strings.Contains(string(got), "operation.added") {
		t.Fatalf("changes.json missing the added op:\n%s", got)
	}
}

func TestDetectRequiresSpecURL(t *testing.T) {
	cmd := newDetectCmd()
	cmd.SetArgs(nil)
	cmd.SilenceUsage, cmd.SilenceErrors = true, true
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "spec-url") {
		t.Fatalf("want a required-flag error naming spec-url, got %v", err)
	}
}
