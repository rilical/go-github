package detector

import (
	"strings"

	"github.com/github/openapi-change-feed/internal/changes"
)

type DetectMeta struct{ BaseSpec, HeadSpec string }

func normalize(rd rawDiff, meta DetectMeta) ([]changes.ChangeRecord, changes.Summary) {
	var recs []changes.ChangeRecord
	add := func(r changes.ChangeRecord) {
		r.BaseSpec, r.HeadSpec = meta.BaseSpec, meta.HeadSpec
		recs = append(recs, r)
	}
	for _, op := range rd.AddedOps {
		add(changes.ChangeRecord{ID: "endpoint-added", Kind: changes.KindOperationAdded,
			Severity: changes.Info, Method: op.Method, Path: op.Path, Section: "paths",
			Text: "operation added: " + op.Method + " " + op.Path})
	}
	removalSeverity := map[opRef]changes.Severity{}
	for _, c := range rd.Changes {
		switch {
		case isRemovalRule(c.ID) && hasDeprecation(c.ID):
			removalSeverity[opRef{c.Method, c.Path}] = changes.Warn
		case isRemovalRule(c.ID):
			removalSeverity[opRef{c.Method, c.Path}] = changes.Breaking
		}
	}
	for _, op := range rd.DeletedOps {
		sev, ok := removalSeverity[op]
		if !ok {
			sev = changes.Breaking
		}
		add(changes.ChangeRecord{ID: "operation-removed", Kind: changes.KindOperationRemoved,
			Severity: sev, Method: op.Method, Path: op.Path, Section: "paths",
			Text: "operation removed: " + op.Method + " " + op.Path})
	}
	return recs, summarize(recs)
}

func isRemovalRule(id string) bool {
	return strings.HasPrefix(id, "api-path-removed") || strings.HasPrefix(id, "api-removed")
}

func hasDeprecation(id string) bool { return strings.HasSuffix(id, "-with-deprecation") }

func summarize(recs []changes.ChangeRecord) changes.Summary {
	s := changes.Summary{Total: len(recs)}
	for _, r := range recs {
		switch r.Severity {
		case changes.Breaking:
			s.Breaking++
		case changes.Warn:
			s.Warn++
		default:
			s.Info++
		}
		switch r.Kind {
		case changes.KindOperationAdded:
			s.Added++
		case changes.KindOperationRemoved:
			s.Removed++
		default:
			s.Modified++
		}
	}
	return s
}
