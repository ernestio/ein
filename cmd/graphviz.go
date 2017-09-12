/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"github.com/r3labs/ein/builds"
	"github.com/spf13/cobra"
)

// graphvizCmd represents the graphviz command
var graphvizCmd = &cobra.Command{
	Use:   "graphviz",
	Short: "graphviz output of a build",
	Long:  "graphviz output of a build",
	Run:   builds.Graphviz,
}

func init() {
	RootCmd.AddCommand(graphvizCmd)
	graphvizCmd.Flags().StringP("component-type", "t", "", "Filter by Component Type")
	graphvizCmd.Flags().StringP("component-tags", "x", "", "Filter by Component Tag")
}
