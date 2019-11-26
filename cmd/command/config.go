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

package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/harrybrwn/apizza/cmd/internal/base"
	"github.com/harrybrwn/apizza/pkg/config"
)

// var cfg = &base.Config{}

type configCmd struct {
	base.CliCommand
	file   bool
	dir    bool
	getall bool
	edit   bool

	card string
	exp  string
}

func (c *configCmd) Run(cmd *cobra.Command, args []string) error {
	if c.edit {
		return config.Edit()
	}
	if c.file {
		c.Println(config.File())
		return nil
	}
	if c.dir {
		c.Println(config.Folder())
		return nil
	}
	if c.getall {
		return config.FprintAll(cmd.OutOrStdout(), config.Object())
	}
	return cmd.Usage()
}

// NewConfigCmd creates a new config command.
func NewConfigCmd(b base.Builder) base.CliCommand {
	c := &configCmd{file: false, dir: false}
	c.CliCommand = b.Build("config", "Configure apizza", c)
	c.SetOutput(b.Output())
	c.Cmd().Long = `The 'config' command is used for accessing the .apizza config file
in your home directory. Feel free to edit the .apizza json file
by hand or use the 'config' command.

ex. 'apizza config get name' or 'apizza config set name=<your name>'`

	c.Flags().BoolVarP(&c.file, "file", "f", c.file, "show the path to the config.json file")
	c.Flags().BoolVarP(&c.dir, "dir", "d", c.dir, "show the apizza config directory path")
	c.Flags().BoolVar(&c.getall, "get-all", c.getall, "show all the contents of the config file")
	c.Flags().BoolVarP(&c.edit, "edit", "e", false, "open the conifg file with the text editor set by $EDITOR")

	c.Flags().StringVar(&c.card, "card", "", "store encrypted credit card number in the database")
	c.Flags().StringVar(&c.exp, "expiration", "", "store the encrypted expiration data of your credit card")

	c.Addcmd(newConfigSet(), newConfigGet())
	return c
}

func newConfigSet() base.CliCommand {
	return base.NewCommand(
		"set", "change variables in the config file",
		func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("no variable given")
			}

			for _, arg := range args {
				keys := strings.Split(arg, "=")
				if len(keys) < 2 || keys[0] == "" || keys[1] == "" {
					return errors.New(`use '<key>=<value>' format (no spaces), use <key>='-' to set as empty`)
				}

				if keys[1] == "-" {
					keys[1] = ""
				}
				err := config.Set(keys[0], keys[1])
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func newConfigGet() base.CliCommand {
	return base.NewCommand(
		"get", "print the specified config variable to screen",
		func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("no variable given")
			}

			for _, arg := range args {
				v := config.Get(arg)
				if v == nil {
					return fmt.Errorf("cannot find %s", arg)
				}
				cmd.Println(v)
			}
			return nil
		},
	)
}