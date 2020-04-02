package generate

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
)


type Command struct {
	projectName string
}


func New() *cobra.Command {
	cmd := Command{}

	cobraCmd := &cobra.Command{
		Use: "generate",
		Short: "generates new project based on the template and provided parameters",
		RunE: cmd.run,
	}

	cobraCmd.Flags().StringVar(&cmd.projectName,"project-name", "", `Project Name`)

	return cobraCmd
}

func (cmd *Command) validate() error {
	if cmd.projectName == "" {
		return errors.New(`[ERROR] "project-name" is required`)
	}
	return nil
}

func (cmd *Command) run(cobraCmd *cobra.Command, args []string) error {
	if err := cmd.validate(); err != nil {
		return err
	}

	if err := generateFromTemplate(cmd.projectName); err != nil {
		return err
	}
	
	log.Println("[INFO] Done!")
	return nil
}
