/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package options

import (
	"github.com/spf13/cobra"
)

// FilenameOptions
type FilenameOptions struct {
	Filenames []string
	Recursive bool
}

func (o *FilenameOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceVarP(&o.Filenames, "filename", "f", o.Filenames,
		"Filename or directory to use")
	cmd.Flags().BoolVarP(&o.Recursive, "recursive", "R", o.Recursive,
		"Process the directory used in -f, --filename recursively.")
}
