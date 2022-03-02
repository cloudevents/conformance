/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/http"
)

func addRaw(topLevel *cobra.Command) {
	po := &options.PortOptions{}
	raw := &cobra.Command{
		Use:   "raw",
		Short: "Dump the raw HTTP request to stdout. (For debugging HTTP requests.)",
		Example: `
  cloudevents raw -P 8181
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			r := http.Raw{
				Out:  cmd.OutOrStdout(),
				Port: po.Port,
			}
			return r.Do()
		},
	}
	po.AddFlags(raw)

	topLevel.AddCommand(raw)
}
