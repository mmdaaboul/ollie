package db

import (
	"bytes"
	"fmt"
	config "ollie/setup"
	"sort"
	"time"

	bolt "github.com/boltdb/bolt"
	"github.com/charmbracelet/log"
)

var db *bolt.DB
var stacksB string = "stacks"
var docsB string = "release_docs"

var stacksTimeFormat string = "2017.09.07 17:06:06"
var docsTimeFormat string = "2006-01-02"

// Define a struct to hold both stack name and date
type stackEntry struct {
	key  string
	date time.Time
}

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

func GetStacks() ([]string, error) {
	var stacks []stackEntry // Define a custom struct with key and date
	err := Update(db, func(tx *bolt.Tx) error {
		b := GetBucket(tx, stacksB)
		if b == nil {
			var err error
			b, err = CreateBucket(tx, stacksB)
			if err != nil {
				return fmt.Errorf("bucket %s not found", stacksB) // Use fmt.Errorf for better error handling
			}
		}

		return b.ForEach(func(k, v []byte) error {
			date, err := time.Parse(stacksTimeFormat, string(v))
			if err != nil {
				return err
			}
			stacks = append(stacks, stackEntry{key: string(k), date: date})
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	// Sort the stacks by date (ascending order)
	sort.Slice(stacks, func(i, j int) bool {
		return stacks[i].date.Before(stacks[j].date)
	})

	// Extract just the stack names from the sorted slice
	var stackNames []string
	for _, entry := range stacks {
		stackNames = append(stackNames, entry.key)
	}

	return stackNames, nil
}

func AddStack(stack string) error {
	return Update(db, func(tx *bolt.Tx) error {
		b := GetBucket(tx, stacksB)
		if b == nil {
			var err error
			b, err = CreateBucket(tx, stacksB)
			if err != nil {
				return err
			}
		}

		value := []byte(time.Now().Format(stacksTimeFormat))
		return Put(b, []byte(stack), value)
	})
}

func UpdateStack(stack string) error {
	return Update(db, func(tx *bolt.Tx) error {
		b := GetBucket(tx, stacksB)
		if b == nil {
			return fmt.Errorf("bucket %s not found", stacksB)
		}

		value := []byte(time.Now().Format(stacksTimeFormat))
		return b.Put([]byte(stack), value)
	})
}

func GetReleaseDocs() ([]string, error) {
	today := time.Now().Format(docsTimeFormat)
	yesterday := time.Now().AddDate(0, 0, -1).Format(docsTimeFormat) // Subtract one day
	var docs []string

	err := Update(db, func(tx *bolt.Tx) error {
		b := GetBucket(tx, docsB)
		if b == nil {
			var err error
			b, err = CreateBucket(tx, docsB)
			if err != nil {
				return err
			}
		}

		return b.ForEach(func(k, v []byte) error {
			if bytes.Equal(v, []byte(today)) || bytes.Equal(v, []byte(yesterday)) {
				docs = append(docs, string(k))
			}
			return nil
		})
	})

	return docs, err
}

func AddReleaseDoc(key string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := GetBucket(tx, docsB)
		if b == nil {
			var err error
			b, err = CreateBucket(tx, docsB)
			if err != nil {
				return err
			}
		}
		value := []byte(time.Now().Format(docsTimeFormat))
		return b.Put([]byte(key), value)
	})
	return err
}
