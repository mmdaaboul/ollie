package repo

import (
	"fmt"
	"os"
	"strings"

	"ollie/db"

	"github.com/boltdb/bolt"
)

func init() {
	CheckAndCreateRecord()
}

func GetCurrentDirectoryAndFolderName() (string, string, error) {

	dir, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("could not get current directory: %v", err)
	}
	folder := dir[strings.LastIndex(dir, "/")+1:]
	return dir, folder, nil
}

func CheckAndCreateRecord() error {
	dir, folder, err := GetCurrentDirectoryAndFolderName()
	if err != nil {
		return err
	}

	db := db.DB()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("repos"))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}

		v := bucket.Get([]byte(folder))
		if v == nil {
			err = bucket.Put([]byte(folder), []byte(dir))
			if err != nil {
				return fmt.Errorf("could not put to bucket: %v", err)
			}
		}

		return nil
	})

	return err
}

func GetAllRecords() (map[string]string, error) {
	records := make(map[string]string)
	db := db.DB()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("repos"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			records[string(k)] = string(v)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not get all records: %v", err)
	}

	return records, nil
}
