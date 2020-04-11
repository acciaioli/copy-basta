package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const (
	flagSrc  = "src"
	flagDest = "dest"
	flagSpec = "spec"

	defaultSpec = "spec.yaml"

	usageSrc  = "Generated Project root directory"
	usageDest = "Specification YAML file, relative to the template root directory"
	usageSpec = "Path to the YAML file with the variables to use in the templates"
)

type Flag struct {
	Ref     *string
	Name    string
	Default *string
	Usage   string
}

type Command struct {
	src                   string
	dest                  string
	spec                  string
	bastaTemplateVarsYAML string
}

/*

	cobraCmd.Flags().StringVar(
		&cmd.src,
		"src",
		"",
		`Template root directory`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.dest,
		"dest",
		"",
		`Project name (root directory for the generation output)`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.spec,
		"spec",
		"spec.yaml",
		`Specification YAML file, relative to the template root directory`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.bastaTemplateVarsYAML,
		"basta-yaml",
		"",
		`Path to the YAML file with the variables to use in the templates`,
	)
*/

func (cmd *Command) Flags() []Flag {
	return []Flag{
		{
			Ref:     &cmd.src,
			Name:    flagSrc,
			Default: nil,
			Usage:   usageSrc,
		},
		{
			Ref:     &cmd.dest,
			Name:    flagDest,
			Default: nil,
			Usage:   usageDest,
		},
		{
			Ref:     &cmd.spec,
			Name:    flagSpec,
			Default: sToP(defaultSpec),
			Usage:   usageSpec,
		},
		{
			Ref:     &cmd.bastaTemplateVarsYAML,
			Name:    "basta-yaml",
			Default: nil,
			Usage:   "Path to the YAML file with the variables to use in the templates",
		},
	}
}

func (cmd *Command) Run() error {
	log.Println("[INFO] Generating new project!")
	if err := cmd.validate(); err != nil {
		return err
	}

	files, err := parse(cmd.src)
	if err != nil {
		return err
	}

	templateVars, err := cmd.loadYAML(cmd.bastaTemplateVarsYAML)
	if err != nil {
		return err
	}

	err = write(cmd.dest, files, templateVars)
	if err != nil {
		return err
	}

	log.Println("[INFO] Done!")
	return nil
}

func (cmd *Command) validate() error {
	if cmd.src == "" {
		return errors.New(`[ERROR] "src" is required`)
	}

	if cmd.dest == "" {
		return errors.New(`[ERROR] "dest" is required`)
	}

	if cmd.spec == "" {
		return errors.New(`[ERROR] "spec" is required`)
	}
	spec := path.Join(cmd.src, cmd.spec)
	if err := fileExists(spec, "spec"); err != nil {
		return err
	}

	if cmd.bastaTemplateVarsYAML != "" {
		if err := fileExists(cmd.bastaTemplateVarsYAML, "basta-yaml"); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Command) loadYAML(filepath string) (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var templateVars = map[string]interface{}{}
	err = yaml.Unmarshal(yamlFile, &templateVars)
	if err != nil {
		return nil, err
	}

	return templateVars, nil
}

func fileExists(filePath string, name string) error {
	fInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf(`[ERROR] "%s" (%s) not found`, name, filePath)
	}
	if fInfo.IsDir() {
		return fmt.Errorf(`[ERROR] "%s" (%s) is not a file`, name, filePath)
	}
	return nil
}

func sToP(s string) *string {
	return &s
}
