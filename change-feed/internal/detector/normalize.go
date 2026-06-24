package detector

import "github.com/github/openapi-change-feed/internal/changes"

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
	return recs, summarize(recs)
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
