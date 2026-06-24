package changes

import "testing"

func TestKindForRuleID_RealOasdiffIDs(t *testing.T) {
	cases := map[string]ChangeKind{
		// IDs observed in the 16-row real POC window (docs/ground-truth.md):
		"endpoint-added":                         KindOperationAdded,
		"api-path-removed-without-deprecation":   KindOperationRemoved,
		"api-removed-without-deprecation":        KindOperationRemoved, // alt spelling in ruleset
		"api-schema-removed":                     KindSchemaChanged,
		"response-required-property-removed":     KindResponseChanged,
		"response-optional-property-added":       KindResponseChanged,
		"response-property-became-nullable":      KindResponseChanged,
		"new-optional-request-property":          KindRequestBodyChanged,
		"request-property-became-optional":       KindRequestBodyChanged,
		"request-property-list-of-types-widened": KindRequestBodyChanged,

		// Additional real oasdiff v1.20.0 IDs exercising heuristic families not present in the POC window:
		"endpoint-deprecated": KindDeprecation, // exercises Contains("deprecat")

		// Synthetic probe for the KindOther fallback (not a real oasdiff id):
		"totally-unknown-future-rule": KindOther,
	}
	for id, want := range cases {
		if got := KindForRuleID(id); got != want {
			t.Errorf("KindForRuleID(%q) = %q, want %q", id, got, want)
		}
	}
}
