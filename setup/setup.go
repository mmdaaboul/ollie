package setup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	DbPath string `json:"dbpath"`
}

var (
	config Config
	once   sync.Once
)

func LoadConfig() (Config, error) {
	var err error
	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			return
		}

		configDir := filepath.Join(home, ".config", "ollie")
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return
		}

		configPath := filepath.Join(configDir, "config")
		var file *os.File
		file, err = os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(&config)
		if err != nil {
			config.DbPath = filepath.Join(home, ".config", "ollie", "db.db")
			file.Seek(0, 0)
			json.NewEncoder(file).Encode(&config)
		}
	})

	return config, err
}
