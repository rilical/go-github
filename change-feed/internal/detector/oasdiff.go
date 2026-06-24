package detector

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oasdiff/oasdiff/load"
)

func loadSpec(data []byte, name string) (*load.SpecInfo, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = false
	si, err := load.NewSpecInfoFromData(loader, data, name)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", name, err)
	}
	return si, nil
}
