package changes

import (
	"encoding/json"
	"testing"
)

func TestChangeRecordRoundTrips(t *testing.T) {
	r := ChangeRecord{ID: "endpoint-added", Kind: KindOperationAdded, Severity: Info,
		Method: "GET", Path: "/ping", Text: "operation added: GET /ping"}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var back ChangeRecord
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatal(err)
	}
	if back.ID != r.ID || back.Kind != r.Kind || back.Path != r.Path {
		t.Errorf("round-trip mismatch: %+v", back)
	}
}
