package main

import "testing"

func TestNewDetectCommandHasFlags(t *testing.T) {
	cmd := newDetectCmd()
	if cmd.Use != "detect" {
		t.Fatalf("Use = %q, want detect", cmd.Use)
	}
	for _, f := range []string{"spec-url", "out", "cursor"} {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("missing --%s flag", f)
		}
	}
}
