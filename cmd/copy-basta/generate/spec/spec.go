package spec

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/prometheus/common/log"

	"gopkg.in/yaml.v2"
)

/*
https://swagger.io/docs/specification/data-models/data-types/
*/

const (
	openAPiString  = "string"
	openAPiNumber  = "number"
	openAPiInteger = "integer"
	openAPiBoolean = "boolean"
	openAPiArray   = "array"
	openAPiObject  = "object"
)

type Variable struct {
	Type        string      `yaml:"type"`
	Default     interface{} `yaml:"default"`
	Description string      `yaml:"description"`
}

func (v *Variable) validate() error {
	log.Warn("[WARN] type checks are currently not supported")

	// type checks
	if v.Type == "" {
		return errors.New("variable validate error: type is required")
	}
	if ok := func(t string) bool {
		for _, t := range []string{openAPiString, openAPiNumber, openAPiInteger, openAPiBoolean, openAPiArray, openAPiObject} {
			if v.Type == t {
				return true
			}
		}
		return false
	}(v.Type); !ok {
		return fmt.Errorf("variable validate error: %s is not a valid type", v.Type)
	}

	// default checks
	log.Warn("[WARN] default values type checks are currently not supported")

	return nil
}

type Spec struct {
	Variables map[string]Variable `yaml:variables`
}

func (spec *Spec) validate() error {
	for _, v := range spec.Variables {
		if err := v.validate(); err != nil {
			return err
		}
	}

	return nil
}

func New(specsYAML string) (*Spec, error) {
	if f, err := os.Open(specsYAML); err != nil {
		return nil, err
	} else {
		return newFromReader(f)
	}
}

func newFromReader(r io.Reader) (*Spec, error) {
	spec := Spec{}
	if err := yaml.NewDecoder(r).Decode(&spec); err != nil {
		return nil, err
	}

	if err := spec.validate(); err != nil {
		return nil, err
	}
	return &spec, nil
}
