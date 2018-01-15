/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"fmt"
	"os"
	"strings"

	"github.com/ernestio/mapping"
	"github.com/fatih/color"
	"github.com/nats-io/nats"
	"github.com/olekukonko/tablewriter"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
	"github.com/spf13/cobra"
)

//type components map[string]

type change struct {
	created bool
	deleted bool
	changes []diff.Change
}

type coutput map[string]map[string][]diff.Change

func (cl *coutput) add(c diff.Change) {
	id := c.Path[0]
	c.Path = c.Path[1:]

	parts := strings.Split(id, "::")
	if (*cl)[parts[0]] == nil {
		(*cl)[parts[0]] = make(map[string][]diff.Change)
	}

	if (*cl)[parts[0]][parts[1]] == nil {
		(*cl)[parts[0]][parts[1]] = make([]diff.Change, 0)
	}

	(*cl)[parts[0]][parts[1]] = append((*cl)[parts[0]][parts[1]], c)
}

func Diff(cmd *cobra.Command, args []string) {
	nuri, _ := cmd.Flags().GetString("nats")
	envName, _ := cmd.Flags().GetString("env")
	fromid, _ := cmd.Flags().GetString("from")
	toid, _ := cmd.Flags().GetString("to")

	nc, err := nats.Connect(nuri)
	if err != nil {
		panic(err)
	}

	m := mapping.New(nc, envName)

	err = m.Diff(fromid, toid)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	g := graph.New()
	err = g.Load(m.Result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	co := make(coutput)

	for _, c := range g.Changelog {
		co.add(c)
	}

	for _, components := range co {
		fmt.Printf("\n")

		for name, changes := range components {
			fmt.Println("\n" + name)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(false)
			table.SetHeaderLine(false)
			table.SetColumnSeparator("")
			table.SetRowLine(false)
			table.SetAutoWrapText(false)

			for _, chng := range changes {
				if chng.From == nil {
					chng.From = ""
					chng.To = color.GreenString(fmt.Sprint(chng.To))
				} else if chng.To == nil {
					chng.To = ""
					chng.From = color.RedString(fmt.Sprint(chng.From))
				} else {
					chng.From = color.RedString(fmt.Sprint(chng.From))
					chng.To = color.GreenString(fmt.Sprint(chng.To))
				}
				table.Append([]string{strings.Join(chng.Path, "."), `"` + fmt.Sprint(chng.From) + `" => "` + fmt.Sprint(chng.To) + `"`})
			}

			table.Render()
		}
	}
}
