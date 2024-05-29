package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"s3-as-drive/backend/types"
	"s3-as-drive/backend/utils"
	"strings"
	"time"
)

type FileService struct {
	ctx context.Context
}

func NewFileService() *FileService {
	return &FileService{}
}

func (fileService *FileService) Startup(ctx context.Context) {
	fileService.ctx = ctx
}

func (fileService *FileService) GetFileList(currentDir string) []types.File {
	depth := strings.Count(currentDir, "/")
	var files []types.File
	db := utils.GetDB()
	db.Where("path like concat(?, '%') and depth = ?", currentDir, depth).Find(&files)
	return files
}

func (fileService *FileService) UploadFiles(currentDir string) (int, error) {
	db := utils.GetDB()
	s3Client := utils.GetS3Client()
	config, err := utils.GetConfig()
	if err != nil {
		return 0, err
	}

	files, err := runtime.OpenMultipleFilesDialog(fileService.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return 0, err
	}

	depth := strings.Count(currentDir, "/")

	for _, filePath := range files {
		// 上传 s3
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return 0, err
		}
		key := fmt.Sprintf("%s%s", newUUID, filepath.Ext(filePath))
		_, err = s3Client.FPutObject(context.Background(), config.S3Config.Bucket, key, filePath, minio.PutObjectOptions{})
		if err != nil {
			return 0, err
		}

		// 插入数据库
		path := filepath.Join(currentDir, filepath.Base(filePath))
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return 0, err
		}
		size := uint64(fileInfo.Size())
		db.Create(&types.File{
			Path:  path,
			Key:   &key,
			IsDir: false,
			Depth: uint(depth),
			Size:  &size,
		})
	}
	return len(files), nil
}

func (fileService *FileService) UploadDir(currentDir string) (string, error) {
	db := utils.GetDB()
	s3Client := utils.GetS3Client()
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	dirPath, err := runtime.OpenDirectoryDialog(fileService.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	if dirPath == "" {
		// 用户取消了
		return "", nil
	}

	base := filepath.Base(dirPath)

	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		newPath := fmt.Sprintf("%s%s", currentDir, path[strings.Index(path, base):])
		depth := strings.Count(newPath, "/")
		if d.IsDir() {
			db.Create(&types.File{
				Path:  newPath,
				IsDir: true,
				Depth: uint(depth),
			})
		} else {
			if filepath.Base(path) == ".DS_Store" {
				// macos 可忽略的文件
				return nil
			}
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			newUUID, err := uuid.NewUUID()
			if err != nil {
				return err
			}
			key := fmt.Sprintf("%s%s", newUUID, filepath.Ext(path))
			size := uint64(fileInfo.Size())
			db.Create(&types.File{
				Path:  newPath,
				Key:   &key,
				IsDir: false,
				Depth: uint(depth),
				Size:  &size,
			})

			_, err = s3Client.FPutObject(context.Background(), config.S3Config.Bucket, key, path, minio.PutObjectOptions{})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return dirPath, nil
}

func (fileService *FileService) DownloadFile(fileId uint) error {
	db := utils.GetDB()
	s3Client := utils.GetS3Client()
	config, err := utils.GetConfig()
	if err != nil {
		return err
	}
	var file types.File
	db.First(&file, fileId)

	if !file.IsDir {
		// 是单个文件
		downloadPath := filepath.Join(config.DownloadConfig.Dir, filepath.Base(file.Path))
		err = s3Client.FGetObject(context.Background(), config.S3Config.Bucket, *file.Key, downloadPath, minio.GetObjectOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	// 是目录
	var files []types.File
	db.Where("path like concat(?, '/%')", file.Path).Find(&files)
	for _, fileItem := range files {
		if fileItem.IsDir {
			continue
		}
		downloadPath := filepath.Join(config.DownloadConfig.Dir, fileItem.Path)
		err := s3Client.FGetObject(context.Background(), config.S3Config.Bucket, *fileItem.Key, downloadPath, minio.GetObjectOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (fileService *FileService) Mkdir(parentDir string, newDir string) error {
	db := utils.GetDB()
	depth := strings.Count(parentDir, "/")
	db.Create(&types.File{
		Path:  fmt.Sprintf("%s%s", parentDir, newDir),
		IsDir: true,
		Depth: uint(depth),
	})
	return nil
}

// TODO 如果是目录怎么办
func (fileService *FileService) GetShareUrl(fileId uint, expireInSecond int) error {
	db := utils.GetDB()
	s3Client := utils.GetS3Client()
	var file types.File
	db.First(&file, fileId)
	config, err := utils.GetConfig()
	if err != nil {
		return err
	}
	fileName := filepath.Base(file.Path)
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	presignedURL, err := s3Client.PresignedGetObject(context.Background(), config.S3Config.Bucket, *file.Key, time.Second*time.Duration(expireInSecond), reqParams)
	if err != nil {
		return err
	}

	// 把分享链接放进剪切板里
	err = runtime.ClipboardSetText(fileService.ctx, presignedURL.String())
	if err != nil {
		return err
	}
	return nil
}

func (fileService *FileService) DeleteFile(fileId uint) error {
	db := utils.GetDB()
	s3Client := utils.GetS3Client()
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	var file types.File
	db.First(&file, fileId)

	if !file.IsDir {
		// 单个文件
		db.Delete(&types.File{}, fileId)

		err = s3Client.RemoveObject(context.Background(), config.S3Config.Bucket, *file.Key, minio.RemoveObjectOptions{})
		if err != nil {
			log.Fatalln(err)
			return err
		}

		return nil
	}

	// 目录
	// 获取所有子文件
	var files []types.File
	db.Where("path like concat(?, '/%')", file.Path).Find(&files)
	// 删除数据库中的子文件和本目录
	db.Where("path like concat(?, '/%') or id = ?", file.Path, fileId).Delete(&types.File{})

	// 删除 s3 中的子文件
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, fileItem := range files {
			if !fileItem.IsDir {
				objectsCh <- minio.ObjectInfo{
					Key: *fileItem.Key,
				}
			}
		}
	}()

	for deleteErr := range s3Client.RemoveObjects(context.Background(), config.S3Config.Bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		log.Println(deleteErr)
	}
	return nil
}
