package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/rueian/rdbfs/utils"
)

// the larger the faster
var syncSize = 81920000

type ObjectAttr struct {
	fuse.Attr
}

func (a ObjectAttr) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ObjectAttr) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", src))
	}

	return json.Unmarshal(b, a)
}

type Object struct {
	nodefs.File
	Dao        *Dao       `gorm:"-"`
	ID         uint       `gorm:"primary_key"`
	Path       string     `gorm:"unique_index:idx_path_name"`
	Name       string     `gorm:"unique_index:idx_path_name"`
	Attr       ObjectAttr `gorm:"type:json"`
	Xattr      []byte
	Data       []byte
	FBuffer    bytes.Buffer `gorm:"-"`
	FBufOffset int64        `gorm:"-"`
	CurrOffset int64        `gorm:"-"`
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

func (o *Object) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	fmt.Println("Read", off, len(dest))
	res, err := o.Dao.ReadBytes(o.ID, dest, off)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	return fuse.ReadResultData(res), fuse.OK
}

func (o *Object) Write(data []byte, off int64) (written uint32, code fuse.Status) {
	fmt.Println("Write", off, len(data))

	o.CurrOffset = off
	over := int(o.CurrOffset - o.FBufOffset)
	if o.CurrOffset > o.FBufOffset && over <= o.FBuffer.Len() {
		// clear unchecked data start from FBufOffset
		o.FBuffer.Truncate(over)
	} else if o.FBuffer.Len() == 0 {
		// all the previous data has been written into DB
		o.FBufOffset = o.CurrOffset
	}

	// write into per object buffer
	n, err := o.FBuffer.Write(data)
	if err != nil {
		return 0, utils.ConvertDaoErr(err)
	}

	// if the buffer is large enough to be written into DB
	if o.FBuffer.Len() > syncSize {
		status := o.Flush()
		if status != fuse.OK {
			return 0, status
		}
	}

	return uint32(n), fuse.OK
}

func (*Object) Flock(flags int) fuse.Status {
	fmt.Println("implement Flock")
	return fuse.OK
}

func (o *Object) Flush() fuse.Status {
	fmt.Println("Flush (temporarily call Fsync())")
	return o.Fsync(0)
}

func (*Object) Release() {
	fmt.Println("implement Release")
}

func (o *Object) Fsync(flags int) (code fuse.Status) {
	fmt.Println("Fsync", int64(o.FBuffer.Len()))

	if o.FBuffer.Len() != 0 {
		// write the data in per object buffer into DB
		written, err := o.Dao.WriteBytes(o.ID, o.FBuffer.Bytes(), o.FBufOffset)
		if err != nil {
			return utils.ConvertDaoErr(err)
		}

		o.Attr.Size = uint64(written) + uint64(o.FBufOffset)
	}

	o.FBufOffset = o.CurrOffset
	o.CurrOffset = 0
	o.FBuffer.Reset()

	if err := o.Dao.UpdateAttr(o.ID, o.Attr); err != nil {
		return utils.ConvertDaoErr(err)
	}

	return fuse.OK
}

func (o *Object) Truncate(size uint64) fuse.Status {
	fmt.Println("Truncate", size)

	if err := o.Dao.Truncate(o.ID, size); err != nil {
		return utils.ConvertDaoErr(err)
	}

	if size != o.Attr.Size {
		o.Attr.Size = size
		if err := o.Dao.UpdateAttr(o.ID, o.Attr); err != nil {
			return utils.ConvertDaoErr(err)
		}
	}

	return fuse.OK
}

func (o *Object) GetAttr(out *fuse.Attr) fuse.Status {
	objValue := reflect.ValueOf(o.Attr.Attr)
	outValue := reflect.ValueOf(out).Elem()
	for i := 0; i < outValue.NumField(); i++ {
		outValue.Field(i).Set(objValue.Field(i))
	}
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
