package changes

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSeverityMarshalsAsString(t *testing.T) {
	for _, c := range []struct {
		sev  Severity
		want string
	}{{Info, `"info"`}, {Warn, `"warn"`}, {Breaking, `"breaking"`}} {
		b, err := json.Marshal(c.sev)
		if err != nil {
			t.Fatalf("marshal %v: %v", c.sev, err)
		}
		if string(b) != c.want {
			t.Errorf("marshal %v = %s, want %s", c.sev, b, c.want)
		}
		var back Severity
		if err := json.Unmarshal(b, &back); err != nil {
			t.Fatalf("unmarshal %s: %v", b, err)
		}
		if back != c.sev {
			t.Errorf("round-trip %s = %v, want %v", b, back, c.sev)
		}
	}
}

func TestSeverityUnmarshalRejectsUnknown(t *testing.T) {
	var s Severity
	if err := json.Unmarshal([]byte(`"bogus"`), &s); err == nil {
		t.Fatal("expected error unmarshaling an unknown severity string")
	}
}

func TestSummaryUsesSnakeCaseKeys(t *testing.T) {
	b, err := json.Marshal(Summary{Total: 9, Breaking: 1, Warn: 2, Info: 6, Added: 3, Removed: 1, Modified: 5})
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, key := range []string{`"total"`, `"breaking"`, `"warn"`, `"info"`, `"added"`, `"removed"`, `"modified"`} {
		if !strings.Contains(got, key) {
			t.Errorf("summary JSON %s missing key %s", got, key)
		}
	}
	if strings.Contains(got, `"Total"`) || strings.Contains(got, `"Added"`) {
		t.Errorf("summary JSON %s still uses capitalized Go field names", got)
	}
}
