package commands

import (
	"github.com/spf13/cobra"

	"copy-basta/services/bootstrap"
)

const (
	commandID          = "init"
	commandDescription = "bootstraps a new copy-basta template project"

	flagName      = "name"
	flagUsageName = "New Project root directory"
)

func InitCommand(globals func() error) *cobra.Command {
	const (
		commandUse         = "init"
		commandDescription = "bootstraps a new copy-basta template project"

		flagName            = "name"
		flagDescriptionName = "New Project root directory"
	)

	var name string

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, what []string) error {
			err := globals()
			if err != nil {
				return err
			}
			return bootstrap.Bootstrap(&bootstrap.Params{
				Name: name,
			})
		},
	}

	cmd.Flags().StringVar(
		&name,
		flagName,
		"",
		flagDescriptionName,
	)

	return cmd
}
