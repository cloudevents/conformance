/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"github.com/spf13/cobra"
)

func AddConformanceCommands(topLevel *cobra.Command) {
	addSend(topLevel)
	addInvoke(topLevel)
	addListener(topLevel)
	addRaw(topLevel)
}
