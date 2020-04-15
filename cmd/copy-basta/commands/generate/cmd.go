package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate/specification"

	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate/parse"
	"github.com/spin14/copy-basta/cmd/copy-basta/commands/generate/write"
	"github.com/spin14/copy-basta/cmd/copy-basta/common"
)

const (
	commandID          = "generate"
	commandDescription = "generates new project based on the template and provided parameters"

	flagSrc         = "src"
	flagDest        = "dest"
	flagSpec        = "spec"
	flagInput       = "input"
	flagDefaultSpec = common.SpecFile
	flagUsageSrc    = "Generated Project root directory"
	flagUsageDest   = "Specification YAML file, relative to the template root directory"
	flagUsageSpec   = "Path to the YAML containing the template specification"
	flagUsageInput  = "Path to the YAML file with the variables to use in the templates"
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
			Usage:   flagUsageSrc,
		},
		{
			Ref:     &cmd.dest,
			Name:    flagDest,
			Default: nil,
			Usage:   flagUsageDest,
		},
		{
			Ref:     &cmd.specYAML,
			Name:    flagSpec,
			Default: sToP(flagDefaultSpec),
			Usage:   flagUsageSpec,
		},
		{
			Ref:     &cmd.inputYAML,
			Name:    flagInput,
			Default: nil,
			Usage:   flagUsageInput,
		},
	}
}

func (cmd *Command) Run(logger *common.Logger) error {
	logger.DebugWithData("user input", common.LoggerData{
		flagSrc:   cmd.src,
		flagDest:  cmd.dest,
		flagSpec:  cmd.specYAML,
		flagInput: cmd.inputYAML,
	})
	logger.Info("validating user input")
	if err := cmd.validate(); err != nil {
		return err
	}

	logger.Info("loading specification file")
	spec, err := specification.New(cmd.specFullPath())
	if err != nil {
		return err
	}

	logger.Info("parsing template files")
	files, err := parse.Parse(cmd.src)
	if err != nil {
		return err
	}
	fdata := common.LoggerData{}
	for _, f := range files {
		fdata[f.Path] = fmt.Sprintf("mode=%v, is-template=%T, byte-counts=%d", f.Mode, f.Template, len(f.Content))
	}
	logger.DebugWithData("parsed files", fdata)

	var input common.InputVariables
	if cmd.inputYAML != "" {
		logger.InfoWithData("loading template variables from file", common.LoggerData{"filepath": cmd.inputYAML})
		fileInput, err := spec.InputFromFile(cmd.inputYAML)
		if err != nil {
			return err
		}
		input = fileInput
	} else {
		logger.Info("getting template variables dynamically")
		stdinInput, err := spec.InputFromStdIn()
		if err != nil {
			return err
		}
		input = stdinInput
	}

	logger.Info("creating new project")
	err = write.Write(cmd.dest, files, input)
	if err != nil {
		return err
	}

	logger.Info("done")
	return nil
}

func (cmd *Command) specFullPath() string {
	return filepath.Join(cmd.src, cmd.specYAML)
}

func (cmd *Command) validate() error {
	if cmd.src == "" {
		return fmt.Errorf(`[ERROR] "%s" is required`, flagSrc)
	}

	if cmd.dest == "" {
		return fmt.Errorf(`[ERROR] "%s" is required`, flagDest)
	}
	if _, err := os.Stat(cmd.dest); err == nil {
		return fmt.Errorf(`[ERROR] "%s" (%s) already exists`, flagDest, cmd.dest)
	}

	if cmd.specYAML == "" {
		return fmt.Errorf(`[ERROR] "%s" is required`, flagSpec)
	}
	spec := cmd.specFullPath()
	if err := fileExistsOrErr(spec, flagSpec); err != nil {
		return err
	}

	if cmd.inputYAML != "" {
		if err := fileExistsOrErr(cmd.inputYAML, flagInput); err != nil {
			return err
		}
	}

	return nil
}

func fileExistsOrErr(filePath string, name string) error {
	fInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf(`[ERROR] "%s" (%s) not found`, name, filePath)
	}
	if fInfo.IsDir() {
		return fmt.Errorf(`[ERROR] "%s" (%s) is not a file`, name, filePath)
	}
	return nil
}

func sToP(s string) *string {
	return &s
}
