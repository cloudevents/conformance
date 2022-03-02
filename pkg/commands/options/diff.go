/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package options

import "github.com/spf13/cobra"

// DiffOptions
type DiffOptions struct {
	FindBy []string

	FullDiff        bool
	IgnoreAdditions bool
}

func (o *DiffOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&o.FindBy, "match", []string{"id"},
		"Find by keys to compare each event set.")

	cmd.Flags().BoolVar(&o.FullDiff, "full", false,
		"Print the full diff.")

	cmd.Flags().BoolVar(&o.IgnoreAdditions, "ignore-additions", false,
		"Ignore additions between file_a and file_b.")
}
