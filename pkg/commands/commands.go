package commands

import (
	"github.com/spf13/cobra"
)

func AddConformanceCommands(topLevel *cobra.Command) {
	addSend(topLevel)
	addInvoke(topLevel)
	addListener(topLevel)
}
