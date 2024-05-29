package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"s3-as-drive/types"
	"s3-as-drive/utils"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

var db *gorm.DB = nil
var s3Client *minio.Client = nil

// NewApp creates a new App application struct
func NewApp() *App {
	initConfig()

	initDB()

	initS3Client()

	return &App{}
}

// 初始化数据库
func initDB() {
	baseDir, err := utils.GetBaseDir()
	if err != nil {
		panic(err)
	}
	dbFilePath := filepath.Join(baseDir, "data.db")
	db, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&types.File{})
	if err != nil {
		panic(err)
	}
}

// 初始化 s3 客户端
func initS3Client() {
	config, err := utils.GetConfig()
	if err != nil {
		panic(err)
	}
	s3Client, err = minio.New(config.S3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3Config.AccessKey, config.S3Config.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		panic(err)
	}
}

func initConfig() {
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
			return
		}
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetConfig() (types.Config, error) {
	config, err := utils.GetConfig()
	if err != nil {
		return types.Config{}, err
	}
	return config, nil
}

func (a *App) UpdateConfig(config types.Config) error {
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

func (a *App) GetFileList(currentDir string) []types.File {
	depth := strings.Count(currentDir, "/")
	var files []types.File
	db.Where("path like concat(?, '%') and depth = ?", currentDir, depth).Find(&files)
	return files
}

func (a *App) UploadFiles(currentDir string) (int, error) {
	config, err := utils.GetConfig()
	if err != nil {
		return 0, err
	}

	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
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

func (a *App) UploadDir(currentDir string) (string, error) {
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
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

func (a *App) DownloadFile(fileId uint) error {
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

func (a *App) Mkdir(parentDir string, newDir string) error {
	depth := strings.Count(parentDir, "/")
	db.Create(&types.File{
		Path:  fmt.Sprintf("%s%s", parentDir, newDir),
		IsDir: true,
		Depth: uint(depth),
	})
	return nil
}

// TODO 如果是目录怎么办
func (a *App) GetShareUrl(fileId uint, expireInSecond int) error {
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
	err = runtime.ClipboardSetText(a.ctx, presignedURL.String())
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteFile(fileId uint) error {
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
