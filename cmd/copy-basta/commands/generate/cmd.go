package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spin14/copy-basta/cmd/copy-basta/common/log"

	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate/specification"

	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate/parse"
	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate/write"
	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

const (
	commandID          = "generate"
	commandDescription = "generates new project based on the template and provided parameters"

	flagSrc            = "src"
	flagDescriptionSrc = "Generated Project root directory"

	flagDest            = "dest"
	flagDescriptionDest = "Specification YAML file, relative to the template root directory"

	flagSpec            = "spec"
	flagDefaultSpec     = common.SpecFile
	flagDescriptionSpec = "Path to the YAML containing the template specification"

	flagInput            = "input"
	flagDescriptionInput = "Path to the YAML file with the variables to use in the templates"
)

type Command struct {
	src       string
	dest      string
	specYAML  string
	inputYAML string
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
			Ref:     &cmd.src,
			Name:    flagSrc,
			Default: nil,
			Usage:   flagDescriptionSrc,
		},
		{
			Ref:     &cmd.dest,
			Name:    flagDest,
			Default: nil,
			Usage:   flagDescriptionDest,
		},
		{
			Ref:     &cmd.specYAML,
			Name:    flagSpec,
			Default: sToP(flagDefaultSpec),
			Usage:   flagDescriptionSpec,
		},
		{
			Ref:     &cmd.inputYAML,
			Name:    flagInput,
			Default: nil,
			Usage:   flagDescriptionInput,
		},
	}
}

func (cmd *Command) Run() error {
	log.TheLogger.DebugWithData("user input", log.LoggerData{
		flagSrc:   cmd.src,
		flagDest:  cmd.dest,
		flagSpec:  cmd.specYAML,
		flagInput: cmd.inputYAML,
	})
	log.TheLogger.Info("validating user input")
	if err := cmd.validate(); err != nil {
		return err
	}

	log.TheLogger.Info("loading specification file")
	spec, err := specification.New(cmd.specFullPath())
	if err != nil {
		return err
	}

	log.TheLogger.Info("parsing template files")
	files, err := parse.Parse(cmd.src)
	if err != nil {
		return err
	}
	fdata := log.LoggerData{}
	for _, f := range files {
		fdata[f.Path] = fmt.Sprintf("mode=%v, is-template=%T, byte-counts=%d", f.Mode, f.Template, len(f.Content))
	}
	log.TheLogger.DebugWithData("parsed files", fdata)

	var input common.InputVariables
	if cmd.inputYAML != "" {
		log.TheLogger.InfoWithData("loading template variables from file", log.LoggerData{"location": cmd.inputYAML})
		fileInput, err := spec.InputFromFile(cmd.inputYAML)
		if err != nil {
			return err
		}
		input = fileInput
	} else {
		log.TheLogger.Info("getting template variables dynamically")
		stdinInput, err := spec.InputFromStdIn()
		if err != nil {
			return err
		}
		input = stdinInput
	}

	log.TheLogger.InfoWithData("creating new project", log.LoggerData{"location": cmd.dest})
	err = write.Write(cmd.dest, files, input)
	if err != nil {
		return err
	}

	log.TheLogger.Info("done")
	return nil
}

func (cmd *Command) specFullPath() string {
	return filepath.Join(cmd.src, cmd.specYAML)
}

func (cmd *Command) validate() error {
	if cmd.src == "" {
		return common.NewFlagValidationError(flagSrc, "is required")
	}
	if _, err := os.Stat(cmd.src); err != nil {
		if os.IsNotExist(err) {
			return common.NewFlagValidationError(flagSrc, fmt.Sprintf("(%s) directory not found", cmd.src))
		} else {
			return err
		}
	}

	if cmd.dest == "" {
		return common.NewFlagValidationError(flagDest, "is required")
	}
	if _, err := os.Stat(cmd.dest); err == nil {
		return common.NewFlagValidationError(flagDest, fmt.Sprintf("(%s) directory already exists", cmd.dest))
	}

	if cmd.specYAML == "" {
		return common.NewFlagValidationError(flagSpec, "is required")
	}
	specYAML := cmd.specFullPath()
	if err := fileExists(flagSpec, specYAML); err != nil {
		return err
	}

	if cmd.inputYAML != "" {
		if err := fileExists(flagInput, cmd.inputYAML); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(flagName string, filePath string) error {
	fInfo, err := os.Stat(filePath)
	if err != nil {
		return common.NewFlagValidationError(flagName, fmt.Sprintf("(%s) file not found", filePath))
	}
	if fInfo.IsDir() {
		return common.NewFlagValidationError(flagName, fmt.Sprintf("(%s) is not a file", filePath))
	}
	return nil
}

func sToP(s string) *string {
	return &s
}
