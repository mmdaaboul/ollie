package db

import (
	config "ollie/setup"

	bolt "github.com/boltdb/bolt"
	"github.com/charmbracelet/log"
)

var db *bolt.DB

func init() {
	var err error
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("There was an issue loading the config: %s", err)
	}

	db, err = bolt.Open(cfg.DbPath, 0666, nil)
	if err != nil {
		log.Fatalf("There was an issue opening up the database: %s", err)
	}
}

func DB() *bolt.DB {
	return db
}

func Close(db *bolt.DB) error {
	return db.Close()
}

func Update(db *bolt.DB, fn func(*bolt.Tx) error) error {
	return db.Update(fn)
}

func View(db *bolt.DB, fn func(*bolt.Tx) error) error {
	return db.View(fn)
}

func CreateBucket(tx *bolt.Tx, name string) (*bolt.Bucket, error) {
	return tx.CreateBucketIfNotExists([]byte(name))
}

func GetBucket(tx *bolt.Tx, name string) *bolt.Bucket {
	return tx.Bucket([]byte(name))
}

func Put(bucket *bolt.Bucket, key, value []byte) error {
	return bucket.Put(key, value)
}

func Get(bucket *bolt.Bucket, key []byte) []byte {
	return bucket.Get(key)
}

func Delete(bucket *bolt.Bucket, key []byte) error {
	return bucket.Delete(key)
}
