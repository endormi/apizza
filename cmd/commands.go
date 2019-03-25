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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/harrybrwn/apizza/cmd/internal/base"
	"github.com/harrybrwn/apizza/cmd/internal/obj"
	"github.com/harrybrwn/apizza/dawg"
	"github.com/harrybrwn/apizza/pkg/cache"
	"github.com/harrybrwn/apizza/pkg/config"
)

var db *cache.DataBase

// Execute runs the root command
func Execute() {
	var err error
	if err = config.SetConfig(".apizza", cfg); err != nil {
		handle(err, "Internal Error", 1)
	}

	builder := newBuilder()

	if db, err = cache.GetDB(builder.dbPath()); err != nil {
		handle(err, "Internal Error", 1)
	}

	defer func() {
		if err = db.Close(); err != nil {
			handle(err, "Internal Error", 1)
		}
		if err = config.Save(); err != nil {
			handle(err, "Internal Error", 1)
		}
	}()

	if _, err = builder.exec(); err != nil {
		handle(err, "Error", 0)
	}
}

func handle(e error, msg string, code int) { fmt.Fprintf(os.Stderr, "%s: %s\n", msg, e); os.Exit(code) }

type basecmd struct {
	*base.Command
	menu    *dawg.Menu
	tsDecay time.Duration
}

func (c *basecmd) getstore() (err error) {
	if store == nil {
		if store, err = dawg.NearestStore(c.Addr, cfg.Service); err != nil {
			return err
		}
	}
	return nil
}

func (c *basecmd) cacheNewMenu() (err error) {
	if err = c.getstore(); err != nil {
		return err
	}

	c.menu, err = store.Menu()
	if err != nil {
		return err
	}
	raw, err := json.Marshal(c.menu)
	if err != nil {
		return err
	}
	return db.Put("menu", raw)
}

func (c *basecmd) getCachedMenu() error {
	if c.menu == nil {
		raw, err := db.Get("menu")
		if err != nil {
			return err
		}
		c.menu = &dawg.Menu{}
		return json.Unmarshal(raw, c.menu)
	}
	return nil
}

var _ cache.Updater = (*basecmd)(nil)

func (c *basecmd) OnUpdate() error {
	return c.cacheNewMenu()
}

func (c *basecmd) NotUpdate() error {
	return c.getCachedMenu()
}

func (c *basecmd) Decay() time.Duration {
	return c.tsDecay
}

func newCommand(use, short string, c base.CliCommand) *basecmd {
	return &basecmd{
		Command: base.NewCommand(use, short, c.Run),
		tsDecay: 12 * time.Hour,
	}
}

type commandBuilder interface {
	exec()
}

type cliBuilder struct {
	root base.CliCommand
	addr *obj.Address
}

func newBuilder() *cliBuilder {
	b := &cliBuilder{root: newApizzaCmd()}

	addrStr := b.root.(*apizzaCmd).address
	if addrStr == "" {
		b.addr = &cfg.Address
	} else {
		b.addr = nil
	}
	return b
}

func (b *cliBuilder) exec() (*cobra.Command, error) {
	b.root.Addcmd(
		b.newCartCmd().Addcmd(
			b.newAddOrderCmd(),
		),
		newConfigCmd().Addcmd(
			newConfigSet(),
			newConfigGet(),
		),
		b.newMenuCmd(),
	)
	return b.root.Cmd().ExecuteC()
}

func (b *cliBuilder) dbPath() string {
	return filepath.Join(config.Folder(), "cache", "apizza.db")
}

func (b *cliBuilder) newCommand(use, short string, c base.CliCommand) *basecmd {
	base := newCommand(use, short, c)
	base.Addr = b.addr
	return base
}
