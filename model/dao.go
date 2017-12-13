package model

import (
	"errors"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Dao struct {
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

	attr := fuse.Attr(object.Attr)

	return &attr, nil
}

func (d *Dao) GetSubTree(path string) ([]*Object, error) {
	var objects []*Object
	if err := d.DbConn.Select("id, name, attr").Where("path = ?", path).Find(&objects).Error; err != nil {
		return nil, err
	}
	return objects, nil
}

func (d *Dao) GetObject(path, name string) (*Object, error) {
	object := &Object{}
	if err := d.DbConn.Select("id, attr").Where("path = ?", path).Where("name = ?", name).First(object).Error; err != nil {
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
			Mode: mode,
		},
	}
	if err := d.DbConn.Create(object).Error; err != nil {
		return nil, err
	}

	return object, nil
}

func (d *Dao) RemoveObject(path, name string) error {
	if err := d.DbConn.Where("path = ?", path).Where("name = ?", name).Delete(Object{}).Error; err != nil {
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
	if err := d.DbConn.Model(Object{}).Where("path = ?", oldPath).Where("name = ?", oldName).Updates(map[string]interface{}{
		"path": newPath,
		"name": newName,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) RenameSubTree(oldPath, newPath string) error {
	if err := d.DbConn.Model(Object{}).Where("path = ?", oldPath).Update("path", newPath).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) ReadBytes(path, name string, dest []byte, off int64) ([]byte, error) {
	return nil, nil
}

func (d *Dao) WriteBytes(path, name string, data []byte, off int64) ([]byte, error) {
	return nil, nil
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
		return &Dao{dbConn}, nil
	}
}
