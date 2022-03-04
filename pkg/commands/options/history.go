/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package options

import "github.com/spf13/cobra"

// HistoryOptions
type HistoryOptions struct {
	Length int
	Retain bool
}

func (o *HistoryOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&o.Length, "history", 50,
		"How many past events to store.")

	cmd.Flags().BoolVar(&o.Retain, "retain", true,
		"Events are retained in history if history is collected.")
}
