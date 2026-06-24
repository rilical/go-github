package detector

import (
	"testing"

	"github.com/github/openapi-change-feed/internal/changes"
)

func TestCoverageGapFiresWhenModifiedOpsUnmapped(t *testing.T) {
	rd := rawDiff{ModifiedOps: 3}
	gap := coverageGap(rd, changes.Summary{}, DetectMeta{})
	if gap == nil || gap.ID != "coverage-gap" || gap.Severity != changes.Warn {
		t.Fatalf("expected a warn coverage-gap record, got %+v", gap)
	}
}

func TestCoverageGapFiresWhenRemovedOpsUnmapped(t *testing.T) {
	// oasdiff structurally saw a deleted endpoint, but no record mapped to the Removed category.
	rd := rawDiff{DeletedOps: []opRef{{"DELETE", "/x"}}}
	gap := coverageGap(rd, changes.Summary{}, DetectMeta{})
	if gap == nil || gap.ID != "coverage-gap" || gap.Severity != changes.Warn {
		t.Fatalf("expected a warn coverage-gap record for unmapped removed op, got %+v", gap)
	}
}

func TestCoverageGapSilentWhenEveryCategoryCovered(t *testing.T) {
	rd := rawDiff{AddedOps: []opRef{{"GET", "/a"}}, ModifiedOps: 1}
	if gap := coverageGap(rd, changes.Summary{Added: 1, Modified: 2}, DetectMeta{}); gap != nil {
		t.Fatalf("no gap expected when each category has records, got %+v", gap)
	}
}
