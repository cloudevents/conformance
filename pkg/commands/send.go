package commands

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/sender"
)

func addSend(topLevel *cobra.Command) {
	ho := &options.HostOptions{}
	eo := &options.EventOptions{}
	do := &options.DeliveryOptions{}
	yo := &options.YAMLOptions{}
	vo := &options.VerboseOptions{}

	invoke := &cobra.Command{
		Use:   "send",
		Short: "Send a cloudevent.",
		Example: `
  cloudevents send http://localhost:8080/ --id abc-123 --source cloudevents.conformance.tool --type foo.bar
  cloudevents send http://localhost:8080/ --id 321-cba --source cloudevents.conformance.tool --type foo.json --mode structured
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
			// Process time now.
			if eo.Now {
				eo.Event.Attributes.Time = time.Now().UTC().Format(time.RFC3339Nano)
			}

			// Process extensions.
			if len(eo.Extensions) > 0 {
				eo.Event.Attributes.Extensions = make(map[string]string)
				for _, ext := range eo.Extensions {
					kv := strings.SplitN(ext, "=", 2)
					if len(kv) == 2 {
						eo.Event.Attributes.Extensions[kv[0]] = kv[1]
					}
				}
			}

			// Build up command.
			i := &sender.Sender{
				URL:     ho.URL,
				Event:   eo.Event,
				YAML:    yo.YAML,
				Verbose: vo.Verbose,
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
			return i.Do()
		},
	}
	eo.AddFlags(invoke)
	yo.AddFlags(invoke)
	vo.AddFlags(invoke)
	do.AddFlags(invoke)

	topLevel.AddCommand(invoke)
}
