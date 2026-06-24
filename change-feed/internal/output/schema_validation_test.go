package output

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateFeedRejectsBadSeverity(t *testing.T) {
	bad := []byte(`{"generated_at":"2020-01-01T00:00:00Z","changes":[{"id":"x","kind":"operation.added","severity":"bogus"}],"summary":{"total":1,"breaking":0,"warn":0,"info":1,"added":1,"removed":0,"modified":0}}`)
	if err := ValidateFeed(bad); err == nil {
		t.Fatal("expected schema validation to reject severity \"bogus\"")
	}
}

func TestValidateFeedRejectsMissingSummary(t *testing.T) {
	bad := []byte(`{"generated_at":"2020-01-01T00:00:00Z","changes":[]}`)
	if err := ValidateFeed(bad); err == nil {
		t.Fatal("expected schema validation to reject a feed with no summary")
	}
}

func TestValidateFeedLineRejectsBadSeverity(t *testing.T) {
	if err := ValidateFeedLine([]byte(`{"id":"x","kind":"operation.added","severity":"bogus"}`)); err == nil {
		t.Fatal("expected feed-line validation to reject severity \"bogus\"")
	}
}

func TestValidateFeedAcceptsRealEmitterOutput(t *testing.T) {
	dir := t.TempDir()
	if err := NewFileEmitter(dir).Emit(context.Background(), sampleBatch()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "changes.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateFeed(data); err != nil {
		t.Fatalf("emitter output should pass schema validation: %v", err)
	}
}
