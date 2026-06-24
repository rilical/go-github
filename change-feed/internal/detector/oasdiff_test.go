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
