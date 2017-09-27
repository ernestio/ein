/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats"
	"github.com/r3labs/graph"
	"github.com/spf13/cobra"
)

func Graphviz(cmd *cobra.Command, args []string) {
	var gg map[string]interface{}

	nuri, _ := cmd.Flags().GetString("nats")
	buildid, _ := cmd.Flags().GetString("build")

	nc, err := nats.Connect(nuri)
	if err != nil {
		panic(err)
	}

	msg, err := nc.Request("build.get.mapping", []byte(`{"id": "`+buildid+`"}`), time.Second)
	if err != nil {
		panic(err)
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

	fmt.Println(g.Graphviz())
}
