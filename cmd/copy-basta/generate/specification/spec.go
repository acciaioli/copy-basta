package specification

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate/common"

	"gopkg.in/yaml.v2"
)

type Spec struct {
	Variables map[string]SpecVariable `yaml:variables`
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

func (spec *Spec) PromptInput() (common.InputVariables, error) {
	r := bufio.NewReader(os.Stdin)

	inputVars := common.InputVariables{}

	for k, v := range spec.Variables {
		userInput, err := promptLoop(r, k, v)
		if err != nil {
			return nil, err
		}

		if userInput != nil {
			value, err := v.process(*userInput)
			if err != nil {
				return nil, err
			}
			inputVars[k] = value
		} else {
			inputVars[k] = v.Default
		}
	}

	return inputVars, nil
}

func promptLoop(r *bufio.Reader, k string, v SpecVariable) (*string, error) {
	for {
		fmt.Print(v.prompt(k))
		userInput, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		userInput = strings.TrimSuffix(userInput, "\n")

		if userInput != "" {
			return &userInput, nil
		}

		if v.Default != nil {
			return nil, nil
		}
	}
}
