package commands

import (
	"errors"
	"net/url"
	"time"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/invoker"
	"github.com/spf13/cobra"
)

func addInvoke(topLevel *cobra.Command) {
	ho := &options.HostOptions{}
	fo := &options.FilenameOptions{}
	do := &options.DeliveryOptions{}
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Build up command.
			i := &invoker.Invoker{
				URL:       ho.URL,
				Files:     fo.Filenames,
				Recursive: fo.Recursive,
				Verbose:   vo.Verbose,
			}

			// Add delay, if specified.
			if len(do.Delay) > 0 {
				d, err := time.ParseDuration(do.Delay)
				if err != nil {
					return err
				}
				i.Delay = &d
			}

			// Run it.
			if err := i.Do(); err != nil {
				return err
			}
			return nil
		},
	}
	options.AddFilenameArg(invoke, fo)
	options.AddDeliveryArg(invoke, do)
	options.AddVerboseArg(invoke, vo)

	topLevel.AddCommand(invoke)
}
