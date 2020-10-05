package database

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var (
	dbName     string = "url.db"
	bucketName []byte = []byte("URL")
)

type boltDB struct {
	*bolt.DB
}

func createBoltDatabase() (Database, error) {
	// create a new client with default options
	var db boltDB
	var err error

	db.DB, err = bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return nil
	})

	return &db, err
}

// Set: Implementation of the database interface
func (db *boltDB) Set(key string, value interface{}) error {
	var valueB []byte
	switch tmp := value.(type) {
	case []byte:
		valueB = tmp
	case string:
		valueB = []byte(tmp)
	default:
		return fmt.Errorf("Value must be of type []byte or string")
	}

	log.Println("(BoltDB): saving data:", value, " for key:", key)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Put([]byte(key), valueB)
		return err
	})

	return err
}

// Get: Implementation of the database interface
func (db *boltDB) Get(key string) (interface{}, error) {
	var v []byte

	log.Println("(BoltDB): retrieving value from key:", key)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		v = b.Get([]byte(key))
		return nil
	})

	return string(v), err
}

// Delete: Implementation of the database interface
func (db *boltDB) Delete(key string) error {
	log.Println("(BoltDB): deleting value from key:", key)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Delete([]byte(key))
		return err
	})

	return err
}
