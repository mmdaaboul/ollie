package setup_test

import (
	"ollie/setup"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()

	// Set up the test config file
	configPath := filepath.Join(tempDir, "config")
	configData := []byte(`{"dbpath": "/path/to/db.db"}`)
	err := os.WriteFile(configPath, configData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Call the LoadConfig function
	config, err := setup.LoadConfig()

	// Verify that the config is loaded correctly
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if config.DbPath == "" {
		t.Errorf("Expected dbpath to be present, got an empty string")
	}
}
