package db

import (
	"fmt"
	"sort"
	"time"

	bolt "github.com/boltdb/bolt"
)

var stacksB string = "stacks"

var stacksTimeFormat string = "2017.09.07"

// Define a struct to hold both stack name and date
type stackEntry struct {
	key  string
	date time.Time
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

func HasStack(stack string) (bool, error) {
	var value []byte
	err := View(db, func(tx *bolt.Tx) error {
		b := GetBucket(tx, stacksB)
		if b == nil {
			return fmt.Errorf("%s bucket doesn't exist", stacksB)
		}

		value = b.Get([]byte(stack))
		return nil
	})

	if err != nil {
		return false, err
	}
	return value != nil, nil
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

func DeleteStack(stack string) error {
	return Update(db, func(tx *bolt.Tx) error {
		b := GetBucket(tx, stacksB)
		if b == nil {
			return fmt.Errorf("Bucket %s not found", stacksB)
		}

		return b.Delete([]byte(stack))
	})
}
