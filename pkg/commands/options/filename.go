package options

import (
	"github.com/spf13/cobra"
)

// FilenameOptions
type FilenameOptions struct {
	Filenames []string
	Recursive bool
}

func AddFilenameArg(cmd *cobra.Command, fo *FilenameOptions) {
	cmd.Flags().StringSliceVarP(&fo.Filenames, "filename", "f", fo.Filenames,
		"Filename or directory to use")
	cmd.Flags().BoolVarP(&fo.Recursive, "recursive", "R", fo.Recursive,
		"Process the directory used in -f, --filename recursively.")
}
