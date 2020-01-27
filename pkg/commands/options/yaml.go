package options

import "github.com/spf13/cobra"

// YAMLOptions
type YAMLOptions struct {
	YAML bool
}

func AddYAMLArg(cmd *cobra.Command, yo *YAMLOptions) {
	cmd.Flags().BoolVar(&yo.YAML, "yaml", false,
		"Output as YAML.")
}
