package utils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
	"s3-as-drive/backend/types"
)

var db *gorm.DB //全局的db对象，我们执行数据库操作主要通过他实现。

func init() {
	var err error

	baseDir, err := GetBaseDir()
	if err != nil {
		panic(err)
	}
	dbFilePath := filepath.Join(baseDir, "data.db")
	db, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&types.File{}, &types.SyncDir{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
