package domain

import (
	"context"
	"fmt"
	"time"

	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
)

type AccountDomain struct {
	accountRepo      repo.AccountRepo
	userRpcDomain    *UserRpcDomain
	departmentDomain *DepartmentDomain
}

func (d *AccountDomain) AccountList(
	organizationCode string,
	memberId int64,
	page int64,
	pageSize int64,
	departmentCode string,
	searchType int32) ([]*data.MemberAccountDisplay, int64, *errs.BError) {
	var condition string
	organizationCodeId := encrypts.DecryptNoErr(organizationCode)
	departmentCodeId := encrypts.DecryptNoErr(departmentCode)
	switch searchType {
	case 1:
		condition = "status = 1"
	case 2:
		condition = "department_code = 0"
	case 3:
		condition = "status = 0"
	case 4:
		condition = fmt.Sprintf("status = 1 and department_code = %d", departmentCodeId)
	default:
		condition = "status = 1"
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, total, err := d.accountRepo.FindList(c, condition, organizationCodeId, departmentCodeId, page, pageSize)
	if err != nil {
		return nil, 0, model.DBError
	}
	var dList []*data.MemberAccountDisplay
	for _, v := range list {
		display := v.ToDisplay()
		memberInfo, _ := d.userRpcDomain.MemberInfo(c, v.MemberCode)
		display.Avatar = memberInfo.Avatar
		if v.DepartmentCode > 0 {
			department, err := d.departmentDomain.FindDepartmentById(v.DepartmentCode)
			if err != nil {
				return nil, 0, err
			}
			display.Departments = department.Name
		}
		dList = append(dList, display)
	}
	return dList, total, nil
}

func (d *AccountDomain) FindAccount(memberId int64) (*data.MemberAccount, *errs.BError) {
	account, err := d.accountRepo.FindByMemberId(context.Background(), memberId)
	if err != nil {
		return nil, model.DBError
	}
	return account, nil
}

func (d *AccountDomain) FindAccountById(id int64) (*data.MemberAccount, *errs.BError) {
	ma, err := d.accountRepo.FindById(context.Background(), id)
	if err != nil {
		return nil, model.DBError
	}
	if ma != nil && ma.MemberCode > 0 {
		memberInfo, _ := d.userRpcDomain.MemberInfo(context.Background(), ma.MemberCode)
		if memberInfo != nil {
			ma.Avatar = memberInfo.Avatar
		}
	}
	return ma, nil
}

func (d *AccountDomain) SaveAccount(ctx context.Context, ma *data.MemberAccount) *errs.BError {
	if err := d.accountRepo.Save(ctx, ma); err != nil {
		return model.DBError
	}
	return nil
}

func (d *AccountDomain) UpdateAccountFields(ctx context.Context, id int64, fields map[string]interface{}) *errs.BError {
	if err := d.accountRepo.UpdateFields(ctx, id, fields); err != nil {
		return model.DBError
	}
	return nil
}

func (d *AccountDomain) ForbidAccount(ctx context.Context, id int64) *errs.BError {
	if err := d.accountRepo.UpdateStatus(ctx, id, 0); err != nil {
		return model.DBError
	}
	return nil
}

func (d *AccountDomain) ResumeAccount(ctx context.Context, id int64) *errs.BError {
	if err := d.accountRepo.UpdateStatus(ctx, id, 1); err != nil {
		return model.DBError
	}
	return nil
}

func (d *AccountDomain) AuthAccount(ctx context.Context, id int64, authorize string) *errs.BError {
	if err := d.accountRepo.UpdateAuthorize(ctx, id, authorize); err != nil {
		return model.DBError
	}
	return nil
}

func (d *AccountDomain) DeleteAccount(ctx context.Context, id int64) *errs.BError {
	if err := d.accountRepo.Delete(ctx, id); err != nil {
		return model.DBError
	}
	return nil
}

func (d *AccountDomain) AllList(ctx context.Context, organizationCode string) ([]*data.MemberAccountDisplay, *errs.BError) {
	orgId := encrypts.DecryptNoErr(organizationCode)
	list, err := d.accountRepo.FindAllByOrganization(ctx, orgId)
	if err != nil {
		return nil, model.DBError
	}
	var dList []*data.MemberAccountDisplay
	for _, v := range list {
		display := v.ToDisplay()
		memberInfo, _ := d.userRpcDomain.MemberInfo(ctx, v.MemberCode)
		if memberInfo != nil {
			display.Avatar = memberInfo.Avatar
		}
		dList = append(dList, display)
	}
	return dList, nil
}

func NewAccountDomain() *AccountDomain {
	return &AccountDomain{
		accountRepo:      dao.NewMemberAccountDao(),
		userRpcDomain:    NewUserRpcDomain(),
		departmentDomain: NewDepartmentDomain(),
	}
}
