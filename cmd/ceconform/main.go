package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/commands"
)

func main() {
	cmds := &cobra.Command{
		Use:   "ceconform",
		Short: "A tool to aid in CloudEvents conformance testing.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	fmt.Println("wat")
	commands.AddConformanceCommands(cmds)

	if err := cmds.Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}
