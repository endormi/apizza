// Copyright © 2019 Harrison Brown harrybrown98@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/harrybrwn/apizza/cmd/internal/base"
	"github.com/harrybrwn/apizza/dawg"
)

var (
	// addr  *dawg.Address
	store *dawg.Store
)

type apizzaCmd struct {
	*basecmd
	address    string
	service    string
	storeID    string
	test       bool
	clearCache bool
}

func (c *apizzaCmd) Run(cmd *cobra.Command, args []string) (err error) {
	if test {
		all, err := db.GetAll()
		if err != nil {
			return err
		}
		for k := range all {
			c.Printf("%v\n", k)
		}
		return nil
	}
	if c.clearCache {
		if err := db.Close(); err != nil {
			return err
		}
		c.Printf("removing %s\n", db.Path)
		return os.Remove(db.Path)
	}
	return cmd.Usage()
}

var test bool

func newApizzaCmd() base.CliCommand {
	c := &apizzaCmd{address: "", service: cfg.Service, clearCache: false}
	c.basecmd = newCommand("apizza", "Dominos pizza from the command line.", c)

	// c.cmd.PersistentFlags().StringVar(&c.address, "address", c.address, "use a specific address")
	c.Cmd().PersistentFlags().StringVar(&c.service, "service", c.service, "select a Dominos service, either 'Delivery' or 'Carryout'")

	c.Cmd().PersistentFlags().BoolVar(&test, "test", false, "testing flag (for development)")
	c.Cmd().PersistentFlags().MarkHidden("test")

	c.Flags().BoolVar(&c.clearCache, "clear-cache", c.clearCache, "delete the database used for caching")
	return c
}
