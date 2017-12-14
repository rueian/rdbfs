package filesystem

import (
	"fmt"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/rueian/rdbfs/model"
	"github.com/rueian/rdbfs/utils"
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

func (fs *RdbFs) GetAttr(fullPath string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)
	//fmt.Println("GetAttr: ", path, name)

	if path == "/" && name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	attr, err := fs.Dao.GetAttr(path, name)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	return attr, fuse.OK
}

func (fs *RdbFs) OpenDir(fullPath string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	fullPath = formatDirPath(fullPath)
	//fmt.Println("OpenDir: ", fullPath)

	objects, err := fs.Dao.GetSubTree(fullPath)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
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

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	return object, fuse.OK
}

func (fs *RdbFs) Mkdir(fullPath string, mode uint32, context *fuse.Context) fuse.Status {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Mkdir: ", path, name)

	_, err := fs.Dao.CreateObject(path, name, mode)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	return fuse.OK
}

func (fs *RdbFs) Rmdir(fullPath string, context *fuse.Context) fuse.Status {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Rmdir: ", path, name)

	if err := fs.Dao.RemoveObject(path, name); err != nil {
		return utils.ConvertDaoErr(err)
	}
	return fuse.OK
}

func (fs *RdbFs) Rename(oldFullPath string, newFullPath string, context *fuse.Context) fuse.Status {
	oldPath, oldName := getPathAndNameFromFullPath(oldFullPath)
	newPath, newName := getPathAndNameFromFullPath(newFullPath)
	fmt.Println("Rename: ", oldFullPath, newFullPath)

	if err := fs.Dao.RenameObject(oldPath, oldName, newPath, newName); err != nil {
		return utils.ConvertDaoErr(err)
	}

	if err := fs.Dao.RenameSubTree(formatDirPath(oldFullPath), formatDirPath(newFullPath)); err != nil {
		return utils.ConvertDaoErr(err)
	}

	return fuse.OK
}

func (fs *RdbFs) Create(fullPath string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Create: ", path, name)

	object, err := fs.Dao.CreateObject(path, name, mode)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	return object, fuse.OK
}

func (fs *RdbFs) Unlink(fullPath string, context *fuse.Context) (code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Unlink: ", path, name)

	if err := fs.Dao.RemoveObject(path, name); err != nil {
		return utils.ConvertDaoErr(err)
	}

	return fuse.OK
}

func (fs *RdbFs) Truncate(fullPath string, size uint64, context *fuse.Context) (code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	return object.Truncate(size)
}

//func (fs *RdbFs) String() string {
//	panic("implement me")
//}
//
//func (fs *RdbFs) SetDebug(debug bool) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
//	panic("implement me")
//}
//
//func (fs *RdbFs) GetXAttr(name string, attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) ListXAttr(name string, context *fuse.Context) (attributes []string, code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
//	panic("implement me")
//}
//
//func (fs *RdbFs) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
//	panic("implement me")
//}
//
//func (fs *RdbFs) OnMount(nodeFs *pathfs.PathNodeFs) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) OnUnmount() {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) StatFs(name string) *fuse.StatfsOut {
//	panic("implement me")
//}
