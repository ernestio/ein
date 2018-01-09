/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ernestio/mapping"
	"github.com/nats-io/nats"
	"github.com/olekukonko/tablewriter"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
	"github.com/spf13/cobra"
)

/*
var groups = map[string][]string{
	"s3":        nil,
	"instances": nil,
}
*/

type UniqueCollection []string

func (uc *UniqueCollection) add(item string) {
	for i := 0; i < len((*uc)); i++ {
		if (*uc)[i] == item {
			return
		}
	}

	(*uc) = append((*uc), item)
}

func componentsByType(t string, cl diff.Changelog) []string {
	var uc UniqueCollection
	for _, c := range cl.Filter([]string{t + "::*"}) {
		name := strings.TrimPrefix(c.Path[0], t+"::")
		uc.add(name)
	}

	return uc
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

	// cl := g.Changelog.Filter([]string{"s3::*"})

	data, err := json.Marshal(g.Changelog)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, name := range componentsByType("s3", g.Changelog) {
		fmt.Println(name)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Field", "From", "To"})

		for _, c := range g.Changelog.Filter([]string{"s3::" + name}) {
			table.Append([]string{strings.Join(c.Path, "."), fmt.Sprint(c.From), fmt.Sprint(c.To)})
		}
		table.Render()
	}

	output(data)
}
