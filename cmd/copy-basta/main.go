package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate"
	"github.com/spin14/copy-basta/cmd/copy-basta/commands/initialize"
	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

func main() {
	err := execute()
	if err != nil {
		log.Fatal(err)
	}
}

func execute() error {
	cmd := &cobra.Command{
		Use:   "copy-basta",
		Short: "copy-basta utility",
		Long: `Basta! Stop copying.

This CLI can be used to bootstrap go projects in seconds, and stop the copy paste madness`,
	}

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
