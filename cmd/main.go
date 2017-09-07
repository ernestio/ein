/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import "github.com/spf13/cobra"

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:   "ein",
	Short: "Ein is a graph inspection tool for ernest",
	Long:  "Ein is a graph inspection tool for ernest",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	RootCmd.PersistentFlags().StringP("nats", "n", "nats://127.0.0.1:4222", "NATS URI")
	RootCmd.PersistentFlags().StringP("build", "b", "", "Build ID")
	RootCmd.PersistentFlags().StringP("env", "e", "", "Environment Name")
}
