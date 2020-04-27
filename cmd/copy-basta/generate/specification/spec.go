package specification

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"copy-basta/cmd/copy-basta/common"
	"copy-basta/cmd/copy-basta/common/log"
)

type Spec struct {
	Variables []SpecVariable `yaml:"variables"`
}

func (spec *Spec) validate() error {
	for _, v := range spec.Variables {
		if err := v.validate(); err != nil {
			return fmt.Errorf("variables error: %s", err.Error())
		}
	}

	return nil
}

type Loader interface {
	LoadReader() (io.Reader, error)
}

func New(loader Loader) (*Spec, error) {
	r, err := loader.LoadReader()
	if err != nil {
		return nil, err
	}
	return newFromReader(r)
}

func newFromReader(r io.Reader) (*Spec, error) {
	spec := Spec{}
	if err := yaml.NewDecoder(r).Decode(&spec); err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
		return nil, errors.New("specification yaml file error: failed to decode yaml")
	}

	if err := spec.validate(); err != nil {
		return nil, fmt.Errorf("specification yaml file error: %s", err.Error())
	}
	return &spec, nil
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
		}
		if err := v.valueOk(value); err != nil {
			return nil, err
		}
	}

	return input, nil
}

func (spec *Spec) InputFromStdIn() (common.InputVariables, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("\n")
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

func promptLoop(r *bufio.Reader, v SpecVariable) (*string, error) {
	for {
		fmt.Print(v.prompt())
		userInput, err := r.ReadString('\n')
		fmt.Print("\n")
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
