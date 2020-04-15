package main

import (
	"github.com/spf13/cobra"

	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate"
	"github.com/spin14/copy-basta/cmd/copy-basta/commands/initialize"
	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

func main() {
	err := execute()
	if err != nil {
		panic(err)
	}
}

func execute() error {
	cmd := &cobra.Command{
		Use:   "copy-basta",
		Short: "copy-basta utility",
		Long: `Basta! Stop copying.

This CLI can be used to bootstrap go projects in seconds, and stop the copy paste madness`,
	}
	var logLevel string
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", `Used to set the logging level. 
Available options: [debug, info, warn, error, fatal]`)

	cmd.AddCommand(newCobraCommand(generate.NewCommand(), logLevel))
	cmd.AddCommand(newCobraCommand(initialize.NewCommand(), logLevel))

	return cmd.Execute()
}

type CommandInterface interface {
	Name() string
	Description() string
	Flags() []common.CommandFlag
	Run(*common.Logger) error
}

func newCobraCommand(cmd CommandInterface, logLevel string) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   cmd.Name(),
		Short: cmd.Description(),
		RunE: func(*cobra.Command, []string) error {
			logger, err := common.NewLogger(common.WithLevelS(logLevel))
			if err != nil {
				return err
			}
			return cmd.Run(logger)
		},
	}

	for _, flag := range cmd.Flags() {
		cobraCmd.Flags().StringVar(
			flag.Ref,
			flag.Name,
			getDefault(flag.Default),
			flag.Usage,
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
