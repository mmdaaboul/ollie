package deploy_test

import (
	"ollie/db"
	"ollie/deploy"
	"ollie/git"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/charmbracelet/huh"
)

func TestDeployStack(t *testing.T) {
	// Set up a temporary test database
	db := db.DB()
	defer db.Close()

	// Add a stack to the database
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("stacks"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("stack1"), []byte{})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Mock user input for selecting an existing stack
	selectedStack := "stack1"
	huh.SetMockResponses([]string{selectedStack})

	// Call the DeployStack function
	deploy.DeployStack()

	// Verify that the selected stack was added to the database
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("stacks"))
		v := b.Get([]byte(selectedStack))
		if v == nil {
			t.Errorf("Expected stack %s to be added to the database", selectedStack)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reset mock responses
	huh.ResetMockResponses()

	// Mock user input for creating a new stack
	newStack := "new-stack"
	huh.SetMockResponses([]string{"new", newStack})

	// Call the DeployStack function again
	deploy.DeployStack()

	// Verify that the new stack was added to the database
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("stacks"))
		v := b.Get([]byte(newStack))
		if v == nil {
			t.Errorf("Expected stack %s to be added to the database", newStack)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reset mock responses
	huh.ResetMockResponses()

	// Mock user input for version bump
	bump := "patch"
	huh.SetMockResponses([]string{bump})

	// Call the DeployStack function again
	DeployStack()

	// Verify that the version bump was applied correctly
	version := git.GetVersion()
	expectedVersion := "1.0.1" // Assuming initial version is "1.0.0"
	if version != expectedVersion {
		t.Errorf("Expected version %s, got %s", expectedVersion, version)
	}
}
