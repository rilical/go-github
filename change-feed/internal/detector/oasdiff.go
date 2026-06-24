package detector

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oasdiff/oasdiff/diff"
	"github.com/oasdiff/oasdiff/load"
)

var DefaultExclude = []string{"description", "examples", "title", "summary"}

type opRef struct{ Method, Path string }

type rawChange struct {
	ID, Text, Method, Path, OperationID, Section, Fingerprint string
	Level                                                     int
}

type rawDiff struct {
	Changes              []rawChange
	AddedOps, DeletedOps []opRef
	ModifiedOps          int
}

func loadSpec(data []byte, name string) (*load.SpecInfo, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = false
	si, err := load.NewSpecInfoFromData(loader, data, name)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", name, err)
	}
	return si, nil
}

func runOasdiff(base, head []byte, exclude []string) (rawDiff, error) {
	baseInfo, err := loadSpec(base, "base.json")
	if err != nil {
		return rawDiff{}, err
	}
	headInfo, err := loadSpec(head, "head.json")
	if err != nil {
		return rawDiff{}, err
	}
	cfg := diff.NewConfig(diff.WithExcludeElements(exclude))
	report, _, err := diff.GetWithOperationsSourcesMap(cfg, baseInfo, headInfo)
	if err != nil {
		return rawDiff{}, fmt.Errorf("diff: %w", err)
	}
	var rd rawDiff
	if report.EndpointsDiff != nil {
		for _, e := range report.EndpointsDiff.Added {
			rd.AddedOps = append(rd.AddedOps, opRef{e.Method, e.Path})
		}
		for _, e := range report.EndpointsDiff.Deleted {
			rd.DeletedOps = append(rd.DeletedOps, opRef{e.Method, e.Path})
		}
		rd.ModifiedOps = len(report.EndpointsDiff.Modified)
	}
	return rd, nil
}
