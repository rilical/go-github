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

func TestNormalizeRemovedOpSeverityFromDeprecation(t *testing.T) {
	rd := rawDiff{
		DeletedOps: []opRef{{"DELETE", "/old"}, {"DELETE", "/gone"}},
		Changes: []rawChange{
			{ID: "api-path-removed-with-deprecation", Method: "DELETE", Path: "/old", Level: 2},
			{ID: "api-path-removed-without-deprecation", Method: "DELETE", Path: "/gone", Level: 3},
		},
	}
	recs, sum := normalize(rd, DetectMeta{})
	got := map[string]changes.Severity{}
	for _, r := range recs {
		if r.Kind == changes.KindOperationRemoved {
			got[r.Path] = r.Severity
		}
	}
	if got["/old"] != changes.Warn {
		t.Errorf("/old = %v, want warn (deprecated removal)", got["/old"])
	}
	if got["/gone"] != changes.Breaking {
		t.Errorf("/gone = %v, want breaking", got["/gone"])
	}
	if sum.Removed != 2 {
		t.Errorf("Removed = %d, want 2", sum.Removed)
	}
}
