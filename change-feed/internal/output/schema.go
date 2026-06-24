package output

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// FeedSchema is the embedded JSON Schema for the change feed, exported so
// external validators and CLI consumers can read it without importing the file.
//
//go:embed feed.schema.json
var FeedSchema []byte

func ValidateFeed(b []byte) error {
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(b, &doc); err != nil {
		return fmt.Errorf("feed not valid JSON: %w", err)
	}
	for _, k := range []string{"generated_at", "changes", "summary"} {
		if _, ok := doc[k]; !ok {
			return fmt.Errorf("feed missing required key %q", k)
		}
	}
	return nil
}

func ValidateFeedLine(b []byte) error {
	var rec map[string]json.RawMessage
	if err := json.Unmarshal(b, &rec); err != nil {
		return fmt.Errorf("feed line not valid JSON: %w", err)
	}
	for _, k := range []string{"id", "kind", "severity"} {
		if _, ok := rec[k]; !ok {
			return fmt.Errorf("feed line missing required key %q", k)
		}
	}
	return nil
}
