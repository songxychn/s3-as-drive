package utils

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"s3-as-drive/types"
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

func GetConfig() (types.Config, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return types.Config{}, err
	}
	configFilePath := filepath.Join(baseDir, "config.json")
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return types.Config{}, err
	}
	config := types.Config{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return types.Config{}, err
	}
	return config, nil
}
