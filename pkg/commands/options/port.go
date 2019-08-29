package options

import "github.com/spf13/cobra"

// PortOptions
type PortOptions struct {
	Port int
}

func AddPortArg(cmd *cobra.Command, po *PortOptions) {
	cmd.Flags().IntVarP(&po.Port, "port", "P", 8080,
		"Port to use")
}
