package main

import (
	"fmt"

	"copy-basta/cmd/common/log"

	"github.com/spf13/cobra"

	"copy-basta/cmd/common"
	"copy-basta/cmd/generate"
	"copy-basta/cmd/initialize"
)

const (
	version = "snapshot" // build-time variable

	cmdUse   = "copy-basta"
	cmdShort = "copy-basta utility"
	cmdLong  = "Basta! Stop copying.\n\nThis CLI can be used to bootstrap go projects in seconds, and stop the copy paste madness"

	flagLogLevel            = "log-level"
	flagLogLevelDefault     = "info"
	flagLogLevelDescription = "Used to set the logging level.\nAvailable options: [debug, info, warn, error, fatal]"
)

func main() {
	if err := execute(); err != nil {
		log.L.Error(err.Error())
		fmt.Println("command failed.")
	}
}

var globals = struct {
	logLevel string
}{}

func execute() error {
	cmd := &cobra.Command{
		Use:     cmdUse,
		Short:   cmdShort,
		Long:    cmdLong,
		Version: version,
	}

	cmd.PersistentFlags().StringVar(&globals.logLevel, flagLogLevel, flagLogLevelDefault, flagLogLevelDescription)

	cmd.AddCommand(newCobraCommand(generate.NewCommand()))
	cmd.AddCommand(newCobraCommand(initialize.NewCommand()))

	return cmd.Execute()
}

type CommandInterface interface {
	Name() string
	Description() string
	Flags() []common.CommandFlag
	Run() error
}

func newCobraCommand(cmd CommandInterface) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   cmd.Name(),
		Short: cmd.Description(),
		RunE: func(*cobra.Command, []string) error {
			logLevel, err := log.StringToLevel(globals.logLevel)
			if err != nil {
				return common.NewFlagValidationError(flagLogLevel, err.Error())
			}
			log.L.SetLevel(logLevel)
			return cmd.Run()
		},
		SilenceErrors: true,
		SilenceUsage:  false,
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
