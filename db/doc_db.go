package db

import (
	"bytes"
	"time"

	bolt "github.com/boltdb/bolt"
)

var docsB string = "release_docs"

var docsTimeFormat string = "2006-01-02"

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
