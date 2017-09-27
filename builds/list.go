/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats"
	"github.com/spf13/cobra"
)

type Env struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Build struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func List(cmd *cobra.Command, args []string) {
	var env Env
	var builds []Build

	nuri, _ := cmd.Flags().GetString("nats")
	envName, _ := cmd.Flags().GetString("env")

	nc, err := nats.Connect(nuri)
	if err != nil {
		panic(err)
	}

	q := []byte(`{"name": "` + envName + `"}`)
	if envName == "" {
		fmt.Println("you must specify an environment")
		os.Exit(1)
	}

	msg, err := nc.Request("environment.get", q, time.Second*10)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(msg.Data, &env)
	if err != nil {
		panic(err)
	}

	q = []byte(`{"environment_id": "` + strconv.Itoa(env.ID) + `"}`)
	msg, err = nc.Request("build.find", q, time.Second*10)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(msg.Data, &builds)
	if err != nil {
		panic(err)
	}

	for _, b := range builds {
		fmt.Printf("%s  %s  %s\n", b.ID, b.Status, b.CreatedAt)
	}
}
