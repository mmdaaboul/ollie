package db_test

import (
	db "ollie/db"
	"reflect"
	"testing"

	bolt "github.com/boltdb/bolt"
)

func TestGetStacks(t *testing.T) {
	// Open a temporary BoltDB for testing
	db := db.DB()
	defer db.Close()

	// Create a test bucket and add some stacks
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("stacks"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("stack1"), []byte{})
		if err != nil {
			return err
		}
		err = b.Put([]byte("stack2"), []byte{})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Call the GetStacks function
	stacks, err := db.GetStacks()
	if err != nil {
		t.Fatal(err)
	}

	// Check the returned stacks
	expectedStacks := []string{"stack1", "stack2"}
	if !reflect.DeepEqual(stacks, expectedStacks) {
		t.Errorf("GetStacks() returned %v, expected %v", stacks, expectedStacks)
	}
}
