package model

import (
	"errors"
	"fmt"
	"time"

	"database/sql/driver"

	"encoding/json"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type ObjectAttr fuse.Attr

func (a *ObjectAttr) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ObjectAttr) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", src))
	}

	return json.Unmarshal(bytes, a)
}

type Object struct {
	Dao   *Dao       `gorm:"-"`
	ID    uint       `gorm:"primary_key"`
	Path  string     `gorm:"unique_index:idx_path_name"`
	Name  string     `gorm:"unique_index:idx_path_name"`
	Attr  ObjectAttr `gorm:"type:json"`
	Xattr []byte
	Data  []byte
}

func (*Object) SetInode(*nodefs.Inode) {
	panic("implement me")
}

func (*Object) String() string {
	panic("implement me")
}

func (*Object) InnerFile() nodefs.File {
	panic("implement me")
}

func (*Object) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	panic("implement me")
}

func (*Object) Write(data []byte, off int64) (written uint32, code fuse.Status) {
	panic("implement me")
}

func (*Object) Flock(flags int) fuse.Status {
	panic("implement me")
}

func (*Object) Flush() fuse.Status {
	panic("implement me")
}

func (*Object) Release() {
	panic("implement me")
}

func (*Object) Fsync(flags int) (code fuse.Status) {
	panic("implement me")
}

func (*Object) Truncate(size uint64) fuse.Status {
	panic("implement me")
}

func (*Object) GetAttr(out *fuse.Attr) fuse.Status {
	panic("implement me")
}

func (*Object) Chown(uid uint32, gid uint32) fuse.Status {
	panic("implement me")
}

func (*Object) Chmod(perms uint32) fuse.Status {
	panic("implement me")
}

func (*Object) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	panic("implement me")
}

func (*Object) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	panic("implement me")
}
