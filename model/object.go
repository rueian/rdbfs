package model

import (
	"time"
	"github.com/hanwen/go-fuse/fuse"
)

const (
	OBJECT_TYPE_DIR = 0
	OBJECT_TYPE_FILE = 1
	OBJECT_TYPE_SYMLINK = 2
)

type Object struct {
	ID        uint   `gorm:"primary_key"`
	Type      uint
	Path      string `gorm:"unique_index:idx_path_name"`
	Name      string `gorm:"unique_index:idx_path_name"`
	Attr      fuse.Attr
	XAttr     []byte
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}