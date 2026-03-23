package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type ProjectAuthDao struct {
	conn *gorms.GormConn
}

func (p *ProjectAuthDao) FindAuthListPage(ctx context.Context, orgCode int64, page int64, pageSize int64) (list []*data.ProjectAuth, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.Model(&data.ProjectAuth{}).
		Where("organization_code=?", orgCode).
		Limit(int(pageSize)).
		Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	err = session.Model(&data.ProjectAuth{}).
		Where("organization_code=?", orgCode).
		Count(&total).Error
	return
}

func (p *ProjectAuthDao) FindAuthList(ctx context.Context, orgCode int64) (list []*data.ProjectAuth, err error) {
	session := p.conn.Session(ctx)
	err = session.Model(&data.ProjectAuth{}).Where("organization_code=?", orgCode).Find(&list).Error
	return
}

func NewProjectAuthDao() *ProjectAuthDao {
	return &ProjectAuthDao{
		conn: gorms.New(),
	}
}

func (p *ProjectAuthDao) SaveAuth(ctx context.Context, pa *data.ProjectAuth) error {
	return p.conn.Session(ctx).Save(pa).Error
}

func (p *ProjectAuthDao) UpdateAuth(ctx context.Context, pa *data.ProjectAuth) error {
	return p.conn.Session(ctx).Model(&data.ProjectAuth{}).Where("id=?", pa.Id).Updates(pa).Error
}

func (p *ProjectAuthDao) SetDefault(ctx context.Context, orgCode int64, authId int64) error {
	session := p.conn.Session(ctx)
	// 先清除该组织所有默认
	if err := session.Model(&data.ProjectAuth{}).Where("organization_code=?", orgCode).Update("is_default", 0).Error; err != nil {
		return err
	}
	// 设置指定 auth 为默认
	return session.Model(&data.ProjectAuth{}).Where("id=?", authId).Update("is_default", 1).Error
}

func (p *ProjectAuthDao) DeleteAuth(ctx context.Context, authId int64) error {
	return p.conn.Session(ctx).Where("id=?", authId).Delete(&data.ProjectAuth{}).Error
}

func (p *ProjectAuthDao) FindDefaultAuth(ctx context.Context, orgCode int64) (*data.ProjectAuth, error) {
	var auth data.ProjectAuth
	err := p.conn.Session(ctx).Where("organization_code=? AND is_default=1", orgCode).Take(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

