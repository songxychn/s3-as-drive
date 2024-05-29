package types

import "time"

type File struct {
	ID        uint      `json:"id"`
	Path      string    `json:"path" gorm:"index;not null"`
	Key       *string   `json:"key"`
	IsDir     bool      `json:"isDir" gorm:"not null"`
	Depth     uint      `json:"depth" gorm:"index;not null"`
	Size      *uint64   `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SyncDir struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name" gorm:"uniqueIndex;not null"`
	Path      string    `json:"path" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
