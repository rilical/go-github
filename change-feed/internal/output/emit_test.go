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

func TestWriteFileAtomicCleansTmpOnRenameFailure(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := writeFileAtomic(target, []byte("data")); err == nil {
		t.Fatal("expected error renaming temp file onto a directory, got nil")
	}
	if _, err := os.Stat(target + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temp file %s.tmp was left behind after rename failure", target)
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
