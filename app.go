package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
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

var db *sql.DB = nil
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
	db, err = sql.Open("sqlite3", dbFilePath)
	if err != nil {
		panic(err)
	}
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS file (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path VARCHAR NOT NULL UNIQUE,
		key VARCHAR NOT NULL,
		is_dir boolean NOT NULL,
		create_time DATETIME DEFAULT CURRENT_TIMESTAMP);
    `
	_, err = db.Exec(createTableSQL)
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

	// 创建配置文件
	configFilePath := filepath.Join(baseDirPath, "config.json")
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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetConfig() types.Result {
	config, err := utils.GetConfig()
	if err != nil {
		return types.Error(err.Error())
	}
	return types.Success(config)
}

func (a *App) UpdateConfig(config types.Config) types.Result {
	bytes, err := json.Marshal(config)
	if err != nil {
		log.Println(err.Error())
		return types.Error(err.Error())
	}
	baseDir, err := utils.GetBaseDir()
	if err != nil {
		return types.Error(err.Error())
	}
	configFilePath := filepath.Join(baseDir, "config.json")
	err = os.WriteFile(configFilePath, bytes, 0644)
	return types.SuccessEmpty()
}

type File struct {
	ID         string    `json:"id"`
	Path       string    `json:"path"`
	Key        string    `json:"key"`
	IsDir      bool      `json:"isDir"`
	CreateTime time.Time `json:"createTime"`
}

func (a *App) GetFileList() types.Result {
	rows, err := db.Query("select id, path, key, is_dir, create_time from file;")
	if err != nil {
		return types.Error(err.Error())
	}
	var files []File
	for rows.Next() {
		var file File
		err = rows.Scan(&file.ID, &file.Path, &file.Key, &file.IsDir, &file.CreateTime)
		if err != nil {
			return types.Error(err.Error())
		}
		files = append(files, file)
	}
	return types.Success(files)
}

func (a *App) SelectFiles(dir string) types.Result {
	config, err := utils.GetConfig()
	if err != nil {
		return types.Error(err.Error())
	}

	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return types.Error(err.Error())
	}

	insertSQL := `
		   INSERT INTO file (path, key, is_dir) VALUES (?, ?, false);
		`
	for _, filePathItem := range files {

		// 上传 s3
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return types.Error(err.Error())
		}
		key := fmt.Sprintf("%s%s", newUUID, filepath.Ext(filePathItem))
		_, err = s3Client.FPutObject(context.Background(), config.S3Config.Bucket, key, filePathItem, minio.PutObjectOptions{})
		if err != nil {
			return types.Error(err.Error())
		}

		// 插入数据库
		index := strings.LastIndex(filePathItem, "/")
		path := filepath.Join(dir, filePathItem[index+1:])
		_, err = db.Exec(insertSQL, path, key)
		if err != nil {
			return types.Error(err.Error())
		}
	}
	return types.Success(len(files))
}

func (a *App) DownloadFile(fileId string) types.Result {
	config, err := utils.GetConfig()
	if err != nil {
		return types.Error(err.Error())
	}
	var file File
	err = db.QueryRow("select id, path, key, is_dir, create_time from file where id = ?", fileId).Scan(&file.ID, &file.Path, &file.Key, &file.IsDir, &file.CreateTime)
	if err != nil {
		return types.Error(err.Error())
	}
	downloadPath := filepath.Join(config.DownloadConfig.Dir, filepath.Base(file.Path))
	err = s3Client.FGetObject(context.Background(), config.S3Config.Bucket, file.Key, downloadPath, minio.GetObjectOptions{})
	if err != nil {
		return types.Error(err.Error())
	}
	return types.SuccessEmpty()
}
