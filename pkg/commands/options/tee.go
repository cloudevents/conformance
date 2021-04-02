package options

import (
	"net/url"

	"github.com/spf13/cobra"
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
