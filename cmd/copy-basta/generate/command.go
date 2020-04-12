package generate

import (
	"fmt"

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
			getDefault(flag.Default),
			getUsage(flag.Usage, flag.Default),
		)
	}

	return cobraCmd
}

func getDefault(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func getUsage(description string, defaultP *string) string {
	if defaultP != nil {
		return fmt.Sprintf("%s (default is %s)", description, *defaultP)
	}
	return description
}
