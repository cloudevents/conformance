package commands

import (
	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/commands/options"
)

func addInvoke(topLevel *cobra.Command) {
	to := &options.TargetOptions{}
	fo := &options.FilenameOptions{}
	invoke := &cobra.Command{
		Use:   "invoke -t <host> -f <file>",
		Short: "Invoke the host with the example input files.",
		Long:  `TODO`,
		Example: `
  TODO
`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Build up command.

			// Run it.

		},
	}
	options.AddTargetArg(invoke, to)
	options.AddFilenameArg(invoke, fo)

	topLevel.AddCommand(invoke)
}
