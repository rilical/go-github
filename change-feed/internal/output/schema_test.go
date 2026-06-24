package output

import (
	"encoding/json"
	"testing"
)

func TestFeedSchemaIsValidJSON(t *testing.T) {
	if len(FeedSchema) == 0 {
		t.Fatal("FeedSchema is empty; the embed did not resolve")
	}
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(FeedSchema, &doc); err != nil {
		t.Fatalf("FeedSchema is not valid JSON: %v", err)
	}
	if _, ok := doc["$schema"]; !ok {
		t.Error("FeedSchema missing $schema key")
	}
}
