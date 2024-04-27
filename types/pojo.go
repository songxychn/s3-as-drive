package types

import (
	"database/sql"
	"time"
)

type File struct {
	ID         string         `json:"id"`
	Path       string         `json:"path"`
	Key        sql.NullString `json:"key"`
	CreateTime time.Time      `json:"createTime"`
	IsDir      bool           `json:"isDir"`
}
