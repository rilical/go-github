package detector

import (
	"os"
	"testing"

	"github.com/github/openapi-change-feed/internal/changes"
)

func TestDetectAddedOperationEndToEnd(t *testing.T) {
	base, err := os.ReadFile("testdata/base_min.json")
	if err != nil {
		t.Fatalf("read base: %v", err)
	}
	head, err := os.ReadFile("testdata/head_add_op.json")
	if err != nil {
		t.Fatalf("read head: %v", err)
	}
	recs, sum, err := New().Detect(base, head, DetectMeta{BaseSpec: "x", HeadSpec: "y"}, DetectOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if sum.Added != 1 || len(recs) == 0 {
		t.Fatalf("Added=%d recs=%d, want Added=1 and >=1 record", sum.Added, len(recs))
	}
}

func TestDetectAppliesSeverityOverrides(t *testing.T) {
	base, err := os.ReadFile("testdata/base_min.json")
	if err != nil {
		t.Fatalf("read base: %v", err)
	}
	head, err := os.ReadFile("testdata/head_add_op.json")
	if err != nil {
		t.Fatalf("read head: %v", err)
	}
	_, sum, err := New().Detect(base, head, DetectMeta{}, DetectOptions{
		SeverityOverrides: map[string]changes.Severity{
			"endpoint-added": changes.Warn,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if sum.Warn < 1 {
		t.Fatalf("override to warn not applied: %+v", sum)
	}
}
