package utils

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"s3-as-drive/backend/types"
)

func init() {
	// 创建 app 根目录
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := currentUser.HomeDir
	baseDirPath := filepath.Join(homeDir, ".s3-as-drive")
	err = os.MkdirAll(baseDirPath, 0755)
	if err != nil {
		panic(err)
	}

	configFilePath := filepath.Join(baseDirPath, "config.json")
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// 配置文件不存在, 需要创建
		config := types.Config{
			S3Config: types.S3Config{
				Endpoint:  "play.min.io",
				AccessKey: "Q3AM3UQ867SPQQA43P2F",
				SecretKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
			},
			DownloadConfig: types.DownloadConfig{
				Dir: filepath.Join(homeDir, "Downloads"),
			},
		}
		bytes, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(configFilePath, bytes, 0644)
		if err != nil {
			panic(err)
		}
	}
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
