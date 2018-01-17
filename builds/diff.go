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

type component struct {
	created bool
	deleted bool
	changes []diff.Change
}

func (c *component) add(dc diff.Change) {
	(*c).changes = append((*c).changes, dc)
}

type group struct {
	components map[string]*component
}

func (g *group) get(c string) *component {
	if (*g).components[c] == nil {
		(*g).components[c] = &component{}
	}

	return (*g).components[c]
}

type doutput struct {
	groups map[string]*group
}

func (o *doutput) get(g string) *group {
	if (*o).groups[g] == nil {
		(*o).groups[g] = &group{
			components: make(map[string]*component),
		}
	}

	return (*o).groups[g]
}

func (o *doutput) add(c diff.Change) {
	id := c.Path[0]
	c.Path = c.Path[1:]

	parts := strings.Split(id, "::")
	ctype := parts[0]
	compid := parts[1]
	last := c.Path[len(c.Path)-1]

	if o.groups == nil {
		o.groups = make(map[string]*group)
	}

	x := o.get(ctype).get(compid)

	if last == "_component_id" {
		switch c.Type {
		case diff.CREATE:
			x.created = true
		case diff.DELETE:
			x.deleted = true
		}
		return
	}

	x.add(c)
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

	var o doutput

	for _, c := range g.Changelog {
		o.add(c)
	}

	for gname, grp := range o.groups {
		fmt.Printf("\n")

		for cname, cmp := range grp.components {
			var header string
			if cmp.created {
				header = fmt.Sprintf("\n + %s", color.GreenString(gname+"."+cname))
			} else if cmp.deleted {
				header = fmt.Sprintf("\n - %s", color.RedString(gname+"."+cname))
			} else {
				header = fmt.Sprintf("\n ~ %s", color.YellowString(gname+"."+cname))
			}

			fmt.Println(header)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(false)
			table.SetHeaderLine(false)
			table.SetColumnSeparator("")
			table.SetRowLine(false)
			table.SetAutoWrapText(false)

			for _, chng := range cmp.changes {
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
				if cmp.created {
					table.Append([]string{color.GreenString("  " + strings.Join(chng.Path, ".")), `"` + fmt.Sprint(chng.From) + `" => "` + fmt.Sprint(chng.To) + `"`})
				} else if cmp.deleted {
					table.Append([]string{color.RedString("  " + strings.Join(chng.Path, ".")), `"` + fmt.Sprint(chng.From) + `" => "` + fmt.Sprint(chng.To) + `"`})
				} else {
					table.Append([]string{color.YellowString("  " + strings.Join(chng.Path, ".")), `"` + fmt.Sprint(chng.From) + `" => "` + fmt.Sprint(chng.To) + `"`})
				}
			}

			table.Render()
		}
	}
}
