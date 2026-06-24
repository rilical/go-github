package monitor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/github/openapi-change-feed/internal/changes"
	"github.com/github/openapi-change-feed/internal/detector"
	"github.com/github/openapi-change-feed/internal/output"
	"github.com/github/openapi-change-feed/internal/specfetch"
)

const specV1 = `{"openapi":"3.0.3","info":{"title":"t","version":"1"},"paths":{}}`

type stubFetcher struct {
	etag      string
	body      string
	unchanged bool
}

func (s *stubFetcher) Fetch(_ context.Context, _, lastETag string) (specfetch.FetchResult, error) {
	if s.unchanged && lastETag == s.etag {
		return specfetch.FetchResult{ETag: s.etag, Unchanged: true}, nil
	}
	return specfetch.FetchResult{Data: []byte(s.body), ETag: s.etag}, nil
}

func newMonitor(dir string, f specfetch.Fetcher) Monitor {
	return Monitor{Fetcher: f, Detector: detector.New(), Emitter: output.NewFileEmitter(dir),
		CursorPath: filepath.Join(dir, "cursor.json"), SpecURL: "http://x/spec.json"}
}

func TestFirstRunEstablishesBaselineNoEmit(t *testing.T) {
	dir := t.TempDir()
	m := newMonitor(dir, &stubFetcher{etag: `"v1"`, body: specV1})
	_, ran, err := m.RunOnce(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ran {
		t.Error("first run should establish baseline, not diff")
	}
	if _, err := os.Stat(filepath.Join(dir, "cursor.json")); err != nil {
		t.Error("cursor must exist after first run")
	}
}

func TestRunOnceIdempotentWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	m := newMonitor(dir, &stubFetcher{etag: `"v1"`, body: specV1, unchanged: true})
	if _, _, err := m.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	_, ran, err := m.RunOnce(context.Background())
	if err != nil || ran {
		t.Fatalf("second run ran=%v err=%v, want ran=false", ran, err)
	}
}

type failEmitter struct{}

func (failEmitter) Emit(context.Context, changes.ChangeBatch) error { return errBoom }

var errBoom = errorString("emit failed")

type errorString string

func (e errorString) Error() string { return string(e) }

func TestEmitFailureDoesNotAdvanceCursor(t *testing.T) {
	dir := t.TempDir()
	base := newMonitor(dir, &stubFetcher{etag: `"v1"`, body: specV1})
	if _, _, err := base.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	cur, _, _ := loadCursor(filepath.Join(dir, "cursor.json"))
	m := newMonitor(dir, &stubFetcher{etag: `"v2"`,
		body: `{"openapi":"3.0.3","info":{"title":"t","version":"1"},"paths":{"/n":{"get":{"responses":{"200":{"description":"ok"}}}}}}`})
	m.Emitter = failEmitter{}
	if _, _, err := m.RunOnce(context.Background()); err == nil {
		t.Fatal("expected emit error")
	}
	after, _, _ := loadCursor(filepath.Join(dir, "cursor.json"))
	if after.ETag != cur.ETag {
		t.Errorf("cursor advanced despite emit failure: %q -> %q", cur.ETag, after.ETag)
	}
}
