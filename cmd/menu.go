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
	"fmt"
	"strings"
	"unicode/utf8"
	"encoding/json"

	"apizza/dawg"
	"github.com/spf13/cobra"
	"github.com/boltdb/bolt"
)

var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Get the Dominos menu.",
	RunE:  func(cmd *cobra.Command, args []string) error {
		var (
			err error
			bflag = cmd.Flags().GetBool
			menuIsCached = true
		)

		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("Menu"))
			return err
		})
		if err != nil {
			return err
		}

		err = db.View(func(tx *bolt.Tx) error {
			var err error
			b := tx.Bucket([]byte("Menu"))
			rawmenu := b.Get([]byte("menu"))
			if rawmenu == nil {
				menuIsCached = false
				menu, err = store.Menu()
				return err
			}
			menu = &dawg.Menu{}
			return json.Unmarshal(rawmenu, menu)
		})
		if err != nil {
			return err
		}

		err = db.Update(func(tx *bolt.Tx) error {
			if menuIsCached {
				return nil
			}
			b := tx.Bucket([]byte("Menu"))
			data, err := json.Marshal(menu)
			if err != nil {
				return err
			}
			return b.Put([]byte("menu"), data)
		})
		if err != nil {
			return err
		}

		if toppings, err := bflag("toppings"); toppings && err == nil {
			printToppings()
		} else if food, err := bflag("food"); food && err == nil {
			printTable()
		}
		return nil
	},
}

func printTable() {
	var f func(map[string]interface{}, string)

	f = func(m map[string]interface{}, spacer string) {
		cats := m["Categories"].([]interface{})
		prods := m["Products"].([]interface{})

		// if there is nothing in that category, dont print the code name
		if len(cats) != 0 || len(prods) != 0 {
			fmt.Print(spacer, m["Code"], "\n")
		}
		if len(cats) > 0 { // the recursive part
			for _, c := range cats {
				f(c.(map[string]interface{}), spacer+"  ")
			}
		} else if len(prods) > 0 { // the printing part
			var prod map[string]interface{}
			max := maxStrLen(prods) + 2
			for _, p := range prods {
				_, ok := menu.Products[p.(string)]
				if ok {
					prod = menu.Products[p.(string)].(map[string]interface{})
				} else {
					continue
				}
				space := strings.Repeat(" ", max-strLen(p.(string)))
				fmt.Print(spacer+"  ", p, space, prod["Name"], "\n")
			}
			print("\n")
		}
	}
	f(menu.Categorization["Food"].(map[string]interface{}), "")
}

func printToppings() {
	indent := strings.Repeat(" ", 4)
	for key, val := range menu.Toppings {
		fmt.Print("  ", key, "\n")
		for k, v := range val.(map[string]interface{}) {
			spacer := strings.Repeat(" ", 3-strLen(k))
			fmt.Print(indent, k, spacer, v.(map[string]interface{})["Name"], "\n")
		}
		print("\n")
	}
}

func maxStrLen(list []interface{}) int {
	max := 0
	for _, s := range list {
		length := strLen(s.(string))
		if length > max {
			max = length
		}
	}

	return max
}

var strLen = utf8.RuneCountInString

func init() {
	var bflag = menuCmd.Flags().BoolP

	bflag("all", "a", false, "show the entire menu")
	bflag("food", "f", true, "print out the food items on the menu")
	bflag("toppings", "t", false, "print out the toppings on the menu")

	rootCmd.AddCommand(menuCmd)
}
