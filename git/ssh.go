package git

import (
	config "ollie/setup"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func publicKey() (*ssh.PublicKeys, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	filePath := cfg.SshKeyPath
	if filePath == "" {
		form := huh.NewInput().Title("Enter path to SSH Key for Git").Value(&filePath)
		form.Run()
		newCfg := cfg
		newCfg.SshKeyPath = filePath
		config.UpdateConfig(newCfg)
	}

	var publicKey *ssh.PublicKeys
	sshKey, _ := os.ReadFile(filePath)
	publicKey, err = ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		return nil, err
	}
	return publicKey, err
}
