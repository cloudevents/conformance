package options

import "github.com/spf13/cobra"

// VerboseOptions
type VerboseOptions struct {
	Verbose bool
}

func (o *VerboseOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&o.Verbose, "verbose", "v", false,
		"Output more debug info to stderr")
}
