package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"

	"apizza/dawg"
	"apizza/pkg/config"
)

func initDatabase() (err error) {
	dir := filepath.Join(config.Folder(), "cache")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0700)
	}
	path := filepath.Join(dir, "apizza.db")
	db, err = bolt.Open(path, 0700, nil)
	return err
}

func bucketExists(tx *bolt.Tx, name string) bool {
	b := tx.Bucket([]byte(name))
	if b == nil {
		return false
	}
	return true
}

func bucketHas(b *bolt.Bucket, name string) bool {
	raw := b.Get([]byte(name))
	if raw == nil {
		return false
	}
	return true
}

// TODO: this is terrible and needs to be split up
func menuManagment() (err error) {
	var menuIsCached = true
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
	return err
}

func dbTest() {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Menu"))
		// fmt.Println("bucket:", b)
		fmt.Println(bucketHas(b, "no menu"))

		return nil
	})
	if err != nil {
		fmt.Println("Caught database error:", err)
	}
}
