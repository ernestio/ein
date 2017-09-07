/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"github.com/r3labs/ein/builds"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a build",
	Long:  "get a build",
	Run:   builds.Get,
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("component-type", "t", "", "Filter by Component Type")
	getCmd.Flags().StringP("component-tags", "x", "", "Filter by Component Tag")
}
