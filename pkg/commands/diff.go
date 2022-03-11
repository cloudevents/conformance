/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/commands/options"
	"github.com/cloudevents/conformance/pkg/diff"
)

func addDiff(topLevel *cobra.Command) {
	do := &options.DiffOptions{}

	cmd := &cobra.Command{
		Use:   "diff file_a file_b",
		Short: "Compare differences between two sets of event data.",
		Example: `
  cloudevents diff ./want.yaml ./got.yaml
`,
		Args: cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			var findBy []string
			for _, fb := range do.FindBy {
				if strings.Contains(fb, ",") {
					findBy = append(findBy, strings.Split(fb, ",")...)
				} else {
					findBy = append(findBy, fb)
				}
			}
			do.FindBy = findBy
		},
		Run: func(cmd *cobra.Command, args []string) {
			r := diff.Diff{
				Out:             cmd.OutOrStdout(),
				FileA:           args[0],
				FileB:           args[1],
				FindBy:          do.FindBy,
				FullDiff:        do.FullDiff,
				IgnoreAdditions: do.IgnoreAdditions,
			}
			if err := r.Do(); err != nil {
				os.Exit(1)
			}
		},
	}

	do.AddFlags(cmd)

	topLevel.AddCommand(cmd)
}
