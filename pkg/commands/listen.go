package commands

import (
	"github.com/spf13/cobra"
	"log"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/listener"
)

func addListener(topLevel *cobra.Command) {
	po := &options.PortOptions{}
	pa := &options.PathOptions{}
	vo := &options.VerboseOptions{}
	listen := &cobra.Command{
		Use:   "listen",
		Short: "Listen to incoming http CloudEvents requests and write out the yaml representation to stdout.",
		Example: `
  cloudevents listen -P 8080 -p incoming -v > got.yaml
`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Build up command.
			i := &listener.Listener{
				Port:    po.Port,
				Path:    pa.Path,
				Verbose: vo.Verbose,
			}

			// Run it.
			if err := i.Do(); err != nil {
				log.Fatalf("error listening: %v", err)
			}
		},
	}
	options.AddPortArg(listen, po)
	options.AddPathArg(listen, pa)
	options.AddVerboseArg(listen, vo)

	topLevel.AddCommand(listen)
}
