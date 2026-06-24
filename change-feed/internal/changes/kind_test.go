package changes

import "testing"

func TestKindForRuleID_RealOasdiffIDs(t *testing.T) {
	// Every id here was OBSERVED from the real POC window (see docs/ground-truth.md).
	cases := map[string]ChangeKind{
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
		"endpoint-deprecated":                    KindDeprecation,
		"totally-unknown-future-rule":            KindOther,
	}
	for id, want := range cases {
		if got := KindForRuleID(id); got != want {
			t.Errorf("KindForRuleID(%q) = %q, want %q", id, got, want)
		}
	}
}
