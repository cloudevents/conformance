package commands

import (
	"errors"
	"log"
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
	yo := &options.YAMLOptions{}
	vo := &options.VerboseOptions{}
	invoke := &cobra.Command{
		Use:   "send",
		Short: "Send a cloudevent.",
		Example: `
  cloudevents send http://localhost:8080/ --id abc-123 --source cloudevents.conformance.tool --type foo.bar
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

			// Run it.
			if err := i.Do(); err != nil {
				log.Fatalf("error sending: %v", err)
			}
		},
	}
	options.AddEventArgs(invoke, eo)
	options.AddYAMLArg(invoke, yo)
	options.AddVerboseArg(invoke, vo)

	topLevel.AddCommand(invoke)
}
