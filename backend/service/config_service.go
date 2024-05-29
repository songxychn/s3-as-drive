package service

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"s3-as-drive/backend/types"
	"s3-as-drive/backend/utils"
)

type ConfigService struct {
	ctx context.Context
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (configService *ConfigService) Startup(ctx context.Context) {
	configService.ctx = ctx
}

func (configService *ConfigService) GetConfig() (types.Config, error) {
	config, err := utils.GetConfig()
	if err != nil {
		return types.Config{}, err
	}
	return config, nil
}

func (configService *ConfigService) UpdateConfig(config types.Config) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	baseDir, err := utils.GetBaseDir()
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(baseDir, "config.json")
	err = os.WriteFile(configFilePath, bytes, 0644)
	return nil
}
