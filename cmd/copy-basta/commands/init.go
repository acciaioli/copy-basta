package commands

import (
	"github.com/spf13/cobra"

	"copy-basta/services/bootstrap"
)

func Init(globals func() error) *cobra.Command {
	const (
		commandUse         = "init"
		commandDescription = "bootstraps a new copy-basta template codebase"

		flagName            = "name"
		flagDescriptionName = "name and root directory of the new template codebase"
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
