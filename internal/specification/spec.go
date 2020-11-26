package specification

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"

	"copy-basta/internal/common/log"
	"copy-basta/internal/crawl"
)

type Spec struct {
	Ignorer   *Ignorer
	Passer    *Passer
	Variables Variables
}

func New(specFileName string, files []crawl.File, overwrite bool) (*Spec, error) {
	var specFile *crawl.File
	for _, f := range files {
		if f.Path == specFileName {
			specFile = &f
			break
		}
	}
	if specFile == nil {
		return nil, fmt.Errorf("specification: failed to find spec file (%s)", specFileName)
	}

	return newFromReader(specFile.Reader, overwrite)
}

func newFromReader(r io.Reader, overwrite bool) (*Spec, error) {
	data := SpecData{}
	err := yaml.NewDecoder(r).Decode(&data)
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
		return nil, errors.New("specification yaml file error: failed to decode yaml")
	}

	var ignoredPatterns []string
	ignoredPatterns = append(ignoredPatterns, data.Ignore...)
	if overwrite {
		ignoredPatterns = append(ignoredPatterns, data.OnOverwrite.Exclude...)
	}

	ignorer, err := NewIgnorer(ignoredPatterns)
	if err != nil {
		return nil, fmt.Errorf("ignorer error: %s", err.Error())
	}

	passer, err := NewPasser(data.PassThrough)
	if err != nil {
		return nil, fmt.Errorf("passer error: %s", err.Error())
	}

	variables, err := NewVariables(data.Variables)
	if err != nil {
		return nil, fmt.Errorf("variables error: %s", err.Error())
	}

	return &Spec{
		Ignorer:   ignorer,
		Passer:    passer,
		Variables: variables,
	}, nil
}
