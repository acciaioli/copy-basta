package initialize

import (
	"fmt"
	"log"
	"os"

	"github.com/spin14/copy-basta/cmd/copy-basta/commands/initialize/bootstrap"

	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

const (
	commandID          = "init"
	commandDescription = "bootstraps a new copy-basta template project"

	flagName      = "name"
	flagUsageName = "New Project root directory"
)

type Command struct {
	name string
}

func NewCommand() *Command {
	return &Command{}
}

func (cmd *Command) Name() string {
	return commandID
}

func (cmd *Command) Description() string {
	return commandDescription
}

func (cmd *Command) Flags() []common.CommandFlag {
	return []common.CommandFlag{
		{
			Ref:     &cmd.name,
			Name:    flagName,
			Default: nil,
			Usage:   flagUsageName,
		},
	}
}

func (cmd *Command) Run() error {
	log.Println("[INFO] Generating new template project!")
	if err := cmd.validate(); err != nil {
		return err
	}

	err := bootstrap.Bootstrap(cmd.name)
	if err != nil {
		return err
	}

	log.Println("[INFO] Done!")
	return nil
}

func (cmd *Command) validate() error {
	if cmd.name == "" {
		return fmt.Errorf(`[ERROR] "%s" is required`, flagName)
	}
	if _, err := os.Stat(cmd.name); err == nil {
		return fmt.Errorf(`[ERROR] "%s" (%s) already exists`, flagName, cmd.name)
	}
	return nil
}
