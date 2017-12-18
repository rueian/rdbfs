package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

type Dao struct {
	Driver string
	DbConn *gorm.DB
}

func (d *Dao) Close() {
	d.DbConn.Close()
}

func (d *Dao) AutoMigrate() error {
	if err := d.DbConn.AutoMigrate(&Object{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) DropTable() error {
	if err := d.DbConn.DropTableIfExists(&Object{}).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetAttr(path, name string) (*fuse.Attr, error) {
	object, err := d.GetObject(path, name)
	if err != nil {
		return &fuse.Attr{}, err
	}
	return &(object.Attr.Attr), nil
}

func (d *Dao) UpdateAttr(id uint, attr ObjectAttr) error {
	if err := d.DbConn.Model(&Object{}).Where("id = ?", id).Update("attr", attr).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) UpdateXattr(id uint, xattr map[string]string) error {
	if err := d.DbConn.Model(&Object{}).Where("id = ?", id).Update("xattr", xattr).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) Link(fid uint, tid uint) error {
	if err := d.DbConn.Model(&Object{}).Where("id = ?", fid).Update("link_id", tid).Error; err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetSubTree(path string) ([]*Object, error) {
	var objects []*Object
	if err := d.DbConn.Select("id, name, attr").Where("path = ?", path).Find(&objects).Error; err != nil {
		return nil, err
	}
	return objects, nil
}

func (d *Dao) GetLinkId(path, name string) (int64, error) {
	object := &Object{}
	if err := d.DbConn.Select("link_id").Where("path = ?", path).Where("name = ?", name).First(object).Error; err != nil {
		return 0, err
	}
	return object.LinkID, nil
}

func (d *Dao) GetObject(path, name string) (*Object, error) {
	object := &Object{}
	if err := d.DbConn.Select("id, attr, xattr, link_id").Where("path = ?", path).Where("name = ?", name).First(object).Error; err != nil {
		return nil, err
	}

	// Follow links chain
	for object.LinkID != 0 {
		obj, err := d.GetObjectById(uint(object.LinkID))
		if err != nil {
			return nil, err
		}
		object = obj
	}

	object.Dao = d

	return object, nil
}

func (d *Dao) GetObjectById(id uint) (*Object, error) {
	object := &Object{}
	if err := d.DbConn.Select("id, attr, xattr, link_id").Where("id = ?", id).First(object).Error; err != nil {
		return nil, err
	}

	object.Dao = d

	return object, nil
}

func (d *Dao) CreateObject(path, name string, mode uint32) (*Object, error) {
	object := &Object{
		Dao:  d,
		Path: path,
		Name: name,
		Attr: ObjectAttr{
			fuse.Attr{
				Mode: mode,
			},
		},
	}

	now := time.Now()
	object.Attr.SetTimes(&now, &now, &now)
	if err := d.DbConn.Create(object).Error; err != nil {
		return nil, err
	}

	return object, nil
}

// Unlink from fs
func (d *Dao) UnlinkName(path, name string) error {
	if err := d.DbConn.Model(Object{}).Where("path = ?", path).Where("name = ?", name).Update(map[string]interface{}{
		"path": "..",
		"name": uuid.NewV4(),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) HasLinkedObject(path, name string) (bool, error) {
	count := 0
	err := d.DbConn.Model(Object{}).Where("link_id = ?", gorm.Expr("(select id from objects where path = ? and name = ?)", path, name)).Count(&count).Error
	// d.DbConn.Table("objects").Select("id").Where("path = ?", path).Where("name = ?", name).QueryExpr()
	fmt.Println("HasLinkedObject", count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (d *Dao) RemoveObject(path, name string) error {
	if err := d.DbConn.Where("path = ?", path).Where("name = ?", name).Delete(Object{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) RemoveObjectById(id uint) error {
	if err := d.DbConn.Where("id = ?", id).Delete(Object{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) RemoveSubTree(path string) error {
	if err := d.DbConn.Where("path = ?", path).Delete(Object{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) RenameObject(oldPath, oldName, newPath, newName string) error {
	if err := d.DbConn.Model(&Object{}).Where("path = ?", oldPath).Where("name = ?", oldName).Updates(map[string]interface{}{
		"path": newPath,
		"name": newName,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) RenameSubTree(oldPath, newPath string) error {
	if err := d.DbConn.Model(&Object{}).Where("path = ?", oldPath).Update("path", newPath).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) ReadBytes(id uint, dest []byte, off int64) ([]byte, error) {
	var row *sql.Row

	row = d.DbConn.Model(&Object{}).Where("id = ?", id).Select("substring(data, ?, ?)", off+1, len(dest)).Row()

	if err := row.Scan(&dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (d *Dao) WriteBytes(id uint, data []byte, off int64) (uint32, error) {
	var err error

	if off == 0 {
		err = d.DbConn.Model(&Object{}).Where("id = ?", id).Update("data", data).Error
	} else {
		if d.Driver == "postgres" {
			err = d.DbConn.Model(&Object{}).Where("id = ?", id).Update("data", gorm.Expr("overlay(data placing ? from ?)", data, off)).Error
		}
		if d.Driver == "mysql" {
			err = d.DbConn.Model(&Object{}).Where("id = ?", id).Update("data", gorm.Expr("INSERT(data, ?, ?, ?)", off, len(data), data)).Error
		}
	}

	return uint32(len(data)), err
}

func (d *Dao) Truncate(id uint, size uint64) error {
	var err error

	if d.Driver == "postgres" {
		err = d.DbConn.Model(&Object{}).Where("id = ?", id).Update("data", gorm.Expr("substring(data from 0 for ?)", size)).Error
	}
	if d.Driver == "mysql" {
		err = d.DbConn.Model(&Object{}).Where("id = ?", id).Update("data", gorm.Expr("SUBSTRING(data, 0, ?)", size)).Error
	}

	return err
}

var supportedDatabase = map[string]bool{
	"mysql":    true,
	"postgres": true,
}

func NewDao(driver, url string) (*Dao, error) {
	if !supportedDatabase[driver] {
		return nil, errors.New("driver not supported")
	}

	if dbConn, err := gorm.Open(driver, url); err != nil {
		return nil, err
	} else {
		return &Dao{driver, dbConn}, nil
	}
}
