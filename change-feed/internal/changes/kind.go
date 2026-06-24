package changes

import "strings"

type ChangeKind string

const (
	KindOperationAdded     ChangeKind = "operation.added"
	KindOperationRemoved   ChangeKind = "operation.removed"
	KindRequestParam       ChangeKind = "request.param.changed"
	KindRequestBodyChanged ChangeKind = "request.body.changed"
	KindResponseChanged    ChangeKind = "response.changed"
	KindSchemaChanged      ChangeKind = "schema.changed"
	KindSecurityChanged    ChangeKind = "security.changed"
	KindDeprecation        ChangeKind = "operation.deprecated"
	KindOther              ChangeKind = "other"
)

// KindForRuleID maps an oasdiff rule id to a coarse family by inspecting id structure.
// oasdiff ships ~992 rule ids; we classify by family rather than enumerate. Order matters:
// the first matching clause wins, most specific first. See docs/ground-truth.md.
func KindForRuleID(id string) ChangeKind {
	switch {
	case id == "endpoint-added":
		return KindOperationAdded
	case strings.HasPrefix(id, "api-path-removed"), strings.HasPrefix(id, "api-removed"),
		id == "api-operation-removed":
		return KindOperationRemoved
	case strings.Contains(id, "deprecat"):
		return KindDeprecation
	case strings.HasPrefix(id, "api-schema"):
		return KindSchemaChanged
	case strings.HasPrefix(id, "api-security"), strings.Contains(id, "security-requirement"):
		return KindSecurityChanged
	case strings.HasPrefix(id, "request-parameter"):
		return KindRequestParam
	case strings.HasPrefix(id, "request-"), strings.Contains(id, "request-property"):
		return KindRequestBodyChanged
	case strings.HasPrefix(id, "response-"):
		return KindResponseChanged
	default:
		return KindOther
	}
}
