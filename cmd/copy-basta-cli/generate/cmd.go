package generate

import (
	"errors"
	"log"

	"github.com/spf13/cobra"
)

type Command struct {
	templateRoot string
	projectName  string
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

	return cobraCmd
}

func (cmd *Command) validate() error {
	if cmd.templateRoot == "" {
		return errors.New(`[ERROR] "template-root" is required`)
	}
	if cmd.projectName == "" {
		return errors.New(`[ERROR] "project-name" is required`)
	}
	return nil
}

func (cmd *Command) run(cobraCmd *cobra.Command, args []string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	files, err := parse(cmd.templateRoot)
	if err != nil {
		return err
	}

	err = write(cmd.projectName, files)
	if err != nil {
		return err
	}

	log.Println("[INFO] Done!")
	return nil
}
