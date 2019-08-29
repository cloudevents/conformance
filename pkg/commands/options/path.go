package options

import "github.com/spf13/cobra"

// PathOptions
type PathOptions struct {
	Path string
}

func AddPathArg(cmd *cobra.Command, po *PathOptions) {
	cmd.Flags().StringVarP(&po.Path, "path", "p", "/",
		"Path to use")
}
