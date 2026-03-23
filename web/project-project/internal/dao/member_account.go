package dao

import (
	"context"

	"gorm.io/gorm"

	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type MemberAccountDao struct {
	conn *gorms.GormConn
}

func (m *MemberAccountDao) FindByMemberId(ctx context.Context, memberId int64) (ma *data.MemberAccount, err error) {
	session := m.conn.Session(ctx)
	err = session.Where("member_code=?", memberId).Take(&ma).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (m *MemberAccountDao) FindList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) (list []*data.MemberAccount, total int64, err error) {
	session := m.conn.Session(ctx)
	offset := (page - 1) * pageSize
	err = session.Model(&data.MemberAccount{}).
		Where("organization_code=?", organizationCode).
		Where(condition).Limit(int(pageSize)).Offset(int(offset)).Find(&list).Error
	err = session.Model(&data.MemberAccount{}).
		Where("organization_code=?", organizationCode).
		Where(condition).Count(&total).Error
	return
}

func (m *MemberAccountDao) FindById(ctx context.Context, id int64) (*data.MemberAccount, error) {
	session := m.conn.Session(ctx)
	var ma data.MemberAccount
	err := session.Where("id=?", id).Take(&ma).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ma, err
}

func (m *MemberAccountDao) Save(ctx context.Context, ma *data.MemberAccount) error {
	session := m.conn.Session(ctx)
	return session.Save(ma).Error
}

func (m *MemberAccountDao) UpdateFields(ctx context.Context, id int64, fields map[string]interface{}) error {
	session := m.conn.Session(ctx)
	return session.Model(&data.MemberAccount{}).Where("id=?", id).Updates(fields).Error
}

func (m *MemberAccountDao) UpdateStatus(ctx context.Context, id int64, status int) error {
	session := m.conn.Session(ctx)
	return session.Model(&data.MemberAccount{}).Where("id=?", id).Update("status", status).Error
}

func (m *MemberAccountDao) UpdateAuthorize(ctx context.Context, id int64, authorize string) error {
	session := m.conn.Session(ctx)
	return session.Model(&data.MemberAccount{}).Where("id=?", id).Update("authorize", authorize).Error
}

func (m *MemberAccountDao) Delete(ctx context.Context, id int64) error {
	session := m.conn.Session(ctx)
	return session.Where("id=?", id).Delete(&data.MemberAccount{}).Error
}

func (m *MemberAccountDao) FindAllByOrganization(ctx context.Context, organizationCode int64) ([]*data.MemberAccount, error) {
	session := m.conn.Session(ctx)
	var list []*data.MemberAccount
	err := session.Where("organization_code=?", organizationCode).Find(&list).Error
	return list, err
}

func NewMemberAccountDao() *MemberAccountDao {
	return &MemberAccountDao{
		conn: gorms.New(),
	}
}
