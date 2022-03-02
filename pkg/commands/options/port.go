package options

import "github.com/spf13/cobra"

// PortOptions
type PortOptions struct {
	Port int
}

func (o *PortOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&o.Port, "port", "P", 8080,
		"Port to use")
}
