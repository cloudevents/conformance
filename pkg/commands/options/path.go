package options

import "github.com/spf13/cobra"

// PathOptions
type PathOptions struct {
	Path string
}

func (o *PathOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.Path, "path", "p", "/",
		"Path to use")
}
