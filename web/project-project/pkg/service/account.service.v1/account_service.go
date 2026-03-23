package account_service_v1

import (
	"context"
	"strconv"

	"github.com/jinzhu/copier"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-grpc/account"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/domain"
	"test.com/project-project/internal/interceptor"
	"test.com/project-project/internal/repo"
)

type AccountService struct {
	account.UnimplementedAccountServiceServer
	cache             repo.Cache
	transaction       tran.Transaction
	accountDomain     *domain.AccountDomain
	projectAuthDomain *domain.ProjectAuthDomain
}

func New() *AccountService {
	return &AccountService{
		cache:             dao.Rc,
		transaction:       dao.NewTransaction(),
		accountDomain:     domain.NewAccountDomain(),
		projectAuthDomain: domain.NewProjectAuthDomain(),
	}
}

func (a *AccountService) Account(ctx context.Context, msg *account.AccountReqMessage) (*account.AccountResponse, error) {
	accountList, total, err := a.accountDomain.AccountList(
		msg.OrganizationCode,
		msg.MemberId,
		msg.Page,
		msg.PageSize,
		msg.DepartmentCode,
		msg.SearchType)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	authList, err := a.projectAuthDomain.AuthList(encrypts.DecryptNoErr(msg.OrganizationCode))
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var maList []*account.MemberAccount
	copier.Copy(&maList, accountList)
	var prList []*account.ProjectAuth
	copier.Copy(&prList, authList)
	return &account.AccountResponse{
		AccountList: maList,
		AuthList:    prList,
		Total:       total,
	}, nil
}

func (a *AccountService) SaveAccount(ctx context.Context, msg *account.AccountReqMessage) (*account.AccountResponse, error) {
	// forbid
	if msg.AccountCode != "" && msg.Authorize == "forbid" {
		id := encrypts.DecryptNoErr(msg.AccountCode)
		if err := a.accountDomain.ForbidAccount(ctx, id); err != nil {
			return nil, errs.GrpcError(err)
		}
		return &account.AccountResponse{}, nil
	}
	// resume
	if msg.AccountCode != "" && msg.Authorize == "resume" {
		id := encrypts.DecryptNoErr(msg.AccountCode)
		if err := a.accountDomain.ResumeAccount(ctx, id); err != nil {
			return nil, errs.GrpcError(err)
		}
		return &account.AccountResponse{}, nil
	}
	// del
	if msg.AccountCode != "" && msg.Authorize == "del" {
		id := encrypts.DecryptNoErr(msg.AccountCode)
		if err := a.accountDomain.DeleteAccount(ctx, id); err != nil {
			return nil, errs.GrpcError(err)
		}
		return &account.AccountResponse{}, nil
	}
	// auth
	if msg.AccountCode != "" && msg.Authorize != "" {
		id := encrypts.DecryptNoErr(msg.AccountCode)
		if err := a.accountDomain.AuthAccount(ctx, id, msg.Authorize); err != nil {
			return nil, errs.GrpcError(err)
		}
		return &account.AccountResponse{}, nil
	}
	// read (single account)
	if msg.AccountCode != "" {
		id := encrypts.DecryptNoErr(msg.AccountCode)
		ma, err := a.accountDomain.FindAccountById(id)
		if err != nil {
			return nil, errs.GrpcError(err)
		}
		var maMsg account.MemberAccount
		if ma != nil {
			copier.Copy(&maMsg, ma.ToDisplay())
		}
		return &account.AccountResponse{AccountList: []*account.MemberAccount{&maMsg}}, nil
	}
	// _allList
	if msg.SearchType == -1 {
		list, err := a.accountDomain.AllList(ctx, msg.OrganizationCode)
		if err != nil {
			return nil, errs.GrpcError(err)
		}
		var maList []*account.MemberAccount
		copier.Copy(&maList, list)
		return &account.AccountResponse{AccountList: maList, Total: int64(len(maList))}, nil
	}
	// add / edit
	orgId := encrypts.DecryptNoErr(msg.OrganizationCode)
	deptId := encrypts.DecryptNoErr(msg.DepartmentCode)
	authorize := msg.Authorize
	if msg.AccountCode != "" {
		id := encrypts.DecryptNoErr(msg.AccountCode)
		fields := make(map[string]interface{})
		if msg.Name != "" {
			fields["name"] = msg.Name
		}
		if msg.Mobile != "" {
			fields["mobile"] = msg.Mobile
		}
		if msg.Email != "" {
			fields["email"] = msg.Email
		}
		if msg.Avatar != "" {
			fields["avatar"] = msg.Avatar
		}
		if msg.Position != "" {
			fields["position"] = msg.Position
		}
		if msg.Authorize != "" {
			fields["authorize"] = msg.Authorize
		}
		if msg.DepartmentCode != "" {
			fields["department_code"] = deptId
		}
		if len(fields) > 0 {
			if err := a.accountDomain.UpdateAccountFields(ctx, id, fields); err != nil {
				return nil, errs.GrpcError(err)
			}
			interceptor.ClearAccountCache()
		}
		return &account.AccountResponse{}, nil
	}
	if authorize == "" && orgId > 0 {
		defaultAuth, err := a.projectAuthDomain.FindDefaultAuth(orgId)
		if err == nil && defaultAuth != nil {
			authorize = strconv.FormatInt(defaultAuth.Id, 10)
		}
	}
	ma := &data.MemberAccount{
		OrganizationCode: orgId,
		DepartmentCode:   deptId,
		Authorize:        authorize,
		Name:             msg.Name,
		Mobile:           msg.Mobile,
		Email:            msg.Email,
		Avatar:           msg.Avatar,
		Position:         msg.Position,
		Status:           1,
	}
	if err := a.accountDomain.SaveAccount(ctx, ma); err != nil {
		return nil, errs.GrpcError(err)
	}
	return &account.AccountResponse{}, nil
}
