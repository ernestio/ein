/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"github.com/r3labs/ein/builds"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "diff output of a build",
	Long:  "diff output of a build",
	Run:   builds.Diff,
}

func init() {
	RootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringP("from", "f", "", "From build ID")
	diffCmd.Flags().StringP("to", "t", "", "To build ID")
}
