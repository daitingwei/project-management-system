package dao

import (
	"context"

	"gorm.io/gorm"

	"test.com/project-user/internal/data/organization"
	"test.com/project-user/internal/database/gorms"
)

type OrganizationDao struct {
	conn *gorms.GormConn
}

func (o *OrganizationDao) FindOrganizationByMemId(ctx context.Context, memId int64) ([]*organization.Organization, error) {
	var orgs []*organization.Organization
	err := o.conn.Session(ctx).Where("member_id=?", memId).Find(&orgs).Error
	return orgs, err
}

func (o *OrganizationDao) FindOrganizationById(ctx context.Context, id int64) (*organization.Organization, error) {
	var org organization.Organization
	err := o.conn.Session(ctx).Where("id=?", id).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (o *OrganizationDao) UpdateOrganization(ctx context.Context, org *organization.Organization) error {
	return o.conn.Session(ctx).Model(&organization.Organization{}).
		Where("id=?", org.Id).
		Updates(map[string]interface{}{
			"name":    org.Name,
			"address": org.Address,
		}).Error
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{
		conn: gorms.New(),
	}
}

// SaveOrganization 直接接受 *gorm.DB 类型
func (o *OrganizationDao) SaveOrganization(tx *gorm.DB, ctx context.Context, org *organization.Organization) error {
	return tx.WithContext(ctx).Create(org).Error
}
