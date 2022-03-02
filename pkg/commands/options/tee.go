/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

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

func (o *TeeOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&o.URLString, "tee", "t", "",
		"Tee to a host url")
}
