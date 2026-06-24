package detector

import (
	"os"
	"strings"
	"testing"
)

func TestAnchorWindow(t *testing.T) {
	base, baseErr := os.ReadFile("../../testdata/anchor/base.json")
	head, headErr := os.ReadFile("../../testdata/anchor/head.json")
	if baseErr != nil || headErr != nil {
		if os.Getenv("ANCHOR_REQUIRED") != "" {
			t.Fatal("ANCHOR_REQUIRED set but anchor specs absent; expected ../../testdata/anchor/{base,head}.json")
		}
		t.Skip("anchor specs absent; expected ../../testdata/anchor/{base,head}.json")
	}

	recs, sum, err := New().Detect(base, head, DetectMeta{BaseSpec: "1d626756", HeadSpec: "3e08e45a"}, DetectOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("anchor summary total=%d breaking=%d warn=%d info=%d added=%d removed=%d", sum.Total, sum.Breaking, sum.Warn, sum.Info, sum.Added, sum.Removed)

	if sum.Total != 348 || sum.Breaking != 36 || sum.Warn != 70 || sum.Info != 242 {
		t.Fatalf("headline mismatch: total=%d breaking=%d warn=%d info=%d; want 348/36/70/242", sum.Total, sum.Breaking, sum.Warn, sum.Info)
	}
	if sum.Added != 6 || sum.Removed != 2 {
		t.Fatalf("added/removed mismatch: %d/%d; want 6/2", sum.Added, sum.Removed)
	}

	foundAccessSource := false
	for _, rec := range recs {
		if strings.Contains(rec.Text, "access_source") {
			foundAccessSource = true
			break
		}
	}
	if !foundAccessSource {
		t.Fatal("expected Team.access_source drift record")
	}
}
