/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"github.com/r3labs/ein/builds"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all builds",
	Long:  "list all builds",
	Run:   builds.List,
}

func init() {
	RootCmd.AddCommand(listCmd)
}
