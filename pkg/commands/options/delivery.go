package options

import "github.com/spf13/cobra"

// DeliveryOptions
type DeliveryOptions struct {
	Delay string
}

func AddDeliveryArg(cmd *cobra.Command, po *DeliveryOptions) {
	cmd.Flags().StringVar(&po.Delay, "delay", "",
		"Delay between sending events such as `300ms`. Valid time units are `ns`, `us`, `ms`, `s`, `m`, `h`.")
}
