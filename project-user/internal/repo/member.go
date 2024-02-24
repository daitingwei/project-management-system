package repo

import (
	"context"
	"gorm.io/gorm"

	"test.com/project-user/internal/data/member"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(tx *gorm.DB, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, pwd string) (mem *member.Member, err error)
	FindMemberById(background context.Context, id int64) (mem *member.Member, err error)
	FindMemberByIds(background context.Context, ids []int64) (list []*member.Member, err error)
	// FindMemberByMobile 根据手机号查询会员
	// 修改说明：新增此接口方法，用于通过手机号查询会员信息
	// 背景：支持手机验证码登录，需要根据手机号查询用户
	FindMemberByMobile(ctx context.Context, mobile string) (mem *member.Member, err error)
	UpdateMemberFields(ctx context.Context, id int64, fields map[string]interface{}) error
}
