/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"fmt"
	"github.com/cloudevents/conformance/pkg/features/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func addTemplate(topLevel *cobra.Command) {
	render := &cobra.Command{
		Use:   "template LANGUAGE",
		Short: "Produce an example template of the sdk features yaml config.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			feat := model.Features{
				Language: args[0],
				Versions: map[string]model.VersionedFeatures{model.Version1dot0: {}},
			}

			out, err := yaml.Marshal(feat)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), string(out))

			return nil
		},
	}

	topLevel.AddCommand(render)
}
