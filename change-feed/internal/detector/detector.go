package detector

import "github.com/github/openapi-change-feed/internal/changes"

type DetectOptions struct {
	ExcludeElements   []string
	SeverityOverrides map[string]changes.Severity
}

type Detector interface {
	Detect(base, head []byte, meta DetectMeta, opts DetectOptions) ([]changes.ChangeRecord, changes.Summary, error)
}

type detector struct{}

func New() Detector { return detector{} }

func (detector) Detect(base, head []byte, meta DetectMeta, opts DetectOptions) ([]changes.ChangeRecord, changes.Summary, error) {
	exclude := opts.ExcludeElements
	if exclude == nil {
		exclude = DefaultExclude
	}
	rd, err := runOasdiff(base, head, exclude)
	if err != nil {
		return nil, changes.Summary{}, err
	}
	recs, sum := normalize(rd, meta)
	if gap := coverageGap(rd, sum, meta); gap != nil {
		recs = append(recs, *gap)
		sum = summarize(recs)
	}
	return recs, sum, nil
}
