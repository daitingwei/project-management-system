package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type AccountRepo interface {
	FindList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) ([]*data.MemberAccount, int64, error)
	FindByMemberId(ctx context.Context, memberId int64) (*data.MemberAccount, error)
	FindById(ctx context.Context, id int64) (*data.MemberAccount, error)
	Save(ctx context.Context, ma *data.MemberAccount) error
	UpdateFields(ctx context.Context, id int64, fields map[string]interface{}) error
	UpdateStatus(ctx context.Context, id int64, status int) error
	UpdateAuthorize(ctx context.Context, id int64, authorize string) error
	Delete(ctx context.Context, id int64) error
	FindAllByOrganization(ctx context.Context, organizationCode int64) ([]*data.MemberAccount, error)
}
