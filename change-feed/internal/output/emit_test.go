package output

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/github/openapi-change-feed/internal/changes"
)

func sampleBatch() changes.ChangeBatch {
	return changes.ChangeBatch{
		GeneratedAt: time.Unix(0, 0).UTC(),
		BaseSpec:    changes.SpecRef{ETag: `"v1"`},
		HeadSpec:    changes.SpecRef{ETag: `"v2"`},
		Changes: []changes.ChangeRecord{{ID: "endpoint-added", Kind: changes.KindOperationAdded,
			Severity: changes.Info, Method: "GET", Path: "/ping", Text: "operation added: GET /ping"}},
		Summary: changes.Summary{Total: 1, Info: 1, Added: 1},
	}
}

func TestEmitWritesThreeArtifacts(t *testing.T) {
	dir := t.TempDir()
	if err := NewFileEmitter(dir).Emit(context.Background(), sampleBatch()); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"changes.json", "feed.jsonl", "report.md"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("missing %s: %v", f, err)
		}
	}
	data, _ := os.ReadFile(filepath.Join(dir, "changes.json"))
	if err := ValidateFeed(data); err != nil {
		t.Errorf("changes.json failed validation: %v", err)
	}
	feed, _ := os.ReadFile(filepath.Join(dir, "feed.jsonl"))
	for i, line := range strings.Split(strings.TrimSpace(string(feed)), "\n") {
		if err := ValidateFeedLine([]byte(line)); err != nil {
			t.Errorf("feed.jsonl line %d invalid: %v", i, err)
		}
	}
}

func TestEmitIsIdempotentPerWindow(t *testing.T) {
	dir := t.TempDir()
	e := NewFileEmitter(dir)
	if err := e.Emit(context.Background(), sampleBatch()); err != nil {
		t.Fatal(err)
	}
	if err := e.Emit(context.Background(), sampleBatch()); err != nil {
		t.Fatal(err)
	}
	feed, _ := os.ReadFile(filepath.Join(dir, "feed.jsonl"))
	lines := strings.Split(strings.TrimSpace(string(feed)), "\n")
	if len(lines) != 1 {
		t.Fatalf("feed.jsonl has %d lines after a retried emit, want 1 (no duplicates)", len(lines))
	}
}
