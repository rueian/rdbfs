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

func formatDirPath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return path
}

func getPathAndNameFromFullPath(fullPath string) (string, string) {
	if fullPath == "/" {
		fullPath = ""
	}

	i := strings.LastIndex(fullPath, "/")

	return formatDirPath(fullPath[:i+1]), fullPath[i+1:]
}

func convertDaoErr(err error) fuse.Status {
	if err == gorm.ErrRecordNotFound {
		return fuse.ENOENT
	}
	fmt.Println(err)
	return fuse.EIO
}

func (fs *RdbFs) GetAttr(fullPath string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("GetAttr: ", path, name)

	if path == "/" && name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	attr, err := fs.Dao.GetAttr(path, name)
	if err != nil {
		return nil, convertDaoErr(err)
	}

	return attr, fuse.OK
}

func (fs *RdbFs) OpenDir(fullPath string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	fullPath = formatDirPath(fullPath)
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

func (fs *RdbFs) Mkdir(fullPath string, mode uint32, context *fuse.Context) fuse.Status {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Mkdir: ", path, name)

	_, err := fs.Dao.CreateObject(path, name, mode)
	if err != nil {
		return convertDaoErr(err)
	}

	return fuse.OK
}

func (fs *RdbFs) Rmdir(fullPath string, context *fuse.Context) fuse.Status {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Rmdir: ", path, name)

	err := fs.Dao.RemoveObject(path, name)
	if err != nil {
		return convertDaoErr(err)
	}

	err = fs.Dao.RemoveSubTree(path)
	if err != nil {
		return convertDaoErr(err)
	}

	return fuse.OK
}

func (fs *RdbFs) Rename(oldFullPath string, newFullPath string, context *fuse.Context) fuse.Status {
	oldPath, oldName := getPathAndNameFromFullPath(oldFullPath)
	newPath, newName := getPathAndNameFromFullPath(newFullPath)
	fmt.Println("Rename: ", oldFullPath, newFullPath)

	if err := fs.Dao.RenameObject(oldPath, oldName, newPath, newName); err != nil {
		return convertDaoErr(err)
	}

	if err := fs.Dao.RenameSubTree(formatDirPath(oldFullPath), formatDirPath(newFullPath)); err != nil {
		return convertDaoErr(err)
	}

	return fuse.OK
}

func (fs *RdbFs) Create(fullPath string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Create: ", path, name)

	object, err := fs.Dao.CreateObject(path, name, mode)
	if err != nil {
		return nil, convertDaoErr(err)
	}

	return object, fuse.OK
}

//Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status)
//Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status)
//Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status)
//Truncate(name string, size uint64, context *fuse.Context) (code fuse.Status)
//Access(name string, mode uint32, context *fuse.Context) (code fuse.Status)
