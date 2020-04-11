package generate

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := Command{}

	cobraCmd := &cobra.Command{
		Use:   "generate",
		Short: "generates new project based on the template and provided parameters",
		RunE: func(*cobra.Command, []string) error {
			return cmd.run()
		},
	}

	cobraCmd.Flags().StringVar(
		&cmd.src,
		"src",
		"",
		`Template root directory`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.dest,
		"dest",
		"",
		`Project name (root directory for the generation output)`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.spec,
		"spec",
		"spec.yaml",
		`Specification YAML file, relative to the template root directory`,
	)
	cobraCmd.Flags().StringVar(
		&cmd.bastaTemplateVarsYAML,
		"basta-yaml",
		"",
		`Path to the YAML file with the variables to use in the templates`,
	)

	return cobraCmd
}
