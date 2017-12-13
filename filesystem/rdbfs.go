package filesystem

import (
	"fmt"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/jinzhu/gorm"
	"github.com/rueian/rdbfs/model"
)

type RdbFs struct {
	pathfs.FileSystem

	Dao *model.Dao
}

func getPathAndNameFromFullPath(fullPath string) (string, string) {
	if fullPath == "/" {
		fullPath = ""
	}

	i := strings.LastIndex(fullPath, "/")

	return fullPath[:i+1], fullPath[i+1:]
}

func convertDaoErr(err error) fuse.Status {
	if err == gorm.ErrRecordNotFound {
		return fuse.ENOENT
	}
	fmt.Println(err)
	return fuse.EIO
}

func (fs *RdbFs) GetAttr(fullPath string, context *fuse.Context) (*fuse.Attr, fuse.Status) {

	if fullPath == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	path, name := getPathAndNameFromFullPath(fullPath)

	fmt.Println("GetAttr: ", path, name)

	attr, err := fs.Dao.GetAttr(path, name)
	if err != nil {
		return nil, convertDaoErr(err)
	}

	return attr, fuse.OK
}

func (fs *RdbFs) OpenDir(fullPath string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	fmt.Println("OpenDir: ", fullPath)

	objects, err := fs.Dao.GetSubTree(fullPath)
	if err != nil {
		return nil, convertDaoErr(err)
	}

	for _, object := range objects {
		c = append(c, fuse.DirEntry{
			Ino:  uint64(object.ID),
			Name: object.Name,
			Mode: object.Attr.Mode,
		})
	}

	return c, fuse.OK
}

func (fs *RdbFs) Open(fullPath string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Open: ", path, name)

	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return nil, convertDaoErr(err)
	}

	return object, fuse.OK
}
