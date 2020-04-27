package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"copy-basta/cmd/copy-basta/clients/github"
	"copy-basta/cmd/copy-basta/common"
	"copy-basta/cmd/copy-basta/common/log"
	"copy-basta/cmd/copy-basta/generate/parse"
	"copy-basta/cmd/copy-basta/generate/parse/parsediskloader"
	"copy-basta/cmd/copy-basta/generate/parse/parsegithubloader"
	"copy-basta/cmd/copy-basta/generate/specification"
	"copy-basta/cmd/copy-basta/generate/specification/specificationdiskloader"
	"copy-basta/cmd/copy-basta/generate/specification/specificationgithubloader"
	"copy-basta/cmd/copy-basta/generate/write"
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

	ghc *github.Client
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
	log.L.DebugWithData("user input", log.Data{
		flagSrc:   cmd.src,
		flagDest:  cmd.dest,
		flagSpec:  cmd.specYAML,
		flagInput: cmd.inputYAML,
	})
	log.L.Info("validating user input")
	if err := cmd.validate(); err != nil {
		return err
	}

	log.L.Info("loading specification file")
	specLoader, err := cmd.getSpecificationLoader()
	if err != nil {
		return err
	}
	spec, err := specification.New(specLoader)
	if err != nil {
		return err
	}

	log.L.Info("parsing template files")
	parseLoader, err := cmd.getParseLoader()
	if err != nil {
		return err
	}
	files, err := parse.Parse(parseLoader)
	if err != nil {
		return err
	}
	fdata := log.Data{}
	for _, f := range files {
		fdata[f.Path] = fmt.Sprintf("mode=%v, is-template=%T, byte-counts=%d", f.Mode, f.Template, len(f.Content))
	}
	log.L.DebugWithData("parsed files", fdata)

	var input common.InputVariables
	if cmd.inputYAML != "" {
		log.L.InfoWithData("loading template variables from file", log.Data{"location": cmd.inputYAML})
		fileInput, err := spec.InputFromFile(cmd.inputYAML)
		if err != nil {
			return err
		}
		input = fileInput
	} else {
		log.L.Info("getting template variables dynamically")
		stdinInput, err := spec.InputFromStdIn()
		if err != nil {
			return err
		}
		input = stdinInput
	}

	log.L.InfoWithData("creating new project", log.Data{"location": cmd.dest})
	err = write.Write(cmd.dest, files, input)
	if err != nil {
		return err
	}

	log.L.Info("done")
	return nil
}

func (cmd *Command) specFullPath() string {
	return filepath.Join(cmd.src, cmd.specYAML)
}

func (cmd *Command) validate() error {
	if cmd.src == "" {
		return common.NewFlagValidationError(flagSrc, "is required")
	}

	if strings.HasPrefix(cmd.src, common.GithubPrefix) {
		// remote github
		ghc, err := github.NewClient(strings.TrimPrefix(cmd.src, common.GithubPrefix))
		if err != nil {
			return err
		}
		cmd.ghc = ghc
		// todo: maybe we should rethink how validations are handled
		log.L.Warn("src is a remote location, skipping validations")
	} else {
		// local
		if _, err := os.Stat(cmd.src); err != nil {
			if os.IsNotExist(err) {
				return common.NewFlagValidationError(flagSrc, fmt.Sprintf("(%s) directory not found", cmd.src))
			}
			return err
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
	}

	if cmd.dest == "" {
		return common.NewFlagValidationError(flagDest, "is required")
	}
	if _, err := os.Stat(cmd.dest); err == nil {
		return common.NewFlagValidationError(flagDest, fmt.Sprintf("(%s) directory already exists", cmd.dest))
	}

	return nil
}

func (cmd *Command) getParseLoader() (parse.Loader, error) {
	switch {
	case strings.HasPrefix(cmd.src, common.GithubPrefix):
		log.L.Debug("using github loader for parsing")
		return parsegithubloader.New(cmd.ghc)
	default:
		log.L.Debug("using disk loader for parsing")
		return parsediskloader.NewLoader(cmd.src)
	}
}

func (cmd *Command) getSpecificationLoader() (specification.Loader, error) {
	switch {
	case strings.HasPrefix(cmd.src, common.GithubPrefix):
		log.L.Debug("using github loader for specification")
		return specificationgithubloader.New(cmd.specYAML, cmd.ghc)
	default:
		log.L.Debug("using disk loader for specification")
		return specificationdiskloader.New(cmd.specFullPath())
	}
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
