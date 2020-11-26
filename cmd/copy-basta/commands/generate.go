package commands

import (
	"copy-basta/services/generate"

	"github.com/spf13/cobra"

	"copy-basta/internal/common"
)

func GenerateCommand(globals func() error) *cobra.Command {

	const (
		commandUse         = "generate"
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

		flagOverwrite            = "overwrite"
		flagDescriptionOverwrite = "Allow overriding files in an existing destination directory"
	)

	var src string
	var dest string
	var specYAML string
	var inputYAML string
	var overwrite bool

	cmd := &cobra.Command{
		Use:   commandUse,
		Short: commandDescription,
		RunE: func(cmd2 *cobra.Command, what []string) error {
			err := globals()
			if err != nil {
				return err
			}
			return generate.Generate(&generate.Params{
				Src:       src,
				Dest:      dest,
				SpecYAML:  specYAML,
				InputYAML: inputYAML,
				Overwrite: overwrite,
			})
		},
	}

	cmd.Flags().StringVar(
		&src,
		flagSrc,
		"",
		flagDescriptionSrc,
	)

	cmd.Flags().StringVar(
		&dest,
		flagDest,
		"",
		flagDescriptionDest,
	)

	cmd.Flags().StringVar(
		&specYAML,
		flagSpec,
		flagDefaultSpec,
		flagDescriptionSpec,
	)

	cmd.Flags().StringVar(
		&inputYAML,
		flagInput,
		"",
		flagDescriptionInput,
	)

	cmd.Flags().BoolVar(
		&overwrite,
		flagOverwrite,
		false,
		flagDescriptionOverwrite,
	)

	return cmd
}
