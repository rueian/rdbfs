package filesystem

import (
	"fmt"
	"strings"
	"time"

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
	fmt.Println("Create: ", path, name, mode)

	object, err := fs.Dao.CreateObject(path, name, mode)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	object.Attr.Ino = uint64(object.ID)
	err = fs.Dao.UpdateAttr(object.ID, object.Attr)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	return object, fuse.OK
}

func (fs *RdbFs) Unlink(fullPath string, context *fuse.Context) (code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)
	fmt.Println("Unlink: ", path, name)
	hasLinked, err := fs.Dao.HasLinkedObject(path, name)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	fmt.Println("Has linked? ", hasLinked)
	if hasLinked {
		if err := fs.Dao.UnlinkName(path, name); err != nil {
			return utils.ConvertDaoErr(err)
		}
	} else {
		linkTo, err := fs.Dao.GetLinkId(path, name)
		if err != nil {
			return utils.ConvertDaoErr(err)
		}
		// Last one link to object, remove old object too
		if linkTo != 0 {
			err = fs.Dao.RemoveObjectById(uint(linkTo))
			if err != nil {
				return utils.ConvertDaoErr(err)
			}
		}

		if err := fs.Dao.RemoveObject(path, name); err != nil {
			return utils.ConvertDaoErr(err)
		}
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
func (fs *RdbFs) Chmod(fullPath string, mode uint32, context *fuse.Context) (code fuse.Status) {
	fmt.Println("Chmod", mode)
	path, name := getPathAndNameFromFullPath(fullPath)

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	return object.Chmod(mode)
}

func (fs *RdbFs) Chown(fullPath string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	fmt.Println("Chown", uid, gid)
	path, name := getPathAndNameFromFullPath(fullPath)

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	return object.Chown(uid, gid)
}

func (fs *RdbFs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	panic("implement me")
}

//
//func (fs *RdbFs) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
//	panic("implement me")
//}
//
func (fs *RdbFs) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	fmt.Println("Link implement me", oldName, newName)
	path, name := getPathAndNameFromFullPath(newName)

	obj, err := fs.Dao.CreateObject(path, name, 420)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	oldPath, oldName := getPathAndNameFromFullPath(oldName)
	object, err := fs.Dao.GetObject(oldPath, oldName)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	err = obj.Dao.Link(obj.ID, object.ID)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	obj.Attr.Ino = uint64(object.ID)
	err = fs.Dao.UpdateAttr(obj.ID, obj.Attr)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	return fuse.OK
}

//
//func (fs *RdbFs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
//	panic("implement me")
//}
//
func (fs *RdbFs) GetXAttr(fullPath string, attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	return []byte(object.Xattr[attribute]), fuse.OK
}

func (fs *RdbFs) ListXAttr(fullPath string, context *fuse.Context) (attributes []string, code fuse.Status) {
	path, name := getPathAndNameFromFullPath(fullPath)

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return nil, utils.ConvertDaoErr(err)
	}

	for k := range object.Xattr {
		attributes = append(attributes, k)
	}

	return attributes, fuse.OK
}

func (fs *RdbFs) RemoveXAttr(fullPath string, attr string, context *fuse.Context) fuse.Status {
	path, name := getPathAndNameFromFullPath(fullPath)

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	delete(object.Xattr, attr)

	err = fs.Dao.UpdateXattr(object.ID, object.Xattr)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	return fuse.OK
}

func (fs *RdbFs) SetXAttr(fullPath string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	path, name := getPathAndNameFromFullPath(fullPath)

	object, err := fs.Dao.GetObject(path, name)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	if object.Xattr == nil {
		object.Xattr = model.ObjectXattr{}
	}
	object.Xattr[attr] = string(data)

	err = fs.Dao.UpdateXattr(object.ID, object.Xattr)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	return fuse.OK
}

//
//func (fs *RdbFs) OnMount(nodeFs *pathfs.PathNodeFs) {
//	panic("implement me")
//}
//
//func (fs *RdbFs) OnUnmount() {
//	panic("implement me")
//}
//
func (fs *RdbFs) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {
	fmt.Println("Symlink implement me", value, linkName)
	path, name := getPathAndNameFromFullPath(linkName)

	obj, err := fs.Dao.CreateObject(path, name, fuse.S_IFLNK|420)
	if err != nil {
		return utils.ConvertDaoErr(err)
	}

	written, code := obj.Write([]byte(value), 0)
	if written < 0 || code != fuse.OK {
		return fuse.EIO
	}

	return obj.Flush()
}

//
func (fs *RdbFs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	fmt.Println("Readlink implement me", name)
	file, code := fs.Open(name, 0, context)
	if code != fuse.OK {
		return "", code
	}

	attr := fuse.Attr{}
	code = file.GetAttr(&attr)
	if code != fuse.OK {
		return "", code
	}

	buf := make([]byte, attr.Size)
	rst, code := file.Read(buf, 0)
	content, code := rst.Bytes(buf)
	fmt.Println("readed in readlink", string(buf), string(content))
	if code != fuse.OK {
		return "", code
	}

	return string(content), fuse.OK
}

//
//func (fs *RdbFs) StatFs(name string) *fuse.StatfsOut {
//	panic("implement me")
//}
