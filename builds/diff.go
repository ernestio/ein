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
	"github.com/r3labs/graph"
	"github.com/spf13/cobra"
)

type change struct {
	To, From string
}

// type :: component :: changes
type changes map[string]map[string]map[string]change

func (cs *changes) addTo(ctype, name, path, to string) {
	if (*cs)[ctype] == nil {
		(*cs)[ctype] = make(map[string]map[string]change)
	}

	if (*cs)[ctype][name] == nil {
		(*cs)[ctype][name] = make(map[string]change)
	}

	_, ok := (*cs)[ctype][name][path]
	if ok {
		x := (*cs)[ctype][name][path]
		x.To = to
		(*cs)[ctype][name][path] = x
	} else {
		(*cs)[ctype][name][path] = change{To: to}
	}
}

func (cs *changes) addFrom(ctype, name, path, from string) {
	if (*cs)[ctype] == nil {
		(*cs)[ctype] = make(map[string]map[string]change)
	}

	if (*cs)[ctype][name] == nil {
		(*cs)[ctype][name] = make(map[string]change)
	}

	_, ok := (*cs)[ctype][name][path]
	if ok {
		x := (*cs)[ctype][name][path]
		x.From = from
		(*cs)[ctype][name][path] = x
	} else {
		(*cs)[ctype][name][path] = change{From: from}
	}
}

func (cs *changes) processFrom(path []string, v interface{}) {
	if v == nil {
		v = "-"
	}

	switch v.(type) {
	case string:
		info := strings.Split(path[0], "::")
		cs.addFrom(info[0], info[1], strings.Join(path[1:], "."), v.(string))
	case map[string]interface{}:
		m := v.(map[string]interface{})
		for k, v := range m {
			cs.processFrom(append(path, k), v)
		}
	default:
		info := strings.Split(path[0], "::")
		cs.addFrom(info[0], info[1], strings.Join(path[1:], "."), fmt.Sprint(v))
	}
}

func (cs *changes) processTo(path []string, v interface{}) {
	if v == nil {
		v = "-"
	}

	switch v.(type) {
	case string:
		info := strings.Split(path[0], "::")
		cs.addTo(info[0], info[1], strings.Join(path[1:], "."), v.(string))
	case map[string]interface{}:
		m := v.(map[string]interface{})
		for k, v := range m {
			cs.processTo(append(path, k), v)
		}
	default:
		info := strings.Split(path[0], "::")
		cs.addTo(info[0], info[1], strings.Join(path[1:], "."), fmt.Sprint(v))
	}
}

func skip(s string) bool {
	return s == "" || s == "-"
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

	cs := make(changes)

	for _, c := range g.Changelog {
		cs.processFrom(c.Path, c.From)
		cs.processTo(c.Path, c.To)
	}

	for t, components := range cs {
		fmt.Println(strings.ToUpper(t) + "'s")

		for name, values := range components {
			fmt.Println("\n" + name)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Field", "From", "To"})

			for path, chng := range values {
				if skip(chng.From) && skip(chng.To) {
					continue
				}

				if skip(chng.From) && !skip(chng.To) {
					path = color.GreenString(path)
					chng.From = color.RedString(`""`)
					chng.To = color.GreenString(`"` + chng.To + `"`)
				} else if !skip(chng.From) && skip(chng.To) {
					path = color.RedString(path)
					chng.From = color.RedString(`"` + chng.From + `"`)
					chng.To = color.RedString(`""`)
				} else if chng.From != chng.To {
					path = color.YellowString(path)
					chng.From = color.RedString(`"` + chng.From + `"`)
					chng.To = color.GreenString(`"` + chng.To + `"`)
				}

				table.Append([]string{path, chng.From, chng.To})
			}

			table.Render()
		}
	}

	/*

	 */

	// /output(data)
}
