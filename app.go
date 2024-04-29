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
		CREATE TABLE IF NOT EXISTS "file" (
		  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
		  "path" VARCHAR NOT NULL,
		  "key" VARCHAR,
		  "create_time" DATETIME DEFAULT CURRENT_TIMESTAMP,
		  "is_dir" boolean NOT NULL,
		  "depth" integer NOT NULL,
		  "size" integer,
		  UNIQUE ("path" ASC)
		);
		
		CREATE INDEX IF NOT EXISTS "idx_depth"
		ON "file" (
		  "depth" ASC
		);
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
	depth := strings.Count(currentDir, "/")
	rows, err := db.Query(`select id, path, key, create_time, is_dir, size from file 
		where path like concat(?, '%') and depth = ?`,
		currentDir, depth)
	if err != nil {
		return types.Error(err.Error())
	}
	var files []types.File
	for rows.Next() {
		var file types.File
		err = rows.Scan(&file.ID, &file.Path, &file.Key, &file.CreateTime, &file.IsDir, &file.Size)
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

	depth := strings.Count(currentDir, "/")
	insertSQL := "INSERT INTO file (path, key, is_dir, depth, size) VALUES (?, ?, false, ?, ?);"
	for _, filePath := range files {
		// 上传 s3
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return types.Error(err.Error())
		}
		key := fmt.Sprintf("%s%s", newUUID, filepath.Ext(filePath))
		_, err = s3Client.FPutObject(context.Background(), config.S3Config.Bucket, key, filePath, minio.PutObjectOptions{})
		if err != nil {
			return types.Error(err.Error())
		}

		// 插入数据库
		path := filepath.Join(currentDir, filepath.Base(filePath))
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return types.Error(err.Error())
		}
		_, err = db.Exec(insertSQL, path, key, depth, fileInfo.Size())
		if err != nil {
			return types.Error(err.Error())
		}
	}
	return types.Success(len(files))
}

func (a *App) UploadDir(currentDir string) types.Result {
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalln(err)
		return types.ErrorEmpty()
	}
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		log.Fatalln(err)
		return types.ErrorEmpty()
	}
	if dirPath == "" {
		// 用户取消了
		return types.Success("")
	}

	base := filepath.Base(dirPath)

	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		newPath := fmt.Sprintf("%s%s", currentDir, path[strings.Index(path, base):])
		depth := strings.Count(newPath, "/")
		if d.IsDir() {
			_, err = db.Exec("insert into file(path, is_dir, depth) values (?, true, ?)", newPath, depth)
			if err != nil {
				return err
			}
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
			_, err = db.Exec("insert into file(path, key, is_dir, depth, size) values (?, ?, false, ?, ?)", newPath, key, depth, fileInfo.Size())
			if err != nil {
				return err
			}

			_, err = s3Client.FPutObject(context.Background(), config.S3Config.Bucket, key, path, minio.PutObjectOptions{})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
		return types.ErrorEmpty()
	}
	return types.Success(dirPath)
}

func (a *App) DownloadFile(fileId string) types.Result {
	config, err := utils.GetConfig()
	if err != nil {
		return types.Error(err.Error())
	}
	var file types.File
	err = db.QueryRow("select id, path, key, is_dir from file where id = ?", fileId).Scan(&file.ID, &file.Path, &file.Key, &file.IsDir)
	if err != nil {
		return types.Error(err.Error())
	}

	if !file.IsDir {
		// 是单个文件
		downloadPath := filepath.Join(config.DownloadConfig.Dir, filepath.Base(file.Path))
		err = s3Client.FGetObject(context.Background(), config.S3Config.Bucket, *file.Key, downloadPath, minio.GetObjectOptions{})
		if err != nil {
			return types.Error(err.Error())
		}
		return types.SuccessEmpty()
	}

	// 是目录
	rows, err := db.Query("select id, path, key, is_dir from file where path like concat(?, '/%')", file.Path)
	if err != nil {
		return types.Error(err.Error())
	}
	for rows.Next() {
		var file types.File
		err = rows.Scan(&file.ID, &file.Path, &file.Key, &file.IsDir)
		if err != nil {
			return types.Error(err.Error())
		}

		if file.IsDir {
			continue
		}
		downloadPath := filepath.Join(config.DownloadConfig.Dir, file.Path)
		err := s3Client.FGetObject(context.Background(), config.S3Config.Bucket, *file.Key, downloadPath, minio.GetObjectOptions{})
		if err != nil {
			return types.Error(err.Error())
		}
	}
	return types.SuccessEmpty()
}

func (a *App) Mkdir(parentDir string, newDir string) types.Result {
	depth := strings.Count(parentDir, "/")
	_, err := db.Exec("insert into file(path, is_dir, depth) values (?, true, ?)", fmt.Sprintf("%s%s", parentDir, newDir), depth)
	if err != nil {
		return types.Error(err.Error())
	}
	return types.SuccessEmpty()
}

func (a *App) GetShareUrl(fileId string, expireInSecond int) types.Result {
	row := db.QueryRow("select id, path, key, is_dir from file where id = ?", fileId)
	var file types.File
	err := row.Scan(&file.ID, &file.Path, &file.Key, &file.IsDir)
	if err != nil {
		return types.Error(err.Error())
	}
	// TODO 如果是目录怎么办

	config, err := utils.GetConfig()
	if err != nil {
		return types.Error(err.Error())
	}
	fileName := filepath.Base(file.Path)
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	presignedURL, err := s3Client.PresignedGetObject(context.Background(), config.S3Config.Bucket, *file.Key, time.Second*time.Duration(expireInSecond), reqParams)
	if err != nil {
		return types.Error(err.Error())
	}

	// 把分享链接放进剪切板里
	err = runtime.ClipboardSetText(a.ctx, presignedURL.String())
	if err != nil {
		return types.Error(err.Error())
	}
	return types.SuccessEmpty()
}

func (a *App) DeleteFile(fileId string) types.Result {
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalln(err)
		return types.ErrorEmpty()
	}

	var file types.File
	row := db.QueryRow("select id, path, key, is_dir, depth from file where id = ?", fileId)
	err = row.Scan(&file.ID, &file.Path, &file.Key, &file.IsDir, &file.Depth)
	if err != nil {
		log.Fatalln(err)
		return types.ErrorEmpty()
	}

	if !file.IsDir {
		// 单个文件
		_, err = db.Exec("delete from file where id = ?", fileId)
		if err != nil {
			log.Fatalln(err)
			return types.ErrorEmpty()
		}

		err = s3Client.RemoveObject(context.Background(), config.S3Config.Bucket, *file.Key, minio.RemoveObjectOptions{})
		if err != nil {
			log.Fatalln(err)
			return types.ErrorEmpty()
		}

		return types.SuccessEmpty()
	}

	// 目录
	// 获取所有子文件
	rows, err := db.Query("select id, path, key, is_dir from file where path like concat(?, '/%')", file.Path)
	if err != nil {
		log.Fatalln(err)
		return types.ErrorEmpty()
	}

	// 删除数据库中的子文件和本目录
	_, err = db.Exec("delete from file where path like concat(?, '/%') or id = ?", file.Path, fileId)
	if err != nil {
		log.Fatalln(err)
		return types.ErrorEmpty()
	}

	// 删除 s3 中的子文件
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)

		for rows.Next() {
			var file types.File
			err = rows.Scan(&file.ID, &file.Path, &file.Key, &file.IsDir)
			if err != nil {
				log.Println(err)
				continue
			}
			if !file.IsDir {
				objectsCh <- minio.ObjectInfo{
					Key: *file.Key,
				}
			}
		}
	}()

	for deleteErr := range s3Client.RemoveObjects(context.Background(), config.S3Config.Bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		log.Println(deleteErr)
	}
	return types.SuccessEmpty()
}
