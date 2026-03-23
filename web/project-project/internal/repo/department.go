package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type DepartmentRepo interface {
	FindDepartmentById(ctx context.Context, id int64) (*data.Department, error)
	FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*data.Department, error)
	Save(dpm *data.Department) error
	ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*data.Department, total int64, err error)
	Update(dpm *data.Department) error
	Delete(id int64) error
	FindChildDepartments(parentCode int64) (list []*data.Department, err error)
	FindAllDepartmentByOrg(organizationCode int64) (list []*data.Department, err error)
	SaveDepartmentMember(dpm *data.DepartmentMember) error
	FindDepartmentMember(ctx context.Context, departmentCode int64) (list []*data.DepartmentMember, err error)
	FindDepartmentMemberByAccount(ctx context.Context, departmentCode int64, accountCode int64) (*data.DepartmentMember, error)
	DeleteDepartmentMember(departmentCode int64, accountCode int64) error
	UpdateDepartmentMember(dpm *data.DepartmentMember) error
	FindDepartmentMembersByOrg(ctx context.Context, organizationCode int64) (list []*data.DepartmentMember, err error)
	FindMemberAccountsByOrg(ctx context.Context, organizationCode int64) (list []*data.MemberAccount, err error)
	FindMemberAccountsByOrgWithKeyword(ctx context.Context, organizationCode int64, keyword string) (list []*data.MemberAccount, err error)
	FindMemberByEmail(ctx context.Context, email string) (*data.Member, error)
	FindMemberAccountByMemberAndOrg(ctx context.Context, memberCode int64, organizationCode int64) (*data.MemberAccount, error)
	SaveMemberAccount(ma *data.MemberAccount) error
	FindMemberAccountsByIds(ctx context.Context, ids []int64) (list []*data.MemberAccount, err error)
	UpdateMemberAccountDepartment(ctx context.Context, accountId int64, departmentCode int64) error
}
