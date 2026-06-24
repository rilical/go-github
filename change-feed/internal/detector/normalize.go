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
	addedKey := opKeySet(rd.AddedOps)
	deletedKey := opKeySet(rd.DeletedOps)
	for _, c := range rd.Changes {
		key := c.Method + " " + c.Path
		if c.ID == "endpoint-added" && addedKey[key] {
			continue
		}
		if isRemovalRule(c.ID) && deletedKey[key] {
			continue
		}
		add(changes.ChangeRecord{ID: c.ID, Fingerprint: c.Fingerprint,
			Kind:     changes.KindForRuleID(c.ID),
			Severity: changes.SeverityFromOasdiffLevel(c.Level),
			Method:   c.Method, Path: c.Path, OperationID: c.OperationID, Section: c.Section,
			Text: c.Text})
	}
	return recs, summarize(recs)
}

func isRemovalRule(id string) bool {
	return strings.HasPrefix(id, "api-path-removed") || strings.HasPrefix(id, "api-removed")
}

func hasDeprecation(id string) bool { return strings.HasSuffix(id, "-with-deprecation") }

func opKeySet(ops []opRef) map[string]bool {
	m := make(map[string]bool, len(ops))
	for _, o := range ops {
		m[o.Method+" "+o.Path] = true
	}
	return m
}

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
