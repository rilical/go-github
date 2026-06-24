package detector

import (
	"os"
	"testing"
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
