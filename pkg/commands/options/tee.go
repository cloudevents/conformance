package options

import (
	"github.com/spf13/cobra"
	"net/url"
)

// TeeOptions
type TeeOptions struct {
	URLString string
	URL       *url.URL
}

func AddTeeArg(cmd *cobra.Command, to *TeeOptions) {
	cmd.Flags().StringVarP(&to.URLString, "tee", "t", "",
		"Tee to a host url")
}
