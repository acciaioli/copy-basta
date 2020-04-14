package specification

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

type Spec struct {
	Variables []SpecVariable `yaml:"variables"`
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

func (spec *Spec) InputFromStdIn() (common.InputVariables, error) {
	r := bufio.NewReader(os.Stdin)

	inputVars := common.InputVariables{}

	for _, v := range spec.Variables {
		userInput, err := promptLoop(r, v)
		if err != nil {
			return nil, err
		}

		if userInput != nil {
			value, err := v.process(*userInput)
			if err != nil {
				return nil, err
			}
			inputVars[v.Name] = value
		} else {
			inputVars[v.Name] = v.Default
		}
	}

	return inputVars, nil
}

func (spec *Spec) InputFromFile(inputYAML string) (common.InputVariables, error) {
	yamlFile, err := ioutil.ReadFile(inputYAML)
	if err != nil {
		return nil, err
	}

	input := common.InputVariables{}
	err = yaml.Unmarshal(yamlFile, &input)
	if err != nil {
		return nil, err
	}

	for _, v := range spec.Variables {
		value, ok := input[v.Name]
		if !ok {
			if v.Default != nil {
				return nil, fmt.Errorf("no value nor default for %s", v.Name)
			}
			input[v.Name] = v.Default
		} else {
			if err := v.valueOk(value); err != nil {
				return nil, err
			}
		}

	}

	return input, nil
}

func promptLoop(r *bufio.Reader, v SpecVariable) (*string, error) {
	for {
		fmt.Print(v.prompt())
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
