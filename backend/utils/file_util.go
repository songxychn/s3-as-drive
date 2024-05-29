package utils

import (
	"os/user"
	"path/filepath"
)

func GetBaseDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	homeDir := currentUser.HomeDir
	dirPath := filepath.Join(homeDir, ".s3-as-drive")
	return dirPath, nil
}
