package monitor

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/github/openapi-change-feed/internal/changes"
)

type cursorState struct {
	Ref changes.SpecRef `json:"ref"`
}

// baselinePath is ALWAYS a fixed file alongside the cursor, never a path read from the cursor JSON.
// This removes the path-traversal / arbitrary-file-read concern both review voices flagged.
func baselinePath(cursorPath string) string {
	return filepath.Join(filepath.Dir(cursorPath), "baseline.json")
}

func loadCursor(path string) (changes.SpecRef, []byte, error) {
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return changes.SpecRef{}, nil, nil
	}
	if err != nil {
		return changes.SpecRef{}, nil, err
	}
	var st cursorState
	if err := json.Unmarshal(b, &st); err != nil {
		return changes.SpecRef{}, nil, err
	}
	spec, err := os.ReadFile(baselinePath(path))
	if errors.Is(err, os.ErrNotExist) {
		return st.Ref, nil, nil
	}
	return st.Ref, spec, err
}

// saveCursor writes the baseline first, then flips the cursor LAST and atomically (temp+rename).
// Ordering note: baseline.json is written before cursor.json, so a crash between the two writes
// leaves cursor.json pointing at the OLD ref while baseline.json already holds the NEW spec. On the
// next run loadCursor returns the OLD ref with the NEW baseline bytes; if HEAD has not moved, the
// diff is NEW-vs-NEW (empty) and the emitter rewrites changes.json/report.md as a zero-change
// snapshot until the next real API change. That is a recoverable blip in the convenience snapshot
// only: feed.jsonl, the append-only source of truth, is never corrupted (an empty diff appends
// nothing and the window marker blocks duplicate appends). The order is deliberately baseline-first:
// flipping it (cursor before baseline) would instead risk re-appending an already-emitted window to
// feed.jsonl if HEAD advanced again before recovery, which corrupts the durable stream. Fully atomic
// two-file recovery (content-addressed baselines) is deferred; this trade keeps the feed correct.
func saveCursor(path string, ref changes.SpecRef, spec []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := writeFileAtomic(baselinePath(path), spec); err != nil {
		return err
	}
	b, err := json.MarshalIndent(cursorState{Ref: ref}, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(path, b)
}

func writeFileAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
