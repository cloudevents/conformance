/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"log"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/listener"
)

func addListener(topLevel *cobra.Command) {
	po := &options.PortOptions{}
	pa := &options.PathOptions{}
	to := &options.TeeOptions{}
	ho := &options.HistoryOptions{}
	vo := &options.VerboseOptions{}
	listen := &cobra.Command{
		Use:   "listen",
		Short: "Listen to incoming http CloudEvents requests and write out the yaml representation to stdout.",
		Example: `
  cloudevents listen -P 8080 -p incoming -v > got.yaml
`,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if to.URLString != "" {
				u, err := url.Parse(to.URLString)
				if err != nil {
					return err
				}
				to.URL = u
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Build up command.
			i := &listener.Listener{
				Port:    po.Port,
				Path:    pa.Path,
				Tee:     to.URL,
				History: ho.Length,
				Retain:  ho.Retain,
				Verbose: vo.Verbose,
			}

			// Run it.
			if err := i.Do(cmd.Context()); err != nil {
				log.Fatalf("error listening: %v", err)
			}
		},
	}
	po.AddFlags(listen)
	pa.AddFlags(listen)
	vo.AddFlags(listen)
	ho.AddFlags(listen)
	to.AddFlags(listen)

	topLevel.AddCommand(listen)
}
