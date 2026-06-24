package monitor

import (
	"context"
	"time"

	"github.com/github/openapi-change-feed/internal/changes"
	"github.com/github/openapi-change-feed/internal/detector"
	"github.com/github/openapi-change-feed/internal/output"
	"github.com/github/openapi-change-feed/internal/specfetch"
)

type Monitor struct {
	Fetcher    specfetch.Fetcher
	Detector   detector.Detector
	Emitter    output.Emitter
	CursorPath string
	SpecURL    string
}

// RunOnce fetches HEAD; if changed, diffs baseline->HEAD, emits, and advances the cursor.
// The cursor advances ONLY after a successful emit. Returns (summary, ranDiff, error).
// If the emit succeeds but the trailing saveCursor fails, RunOnce returns (summary, true, err):
// ran=true reflects that the feed was already written, and the window marker keeps the inevitable
// retry from duplicating it.
func (m Monitor) RunOnce(ctx context.Context) (changes.Summary, bool, error) {
	baseRef, baseSpec, err := loadCursor(m.CursorPath)
	if err != nil {
		return changes.Summary{}, false, err
	}
	res, err := m.Fetcher.Fetch(ctx, m.SpecURL, baseRef.ETag)
	if err != nil {
		return changes.Summary{}, false, err
	}
	if res.Unchanged {
		return changes.Summary{}, false, nil
	}
	if baseRef.ETag != "" && res.ETag == baseRef.ETag {
		// Origins without ETag support never return 304; specfetch hands back a
		// content-hash ETag instead. A matching hash means identical content, so
		// there is nothing to diff and no snapshot to rewrite.
		return changes.Summary{}, false, nil
	}
	headRef := changes.SpecRef{URL: m.SpecURL, ETag: res.ETag, FetchedAt: time.Now().UTC()}
	if baseSpec == nil {
		return changes.Summary{}, false, saveCursor(m.CursorPath, headRef, res.Data)
	}
	recs, sum, err := m.Detector.Detect(baseSpec, res.Data,
		detector.DetectMeta{BaseSpec: baseRef.ETag, HeadSpec: res.ETag}, detector.DetectOptions{})
	if err != nil {
		return changes.Summary{}, false, err
	}
	batch := changes.ChangeBatch{GeneratedAt: time.Now().UTC(), BaseSpec: baseRef,
		HeadSpec: headRef, Changes: recs, Summary: sum}
	if err := m.Emitter.Emit(ctx, batch); err != nil {
		return changes.Summary{}, false, err
	}
	return sum, true, saveCursor(m.CursorPath, headRef, res.Data)
}

// Run satisfies the locked Master section 7 contract `Monitor interface { Run(ctx) error }`.
// It runs one detection cycle and discards the richer (summary, ran) detail that RunOnce exposes for the CLI.
func (m Monitor) Run(ctx context.Context) error {
	_, _, err := m.RunOnce(ctx)
	return err
}
