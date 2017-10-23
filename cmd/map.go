/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"github.com/r3labs/ein/builds"
	"github.com/spf13/cobra"
)

// mapCmd represents the map command
var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "map a build",
	Long:  "map a build",
	Run:   builds.Map,
}

func init() {
	RootCmd.AddCommand(mapCmd)
	mapCmd.Flags().BoolP("graphviz", "g", false, "Show graphviz output")
}
