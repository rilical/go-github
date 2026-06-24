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
	basePath := "testdata/base_min.json"
	base, err := os.ReadFile(basePath)
	if err != nil {
		t.Fatalf("read %s: %v", basePath, err)
	}
	headPath := "testdata/head_add_op.json"
	head, err := os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("read %s: %v", headPath, err)
	}
	rd, err := runOasdiff(base, head, DefaultExclude)
	if err != nil {
		t.Fatal(err)
	}
	if len(rd.AddedOps) != 1 || rd.AddedOps[0].Method != "GET" || rd.AddedOps[0].Path != "/ping" {
		t.Fatalf("AddedOps = %+v, want [GET /ping]", rd.AddedOps)
	}
}

func TestRunOasdiffPopulatesCheckerRowFields(t *testing.T) {
	basePath := "testdata/base_min.json"
	base, err := os.ReadFile(basePath)
	if err != nil {
		t.Fatalf("read %s: %v", basePath, err)
	}
	headPath := "testdata/head_add_op.json"
	head, err := os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("read %s: %v", headPath, err)
	}
	rd, err := runOasdiff(base, head, DefaultExclude)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, c := range rd.Changes {
		if c.ID == "endpoint-added" && c.Method == "GET" && c.Path == "/ping" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected an endpoint-added checker row for GET /ping; got %+v", rd.Changes)
	}
}
