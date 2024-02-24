package dao

import (
	"context"

	"gorm.io/gorm"

	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/database/gorms"
)

type MemberDao struct {
	conn *gorms.GormConn
}

func NewMemberDao() *MemberDao {
	return &MemberDao{
		conn: gorms.New(),
	}
}

func (m *MemberDao) FindMemberByIds(background context.Context, ids []int64) (list []*member.Member, err error) {
	if len(ids) <= 0 {
		return nil, nil
	}
	err = m.conn.Session(background).Model(&member.Member{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func (m *MemberDao) FindMemberById(ctx context.Context, id int64) (mem *member.Member, err error) {
	err = m.conn.Session(ctx).Where("id=?", id).First(&mem).Error
	return
}

func (m *MemberDao) FindMember(ctx context.Context, account string, pwd string) (*member.Member, error) {
	var mem member.Member
	err := m.conn.Session(ctx).Where("account=? and password=?", account, pwd).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mem, err
}

// SaveMember 直接接受 *gorm.DB 类型
func (m *MemberDao) SaveMember(tx *gorm.DB, ctx context.Context, mem *member.Member) error {
	return tx.WithContext(ctx).Create(mem).Error
}

func (m *MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("email=?", email).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) GetMemberByAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("account=?", account).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("mobile=?", mobile).Count(&count).Error
	return count > 0, err
}

// FindMemberByMobile 根据手机号查询会员
// 修改说明：新增此方法，用于通过手机号查询会员信息
// 背景：支持手机验证码登录，需要根据手机号查询用户
// 实现：根据手机号查询ms_member表，返回会员信息
func (m *MemberDao) FindMemberByMobile(ctx context.Context, mobile string) (*member.Member, error) {
	var mem member.Member
	err := m.conn.Session(ctx).Where("mobile=?", mobile).First(&mem).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mem, err
}

func (m *MemberDao) UpdateMemberFields(ctx context.Context, id int64, fields map[string]interface{}) error {
	return m.conn.Session(ctx).Model(&member.Member{}).Where("id = ?", id).Updates(fields).Error
}
