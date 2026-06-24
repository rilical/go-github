package output

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/github/openapi-change-feed/internal/changes"
)

type Emitter interface {
	Emit(ctx context.Context, batch changes.ChangeBatch) error
}

type fileEmitter struct {
	dir string
}

func NewFileEmitter(dir string) Emitter {
	return fileEmitter{dir: dir}
}

func (e fileEmitter) Emit(_ context.Context, batch changes.ChangeBatch) error {
	if err := os.MkdirAll(e.dir, 0o755); err != nil {
		return err
	}

	pretty, err := json.MarshalIndent(batch, "", "  ")
	if err != nil {
		return err
	}
	if err := writeFileAtomic(filepath.Join(e.dir, "changes.json"), pretty); err != nil {
		return err
	}
	if err := writeFileAtomic(filepath.Join(e.dir, "report.md"), []byte(renderReport(batch))); err != nil {
		return err
	}

	window := batch.BaseSpec.ETag + ".." + batch.HeadSpec.ETag
	markerPath := filepath.Join(e.dir, ".feed_window")
	shouldAppend := true
	if prev, err := os.ReadFile(markerPath); err == nil && strings.TrimSpace(string(prev)) == window && window != ".." {
		shouldAppend = false
	}

	if shouldAppend {
		f, err := os.OpenFile(filepath.Join(e.dir, "feed.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(f)
		for _, c := range batch.Changes {
			if err := enc.Encode(c); err != nil {
				_ = f.Close()
				return err
			}
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	return os.WriteFile(markerPath, []byte(window), 0o644)
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
