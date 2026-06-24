package detector

import (
	"testing"

	"github.com/github/openapi-change-feed/internal/changes"
)

func TestNormalizeSyntheticRecordsHaveUniqueFingerprints(t *testing.T) {
	rd := rawDiff{
		AddedOps:   []opRef{{"GET", "/a"}, {"POST", "/b"}},
		DeletedOps: []opRef{{"DELETE", "/c"}},
	}
	recs, _ := normalize(rd, DetectMeta{})
	seen := map[string]string{}
	for _, r := range recs {
		if r.Fingerprint == "" {
			t.Errorf("record %s %s (%s) has empty fingerprint", r.Method, r.Path, r.ID)
			continue
		}
		key := r.Method + " " + r.Path
		if prev, ok := seen[r.Fingerprint]; ok && prev != key {
			t.Errorf("fingerprint %s shared by %q and %q", r.Fingerprint, prev, key)
		}
		seen[r.Fingerprint] = key
	}
	if len(seen) != 3 {
		t.Fatalf("want 3 distinct fingerprints, got %d", len(seen))
	}
}

func TestSyntheticFingerprintIsDeterministicAndUnique(t *testing.T) {
	a := syntheticFingerprint(changes.KindOperationAdded, "GET", "/a")
	again := syntheticFingerprint(changes.KindOperationAdded, "GET", "/a")
	byMethod := syntheticFingerprint(changes.KindOperationAdded, "POST", "/a")
	byKind := syntheticFingerprint(changes.KindOperationRemoved, "GET", "/a")
	if a == "" || a != again {
		t.Fatalf("fingerprint not deterministic: %q vs %q", a, again)
	}
	if a == byMethod || a == byKind {
		t.Fatalf("distinct ops collide: %q / %q / %q", a, byMethod, byKind)
	}
}

func TestNormalizeGenericRecordFallsBackToFingerprint(t *testing.T) {
	rd := rawDiff{Changes: []rawChange{
		{ID: "response-optional-property-added", Method: "GET", Path: "/x", Level: 1, Text: "t"},
	}}
	recs, _ := normalize(rd, DetectMeta{})
	if len(recs) != 1 {
		t.Fatalf("recs = %d, want 1", len(recs))
	}
	if recs[0].Fingerprint == "" {
		t.Error("generic record with no oasdiff fingerprint should fall back to a synthetic one")
	}
}
