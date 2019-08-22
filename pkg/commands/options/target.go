package options

import (
	"github.com/spf13/cobra"
)

// TargetOptions
type TargetOptions struct {
	Target string
}

func AddTargetArg(cmd *cobra.Command, to *TargetOptions) {
	cmd.Flags().StringVarP(&to.Target, "target", "t", "",
		"Target a host to invoke. Must be a URL.")
}
