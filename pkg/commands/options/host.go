package options

import (
	"github.com/spf13/cobra"
)

// HostOptions
type HostOptions struct {
	URL string
}

func AddHostArg(cmd *cobra.Command, to *HostOptions) {
	cmd.Flags().StringVarP(&to.URL, "host", "t", "",
		"Target a host to invoke. Must be a fully qualified URL.")
}
