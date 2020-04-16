package initialize

import (
	"fmt"
	"os"

	"github.com/spin14/copy-basta/cmd/copy-basta/common/log"

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
	log.TheLogger.DebugWithData("user input", log.LoggerData{
		flagName: cmd.name,
	})
	log.TheLogger.Info("validating user input")
	if err := cmd.validate(); err != nil {
		return err
	}

	log.TheLogger.InfoWithData("bootstrapping new template project", log.LoggerData{"location": cmd.name})
	err := bootstrap.Bootstrap(cmd.name)
	if err != nil {
		return err
	}

	log.TheLogger.Info("done")
	return nil
}

func (cmd *Command) validate() error {
	if cmd.name == "" {
		return common.NewFlagValidationError(flagName, "is required")
	}

	if _, err := os.Stat(cmd.name); err == nil {
		return common.NewFlagValidationError(flagName, fmt.Sprintf("(%s) directory already exists", cmd.name))
	}
	return nil
}
