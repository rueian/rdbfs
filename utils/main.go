package utils

import (
	"fmt"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/jinzhu/gorm"
)

func ConvertDaoErr(err error) fuse.Status {
	if err == gorm.ErrRecordNotFound {
		return fuse.ENOENT
	}
	fmt.Println(err)
	return fuse.EIO
}
