package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Command struct {
	templateRoot          string
	projectName           string
	bastaTemplateVarsYAML string
}

func New() *cobra.Command {
	cmd := Command{}

	cobraCmd := &cobra.Command{
		Use:   "generate",
		Short: "generates new project based on the template and provided parameters",
		RunE:  cmd.run,
	}

	cobraCmd.Flags().StringVar(
		&cmd.templateRoot,
		"template-root",
		"",
		`Template root directory`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.projectName,
		"project-name",
		"",
		`Project name (root directory for the generation output)`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.bastaTemplateVarsYAML,
		"basta-yaml",
		"",
		`Path to the YAML file with the variables to use in the templates`,
	)

	return cobraCmd
}

func (cmd *Command) run(cobraCmd *cobra.Command, args []string) error {
	log.Println("[INFO] Generating new project!")
	if err := cmd.validate(); err != nil {
		return err
	}

	files, err := parse(cmd.templateRoot)
	if err != nil {
		return err
	}

	templateVars, err := cmd.loadYAML(cmd.bastaTemplateVarsYAML)
	if err != nil {
		return err
	}

	err = write(cmd.projectName, files, templateVars)
	if err != nil {
		return err
	}

	log.Println("[INFO] Done!")
	return nil
}

func (cmd *Command) validate() error {
	if cmd.templateRoot == "" {
		return errors.New(`[ERROR] "template-root" is required`)
	}
	if cmd.projectName == "" {
		return errors.New(`[ERROR] "project-name" is required`)
	}
	if cmd.bastaTemplateVarsYAML == "" {
		return errors.New(`[ERROR] "basta-yaml" is required`)
	}
	if _, err := os.Stat(cmd.bastaTemplateVarsYAML); os.IsNotExist(err) {
		return fmt.Errorf(`[ERROR] "basta-yaml" (%s) not found`, cmd.bastaTemplateVarsYAML)
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
