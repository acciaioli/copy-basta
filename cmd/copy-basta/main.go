package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"copy-basta/cmd/copy-basta/commands"
	"copy-basta/internal/common/log"
)

var version = "snapshot" // build-time variable

const (
	cmdUse   = "copy-basta"
	cmdShort = "copy-basta utility"
	cmdLong  = `Basta! Stop copying.

This CLI can be used to bootstrap go projects in seconds, and stop the copy paste madness`
)

func main() {
	if err := execute(); err != nil {
		log.L.Error(err.Error())
		fmt.Println("command failed.")
	}
}

type globals struct {
	logLevel string
}

func (g *globals) register(cmd *cobra.Command) {
	const flag = "log-level"
	cmd.PersistentFlags().StringVar(
		&g.logLevel,
		flag,
		log.Error,
		fmt.Sprintf(
			"global logging level. one of [%s, %s, %s, %s, %s]",
			log.Debug,
			log.Info,
			log.Warn,
			log.Error,
			log.Fatal,
		),
	)
}

func (g *globals) process() error {
	lvl, err := log.ToLevel(g.logLevel)
	if err != nil {
		return err
	}
	log.L.SetLevel(lvl)
	return nil
}

func execute() error {
	cmd := &cobra.Command{
		Use:           cmdUse,
		Short:         cmdShort,
		Long:          cmdLong,
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	globals := globals{}
	globals.register(cmd)

	cmd.AddCommand(commands.Init(globals.process))
	cmd.AddCommand(commands.Generate(globals.process))

	return cmd.Execute()
}
