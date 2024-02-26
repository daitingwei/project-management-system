package dao

import (
	"context"
	"test.com/project-common/encrypts"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type SourceLinkDao struct {
	conn *gorms.GormConn
}

// func (s *SourceLinkDao) Save(ctx context.Context, link *data.SourceLink) error {
// 	return s.conn.Session(ctx).Save(&link).Error
// }

func (s *SourceLinkDao) SaveTx(ctx context.Context, conn database.DbConn, link *data.SourceLink) error {
	s.conn = conn.(*gorms.GormConn)
	err := s.conn.Tx(ctx).Save(&link).Error
	return err
}

func (s *SourceLinkDao) FindByTaskCode(ctx context.Context, taskCode int64) (list []*data.SourceLink, err error) {
	session := s.conn.Session(ctx)
	// TODO: 这里"task"写死了，需要改为动态传入以支持其他模块的文件查询
	err = session.Model(&data.SourceLink{}).Where("link_type=? and link_code=?", "task", taskCode).Find(&list).Error
	return
}

func NewSourceLinkDao() *SourceLinkDao {
	return &SourceLinkDao{
		conn: gorms.New(),
	}
}

func (s *SourceLinkDao) DeleteByCode(ctx context.Context, code string) error {
	session := s.conn.Session(ctx)
	id := encrypts.DecryptNoErr(code)
	return session.Where("id=?", id).Delete(&data.SourceLink{}).Error
}
