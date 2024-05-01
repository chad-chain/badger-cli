package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
)

func main() {
	// get all the accounts from badger db

	// get path to the database from args
	path := os.Args[1]
	opts := badger.DefaultOptions(path)
	opts.Dir = path
	opts.ValueDir = path
	opts.ReadOnly = true
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// stop for 1 second to allow the db to be opened
	time.Sleep(1 * time.Second)
	fmt.Printf("Reading all the keys from the DB\n\n")

	err = db.View(func(txn *badger.Txn) error {
		// Iterate over the keys in the DB
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		// make a table of all the keys and values

		fmt.Println("Key\tValue")
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("%s\t%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
