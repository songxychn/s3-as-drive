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
		key VARCHAR,
		is_dir boolean NOT NULL DEFAULT false,
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

func (a *App) GetFileList(currentDir string) types.Result {
	slashCount := strings.Count(currentDir, "/")
	rows, err := db.Query(`select id, path, key, create_time, is_dir from file 
		where path like concat(?, '%') 
		and (length(path) - length(replace(path,'/',''))) = ?`,
		currentDir, slashCount)
	if err != nil {
		return types.Error(err.Error())
	}
	var files []types.File
	for rows.Next() {
		var file types.File
		err = rows.Scan(&file.ID, &file.Path, &file.Key, &file.CreateTime, &file.IsDir)
		if err != nil {
			return types.Error(err.Error())
		}
		files = append(files, file)
	}
	return types.Success(files)
}

func (a *App) UploadFiles(currentDir string) types.Result {
	config, err := utils.GetConfig()
	if err != nil {
		return types.Error(err.Error())
	}

	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return types.Error(err.Error())
	}

	insertSQL := "INSERT INTO file (path, key) VALUES (?, ?);"
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
		path := filepath.Join(currentDir, filepath.Base(filePathItem))
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
	var file types.File
	err = db.QueryRow("select id, path, key, create_time, is_dir from file where id = ?", fileId).Scan(&file.ID, &file.Path, &file.Key, &file.CreateTime, &file.IsDir)
	if err != nil {
		return types.Error(err.Error())
	}

	if !file.IsDir {
		// 是单个文件
		downloadPath := filepath.Join(config.DownloadConfig.Dir, filepath.Base(file.Path))
		err = s3Client.FGetObject(context.Background(), config.S3Config.Bucket, file.Key.String, downloadPath, minio.GetObjectOptions{})
		if err != nil {
			return types.Error(err.Error())
		}
		return types.SuccessEmpty()
	}

	// 是目录
	rows, err := db.Query("select id, path, key, create_time, is_dir from file where path like concat(?, '%')", file.Path)
	if err != nil {
		return types.Error(err.Error())
	}
	for rows.Next() {
		var file types.File
		err = rows.Scan(&file.ID, &file.Path, &file.Key, &file.CreateTime, &file.IsDir)
		if err != nil {
			return types.Error(err.Error())
		}
		fmt.Println(file)

		if file.IsDir {
			continue
		}
		downloadPath := filepath.Join(config.DownloadConfig.Dir, file.Path)
		err := s3Client.FGetObject(context.Background(), config.S3Config.Bucket, file.Key.String, downloadPath, minio.GetObjectOptions{})
		if err != nil {
			return types.Error(err.Error())
		}
	}
	return types.SuccessEmpty()
}

func (a *App) Mkdir(parentDir string, newDir string) types.Result {
	_, err := db.Exec("insert into file(path, is_dir) values (?, true)", fmt.Sprintf("%s%s", parentDir, newDir))
	if err != nil {
		return types.Error(err.Error())
	}
	return types.SuccessEmpty()
}
