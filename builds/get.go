/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats"
	"github.com/r3labs/graph"
	"github.com/spf13/cobra"
)

func Get(cmd *cobra.Command, args []string) {
	var gg map[string]interface{}

	nuri, _ := cmd.Flags().GetString("nats")
	buildid, _ := cmd.Flags().GetString("build")
	ctags, _ := cmd.Flags().GetString("component-tags")
	ctype, _ := cmd.Flags().GetString("component-type")

	nc, err := nats.Connect(nuri)
	if err != nil {
		panic(err)
	}

	msg, err := nc.Request("build.get.mapping", []byte(`{"id": "`+buildid+`"}`), time.Second)
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		output(msg.Data)
		os.Exit(0)
	}

	err = json.Unmarshal(msg.Data, &gg)
	if err != nil {
		panic(err)
	}

	g := graph.New()

	err = g.Load(gg)
	if err != nil {
		panic(err)
	}

	query := strings.Split(args[0], ".")

	var cg graph.ComponentGroup

	switch query[0] {
	case "components":
		cg = g.GetComponents()
	case "changes":
		cg = g.GetChanges()
	default:
		panic("no valid query specified")
	}

	if ctags != "" {
	}

	if ctype != "" {
		cg = cg.ByType(ctype)
	}

	if len(query) > 1 {
		for i := len(cg) - 1; i >= 0; i-- {
			name := strings.Split(cg[i].GetID(), "::")[1]
			if name != query[1] {
				cg = append(cg[:i], cg[i+1:]...)
			}
		}
	}

	for _, c := range cg {
		data, err := json.Marshal(c)
		if err != nil {
			panic(err)
		}

		output(data)
	}
}

func output(data []byte) {
	var out bytes.Buffer

	err := json.Indent(&out, data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out.Bytes()))
}
