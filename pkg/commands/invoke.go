package commands

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"net/url"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/invoker"
)

func addInvoke(topLevel *cobra.Command) {
	ho := &options.HostOptions{}
	fo := &options.FilenameOptions{}
	vo := &options.VerboseOptions{}
	invoke := &cobra.Command{
		Use:   "invoke",
		Short: "Invoke the host with the example input files.",
		Example: `
  cloudevents invoke http://localhost:8008/ -f ./yaml/v1.0
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a host argument")
			}
			u, err := url.Parse(args[0])
			if err != nil {
				return err
			}
			ho.URL = u
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Build up command.
			i := &invoker.Invoker{
				URL:       ho.URL,
				Files:     fo.Filenames,
				Recursive: fo.Recursive,
				Verbose:   vo.Verbose,
			}

			// Run it.
			if err := i.Do(); err != nil {
				log.Fatalf("error invoking target: %v", err)
			}
		},
	}
	options.AddFilenameArg(invoke, fo)
	options.AddVerboseArg(invoke, vo)

	topLevel.AddCommand(invoke)
}
