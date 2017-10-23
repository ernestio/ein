/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ernestio/mapping"
	"github.com/ernestio/mapping/definition"
	"github.com/ghodss/yaml"
	"github.com/nats-io/nats"
	"github.com/r3labs/graph"
	"github.com/spf13/cobra"
)

func Map(cmd *cobra.Command, args []string) {
	var def definition.Definition

	nuri, _ := cmd.Flags().GetString("nats")
	graphviz, _ := cmd.Flags().GetBool("graphviz")

	nc, err := nats.Connect(nuri)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &def)
	if err != nil {
		panic(err)
	}

	m := mapping.New(nc, def.FullName())

	err = m.Apply(&def)
	if err != nil {
		panic(err)
	}

	if graphviz {
		g := graph.New()
		err = g.Load(m.Result)
		if err != nil {
			panic(err)
		}

		fmt.Println(g.Graphviz())
		return
	}

	data, err = json.Marshal(m.Result)
	if err != nil {
		panic(err)
	}

	output(data)
}
