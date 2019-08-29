package options

import "github.com/spf13/cobra"

// VerboseOptions
type VerboseOptions struct {
	Verbose bool
}

func AddVerboseArg(cmd *cobra.Command, vo *VerboseOptions) {
	cmd.Flags().BoolVarP(&vo.Verbose, "verbose", "v", false,
		"Output more debug info to stderr")
}
