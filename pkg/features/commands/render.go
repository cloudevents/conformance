/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"github.com/cloudevents/conformance/pkg/features/model"
	"github.com/cloudevents/conformance/pkg/features/renderer"
	"github.com/spf13/cobra"
)

func addRender(topLevel *cobra.Command) {
	var input []model.Features
	render := &cobra.Command{
		Use:   "render FEATURE_FILE...",
		Short: "Render the markdown table for the given feature file(s).",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			input, err = model.FromYaml(true, args...)
			return
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Build up command.
			i := &renderer.Renderer{
				// TODO: we might have formatting options.
			}

			// Run it.
			return i.Render(input, cmd.OutOrStdout())
		},
	}

	topLevel.AddCommand(render)
}
