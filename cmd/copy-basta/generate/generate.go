package generate

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate/write"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate/common"

	"github.com/spin14/copy-basta/cmd/copy-basta/generate/specification"
)

const (
	CommandFlag      = "generate"
	CommandFlagShort = "generates new project based on the template and provided parameters"

	flagSrc         = "src"
	flagDest        = "dest"
	flagSpec        = "spec"
	flagInput       = "input"
	flagDefaultSpec = "spec.yaml"
	flagUsageSrc    = "Generated Project root directory"
	flagUsageDest   = "Specification YAML file, relative to the template root directory"
	flagUsageSpec   = "Path to the YAML containing the template specification"
	flagUsageInput  = "Path to the YAML file with the variables to use in the templates"
)

type Flag struct {
	Ref     *string
	Name    string
	Default *string
	Usage   string
}

type Command struct {
	src       string
	dest      string
	specYAML  string
	inputYAML string
}

func (cmd *Command) Flags() []Flag {
	return []Flag{
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

func (cmd *Command) Run() error {
	log.Println("[INFO] Generating new project!")
	if err := cmd.validate(); err != nil {
		return err
	}

	files, err := parse(cmd.src)
	if err != nil {
		return err
	}

	spec, err := specification.New(cmd.specFullPath())
	if err != nil {
		return err
	}

	var input common.InputVariables
	if cmd.inputYAML != "" {
		fileInput, err := spec.InputFromFile(cmd.inputYAML)
		if err != nil {
			return err
		}
		input = fileInput
	} else {
		stdinInput, err := spec.InputFromStdIn()
		if err != nil {
			return err
		}
		input = stdinInput
	}

	err = write.Write(cmd.dest, files, input)
	if err != nil {
		return err
	}

	log.Println("[INFO] Done!")
	return nil
}

func (cmd *Command) specFullPath() string {
	return path.Join(cmd.src, cmd.specYAML)
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
