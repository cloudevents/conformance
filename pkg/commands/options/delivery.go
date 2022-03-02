package options

import "github.com/spf13/cobra"

// DeliveryOptions
type DeliveryOptions struct {
	Delay string
}

func (o *DeliveryOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Delay, "delay", "",
		"Delay between sending events such as `300ms`. Valid time units are `ns`, `us`, `ms`, `s`, `m`, `h`.")
}
