package commands

import (
	"github.com/spf13/cobra"
	"log"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/invoker"
)

func addInvoke(topLevel *cobra.Command) {
	ho := &options.HostOptions{}
	fo := &options.FilenameOptions{}
	invoke := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke the host with the example input files.",
		Example: `
  ceconform invoke -t http://localhost:8008/ -f ./yaml/v0.3
`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Build up command.
			i := &invoker.Invoker{
				URL:       ho.URL,
				Files:     fo.Filenames,
				Recursive: fo.Recursive,
			}

			// Run it.
			if err := i.Do(); err != nil {
				log.Fatalf("error invoking target: %v", err)
			}
		},
	}
	options.AddHostArg(invoke, ho)
	options.AddFilenameArg(invoke, fo)

	topLevel.AddCommand(invoke)
}
