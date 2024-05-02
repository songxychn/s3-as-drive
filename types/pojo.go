package types

import "time"

type File struct {
	ID        uint      `json:"id"`
	Path      string    `json:"path"`
	Key       *string   `json:"key"`
	IsDir     bool      `json:"isDir"`
	Depth     uint      `json:"depth"`
	Size      *uint64   `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
