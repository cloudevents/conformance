/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/cloudevents/conformance/pkg/features/commands"
	"github.com/spf13/cobra"
)

func main() {
	cmds := &cobra.Command{
		Use:   "features",
		Short: "A tool to generate feature tables for CloudEvents SDKs.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		SilenceUsage: true,
	}
	commands.AddFeaturesCommands(cmds)

	if err := cmds.Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}
