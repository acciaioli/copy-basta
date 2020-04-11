package generate

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := Command{}

	cobraCmd := &cobra.Command{
		Use:   CommandFlag,
		Short: CommandFlagShort,
		RunE: func(*cobra.Command, []string) error {
			return cmd.Run()
		},
	}

	for _, flag := range cmd.Flags() {
		cobraCmd.Flags().StringVar(
			flag.Ref,
			flag.Name,
			func(p *string) string {
				if p == nil {
					return ""
				}
				return *p
			}(flag.Default),
			flag.Usage,
		)
	}

	return cobraCmd
}
