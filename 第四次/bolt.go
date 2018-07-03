package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {

	println("start:")
	db, err := bolt.Open("blockchain.db", 0600, nil)

	fmt.Print("%v", db)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {

		//b, err := tx.Bucket([]byte("tests2"))
		b := tx.Bucket([]byte("blocks"))

		/*

			if err != nil {
				log.Panic(err)
			}
			if b != nil {
				b.Put([]byte("mykey"), []byte("1234"))
				if err != nil {
					log.Panic(err)
				}
			}else{
				log.Panic("no database!")
			}
		*/

		//get
		tip := b.Get([]byte("l"))

		println(string(tip))

		return nil
	})

	db.Close()

}
