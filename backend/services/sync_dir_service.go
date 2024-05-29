package services

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"s3-as-drive/backend/types"
	"s3-as-drive/backend/utils"
)

type SyncDirService struct {
	ctx context.Context
}

func NewSyncDirService() *SyncDirService {
	return &SyncDirService{}
}

func (syncDirService *SyncDirService) Startup(ctx context.Context) {
	syncDirService.ctx = ctx
}

func (syncDirService *SyncDirService) CreateSyncDir(name string, path string) (types.SyncDir, error) {
	db := utils.GetDB()
	syncDir := types.SyncDir{
		Name: name,
		Path: path,
	}
	db.Create(&syncDir)
	return syncDir, nil
}

func (syncDirService *SyncDirService) ChoseDir() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(syncDirService.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return "", err
	}
	return dir, nil
}
