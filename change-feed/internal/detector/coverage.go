package detector

import (
	"fmt"
	"strings"

	"github.com/github/openapi-change-feed/internal/changes"
)

// coverageGap compares per-category structural counts from the raw diff against what normalization
// emitted. A non-empty category with zero emitted records is a silent false-negative; surface it.
func coverageGap(rd rawDiff, sum changes.Summary, meta DetectMeta) *changes.ChangeRecord {
	// modifiedRaw is a zero/non-zero sentinel for "the diff structurally contains modified or
	// component-level changes", not a matched count: it sums modified endpoints (rd.ModifiedOps)
	// with every non-add/non-removal checker row. Mixed units are intentional and only ever
	// compared against 0 below.
	modifiedRaw := rd.ModifiedOps
	for _, c := range rd.Changes {
		if c.ID != "endpoint-added" && !isRemovalRule(c.ID) {
			modifiedRaw++
		}
	}

	var gaps []string
	if len(rd.AddedOps) > 0 && sum.Added == 0 {
		gaps = append(gaps, fmt.Sprintf("%d added ops unmapped", len(rd.AddedOps)))
	}
	if len(rd.DeletedOps) > 0 && sum.Removed == 0 {
		gaps = append(gaps, fmt.Sprintf("%d removed ops unmapped", len(rd.DeletedOps)))
	}
	if modifiedRaw > 0 && sum.Modified == 0 {
		gaps = append(gaps, fmt.Sprintf("%d modified changes unmapped", modifiedRaw))
	}
	if len(gaps) == 0 {
		return nil
	}
	return &changes.ChangeRecord{
		ID:       "coverage-gap",
		Kind:     changes.KindOther,
		Severity: changes.Warn,
		Text:     "coverage gap: " + strings.Join(gaps, "; "),
		BaseSpec: meta.BaseSpec,
		HeadSpec: meta.HeadSpec,
	}
}
