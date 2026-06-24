package detector

import (
	"os"
	"testing"
)

func TestLoadRealSpecsSucceeds(t *testing.T) {
	for _, f := range []string{"testdata/base_min.json", "testdata/head_add_op.json"} {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		if _, err := loadSpec(b, f); err != nil {
			t.Fatalf("loadSpec(%s): %v", f, err)
		}
	}
}

func TestRunOasdiffDetectsAddedOperation(t *testing.T) {
	base, _ := os.ReadFile("testdata/base_min.json")
	head, _ := os.ReadFile("testdata/head_add_op.json")
	rd, err := runOasdiff(base, head, DefaultExclude)
	if err != nil {
		t.Fatal(err)
	}
	if len(rd.AddedOps) != 1 || rd.AddedOps[0].Method != "GET" || rd.AddedOps[0].Path != "/ping" {
		t.Fatalf("AddedOps = %+v, want [GET /ping]", rd.AddedOps)
	}
}
