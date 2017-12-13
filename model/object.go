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

func (a ObjectAttr) Value() (driver.Value, error) {
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
	nodefs.File
	Dao   *Dao       `gorm:"-"`
	ID    uint       `gorm:"primary_key"`
	Path  string     `gorm:"unique_index:idx_path_name"`
	Name  string     `gorm:"unique_index:idx_path_name"`
	Attr  ObjectAttr `gorm:"type:json"`
	Xattr []byte
	Data  []byte
}

func (*Object) SetInode(*nodefs.Inode) {
	fmt.Println("SetInode")
}

func (o *Object) String() string {
	return o.Path + "/" + o.Name
}

func (*Object) InnerFile() nodefs.File {
	fmt.Println("implement InnerFile")
	return nil
}

func (*Object) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	fmt.Println("implement Read")
	return nil, fuse.OK
}

func (*Object) Write(data []byte, off int64) (written uint32, code fuse.Status) {
	fmt.Println("implement Write")
	return 0, fuse.OK
}

func (*Object) Flock(flags int) fuse.Status {
	fmt.Println("implement Flock")
	return fuse.OK
}

func (*Object) Flush() fuse.Status {
	fmt.Println("implement Flush")
	return fuse.OK
}

func (*Object) Release() {
	fmt.Println("implement Release")
}

func (*Object) Fsync(flags int) (code fuse.Status) {
	fmt.Println("implement Fsync")
	return fuse.OK
}

func (*Object) Truncate(size uint64) fuse.Status {
	fmt.Println("implement Truncate")
	return fuse.OK
}

func (o *Object) GetAttr(out *fuse.Attr) fuse.Status {
	fmt.Println("implement GetAttr")
	return fuse.OK
}

func (*Object) Chown(uid uint32, gid uint32) fuse.Status {
	fmt.Println("implement Chown")
	return fuse.OK
}

func (*Object) Chmod(perms uint32) fuse.Status {
	fmt.Println("implement Chmod")
	return fuse.OK
}

func (o *Object) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	fmt.Println("implement Utimens")
	return fuse.OK
}

func (*Object) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	fmt.Println("implement Allocate")
	return fuse.OK
}