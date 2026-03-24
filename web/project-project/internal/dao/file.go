package dao

import (
	"context"
	"time"

	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type FileDao struct {
	conn *gorms.GormConn
}

func (f *FileDao) FindByIds(ctx context.Context, ids []int64) (list []*data.File, err error) {
	session := f.conn.Session(ctx)
	err = session.Model(&data.File{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func (f *FileDao) FindByProjectCode(ctx context.Context, projectCode int64, deleted int, page int64, pageSize int64) (list []*data.File, total int64, err error) {
	session := f.conn.Session(ctx)
	offset := (page - 1) * pageSize
	err = session.Model(&data.File{}).Where("project_code = ? and deleted = ?", projectCode, deleted).
		Count(&total).Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error
	return
}

func (f *FileDao) FindById(ctx context.Context, id int64) (*data.File, error) {
	session := f.conn.Session(ctx)
	var file data.File
	err := session.Model(&data.File{}).Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (f *FileDao) UpdateTitle(ctx context.Context, id int64, title string) error {
	session := f.conn.Session(ctx)
	return session.Model(&data.File{}).Where("id = ?", id).Update("title", title).Error
}

func (f *FileDao) UpdateDeleted(ctx context.Context, id int64, deleted int) error {
	session := f.conn.Session(ctx)
	return session.Model(&data.File{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted":      deleted,
		"deleted_time": func() int64 {
			if deleted == 1 {
				return time.Now().UnixMilli()
			}
			return 0
		}(),
	}).Error
}

func (f *FileDao) DeleteById(ctx context.Context, id int64) error {
	session := f.conn.Session(ctx)
	return session.Where("id = ?", id).Delete(&data.File{}).Error
}

// func (f *FileDao) Save(ctx context.Context, file *data.File) error {
// 	err := f.conn.Session(ctx).Save(&file).Error
// 	return err
// }

func (f *FileDao) SaveTx(ctx context.Context, conn database.DbConn, file *data.File) error {
	f.conn = conn.(*gorms.GormConn)
	err := f.conn.Tx(ctx).Save(&file).Error
	return err
}

func NewFileDao() *FileDao {
	return &FileDao{
		conn: gorms.New(),
	}
}
