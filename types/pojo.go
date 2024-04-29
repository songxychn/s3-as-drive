package types

import (
	"time"
)

type File struct {
	ID         string    `json:"id"`
	Path       string    `json:"path"`
	Key        *string   `json:"key"`
	CreateTime time.Time `json:"createTime"`
	IsDir      bool      `json:"isDir"`
	Depth      int       `json:"depth"`
	Size       *int64    `json:"size"`
}
