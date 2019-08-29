package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/commands"
)

func main() {
	cmds := &cobra.Command{
		Use:   "cloudevents",
		Short: "A tool to aid in CloudEvents conformance testing.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	commands.AddConformanceCommands(cmds)

	if err := cmds.Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}
