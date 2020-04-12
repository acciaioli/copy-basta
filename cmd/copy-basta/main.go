package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate"
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

	cmd.AddCommand(generate.New())

	return cmd.Execute()
}
