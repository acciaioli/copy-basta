package commands

import (
	"copy-basta/services/generate"

	"github.com/spf13/cobra"

	"copy-basta/internal/common"
)

func Generate(globals func() error) *cobra.Command {

	const (
		commandUse         = "generate"
		commandDescription = "generates a new codebase based on the template and on the provided variables"

		flagSrc            = "src"
		flagDescriptionSrc = "root directory of the template codebase"

		flagDest            = "dest"
		flagDescriptionDest = "to be root directory of the generated codebase"

		flagSpec            = "spec"
		flagDefaultSpec     = common.SpecFile
		flagDescriptionSpec = "path (relative to src) to the YAML containing the template specification"

		flagInput            = "input"
		flagDescriptionInput = "path to the YAML file with the variables to use in the templates"

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
