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

func TestNormalizeDedupsEndpointAddRemoveRows(t *testing.T) {
	rd := rawDiff{
		AddedOps: []opRef{{"GET", "/ping"}},
		Changes:  []rawChange{{ID: "endpoint-added", Method: "GET", Path: "/ping", Level: 1}},
	}
	recs, sum := normalize(rd, DetectMeta{})
	added := 0
	for _, r := range recs {
		if r.Kind == changes.KindOperationAdded {
			added++
		}
	}
	if added != 1 || sum.Added != 1 {
		t.Fatalf("added records = %d (sum %d), want exactly 1 (no double-count)", added, sum.Added)
	}
}

func TestNormalizeKeepsRemovalRowNotInEndpointsDiff(t *testing.T) {
	// A removal-family checker row whose op is NOT in DeletedOps is distinct and must survive.
	rd := rawDiff{
		DeletedOps: []opRef{{"DELETE", "/in-set"}},
		Changes: []rawChange{
			{ID: "api-path-removed-without-deprecation", Method: "DELETE", Path: "/in-set", Level: 3},
			{ID: "api-removed-without-deprecation", Method: "DELETE", Path: "/not-in-set", Level: 3},
		},
	}
	recs, _ := normalize(rd, DetectMeta{})
	var sawNotInSet bool
	removedRecords := 0
	for _, r := range recs {
		if r.Kind == changes.KindOperationRemoved {
			removedRecords++
		}
		if r.Path == "/not-in-set" {
			sawNotInSet = true
		}
	}
	if !sawNotInSet {
		t.Error("removal row for an op absent from DeletedOps was wrongly dropped (undercount)")
	}
	// /in-set is represented once via DeletedOps; its checker row is the only one suppressed.
	if removedRecords < 1 {
		t.Errorf("removed records = %d, want >=1", removedRecords)
	}
}
