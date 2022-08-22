/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package commands

import (
	"github.com/spf13/cobra"
)

func AddFeaturesCommands(topLevel *cobra.Command) {
	addRender(topLevel)
	addTemplate(topLevel)
}
