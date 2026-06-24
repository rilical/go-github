package detector

import (
	"testing"

	"github.com/github/openapi-change-feed/internal/changes"
)

func TestNormalizeAddedOps(t *testing.T) {
	rd := rawDiff{AddedOps: []opRef{{"GET", "/a"}, {"POST", "/b"}}}
	recs, sum := normalize(rd, DetectMeta{BaseSpec: "x", HeadSpec: "y"})
	if sum.Added != 2 {
		t.Fatalf("Added = %d, want 2", sum.Added)
	}
	for _, r := range recs {
		if r.Kind != changes.KindOperationAdded || r.Severity != changes.Info {
			t.Errorf("bad added record: %+v", r)
		}
	}
}
