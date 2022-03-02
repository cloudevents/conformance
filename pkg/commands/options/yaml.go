/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package options

import "github.com/spf13/cobra"

// YAMLOptions
type YAMLOptions struct {
	YAML bool
}

func (o *YAMLOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.YAML, "yaml", false,
		"Output as YAML.")
}
