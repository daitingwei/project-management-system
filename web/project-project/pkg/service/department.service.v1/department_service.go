package department_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/kk"
	"test.com/project-grpc/department"
	"test.com/project-grpc/user/login"
	"test.com/project-project/config"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/domain"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/rpc"
	"test.com/project-project/pkg/model"
	"time"
)

type DepartmentService struct {
	department.UnimplementedDepartmentServiceServer
	cache            repo.Cache
	transaction      tran.Transaction
	departmentDomain *domain.DepartmentDomain
	departmentRepo   repo.DepartmentRepo
}

func New() *DepartmentService {
	return &DepartmentService{
		cache:            dao.Rc,
		transaction:      dao.NewTransaction(),
		departmentDomain: domain.NewDepartmentDomain(),
		departmentRepo:   dao.NewDepartmentDao(),
	}
}
func (d *DepartmentService) List(ctx context.Context, msg *department.DepartmentReqMessage) (*department.ListDepartmentMessage, error) {
	organizationCode := encrypts.DecryptNoErr(msg.OrganizationCode)
	var parentDepartmentCode int64
	if msg.ParentDepartmentCode != "" {
		parentDepartmentCode = encrypts.DecryptNoErr(msg.ParentDepartmentCode)
	}
	dps, total, err := d.departmentDomain.List(
		organizationCode,
		parentDepartmentCode,
		msg.Page,
		msg.PageSize)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var list []*department.DepartmentMessage
	copier.Copy(&list, dps)
	config.SendLog(kk.Info("List", "DepartmentService.List", kk.FieldMap{
		"organizationCode":     organizationCode,
		"parentDepartmentCode": parentDepartmentCode,
		"page":                msg.Page,
	}))
	return &department.ListDepartmentMessage{List: list, Total: total}, nil
}

func (d *DepartmentService) Save(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
	organizationCode := encrypts.DecryptNoErr(msg.OrganizationCode)
	var departmentCode int64
	if msg.DepartmentCode != "" {
		departmentCode = encrypts.DecryptNoErr(msg.DepartmentCode)
	}
	var parentDepartmentCode int64
	if msg.ParentDepartmentCode != "" {
		parentDepartmentCode = encrypts.DecryptNoErr(msg.ParentDepartmentCode)
	}
	dp, err := d.departmentDomain.Save(
		organizationCode,
		departmentCode,
		parentDepartmentCode,
		msg.Name)
	if err != nil {
		return &department.DepartmentMessage{}, errs.GrpcError(err)
	}
	var res = &department.DepartmentMessage{}
	copier.Copy(res, dp)
	return res, nil
}

func (d *DepartmentService) Read(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
	dp, err := d.departmentDomain.FindDepartmentById(departmentCode)
	if err != nil {
		return &department.DepartmentMessage{}, errs.GrpcError(err)
	}
	var res = &department.DepartmentMessage{}
	copier.Copy(res, dp.ToDisplay())
	return res, nil
}

// TODO: Delete 删除部门
// func (d *DepartmentService) Delete(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
// 	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
// 	err := d.departmentDomain.Delete(departmentCode)
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}
// 	return &department.DepartmentMessage{}, nil
// }

func (d *DepartmentService) DeleteDepartment(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DeleteDepartmentResponse, error) {
	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
	if err := d.departmentDomain.Delete(departmentCode); err != nil {
		return nil, errs.GrpcError(err)
	}
	return &department.DeleteDepartmentResponse{}, nil
}

// TODO: Update 更新部门
// func (d *DepartmentService) Update(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
// 	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
// 	parentDepartmentCode := encrypts.DecryptNoErr(msg.ParentDepartmentCode)

// 	dp, err := d.departmentDomain.Update(
// 		departmentCode,
// 		msg.Name,
// 		parentDepartmentCode,
// 		msg.Icon,
// 		int(msg.Sort),
// 	)
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}

// 	var res = &department.DepartmentMessage{}
// 	copier.Copy(res, dp)
// 	return res, nil
// }

// TODO: AddDepartmentMember 添加部门成员
// func (d *DepartmentService) AddDepartmentMember(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
// 	organizationCode := encrypts.DecryptNoErr(msg.OrganizationCode)
// 	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
// 	accountCode := encrypts.DecryptNoErr(msg.AccountCode)

// 	member, err := d.departmentDomain.AddMember(
// 		organizationCode,
// 		departmentCode,
// 		accountCode,
// 		int(msg.IsOwner),
// 	)
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}

// 	var res = &department.DepartmentMessage{}
// 	copier.Copy(res, member)
// 	return res, nil
// }

// TODO: RemoveDepartmentMember 移除部门成员
// func (d *DepartmentService) RemoveDepartmentMember(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
// 	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
// 	accountCode := encrypts.DecryptNoErr(msg.AccountCode)

// 	err := d.departmentDomain.RemoveMember(departmentCode, accountCode)
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}

// 	return &department.DepartmentMessage{}, nil
// }

// TODO: ListDepartmentMembers 获取部门成员列表
// func (d *DepartmentService) ListDepartmentMembers(ctx context.Context, msg *department.DepartmentReqMessage) (*department.ListDepartmentMemberMessage, error) {
// 	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)

// 	members, err := d.departmentDomain.ListMembers(departmentCode)
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}

// 	var list []*department.DepartmentMemberMessage
// 	copier.Copy(&list, members)

// 	return &department.ListDepartmentMemberMessage{List: list}, nil
// }

// TODO: SetDepartmentPrincipal 设置部门负责人
// func (d *DepartmentService) SetDepartmentPrincipal(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
// 	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
// 	accountCode := encrypts.DecryptNoErr(msg.AccountCode)

// 	err := d.departmentDomain.SetPrincipal(departmentCode, accountCode, int(msg.IsPrincipal))
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}

// 	return &department.DepartmentMessage{}, nil
// }

// TODO: GetDepartmentTree 获取部门树形结构
// func (d *DepartmentService) GetDepartmentTree(ctx context.Context, msg *department.DepartmentReqMessage) (*department.ListDepartmentMessage, error) {
// 	organizationCode := encrypts.DecryptNoErr(msg.OrganizationCode)

// 	tree, err := d.departmentDomain.GetDepartmentTree(organizationCode)
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}

// 	var list []*department.DepartmentMessage
// 	for _, node := range tree {
// 		var dm department.DepartmentMessage
// 		copier.Copy(&dm, node.DepartmentDisplay)
// 		list = append(list, &dm)
// 	}

// 	return &department.ListDepartmentMessage{List: list, Total: int64(len(list))}, nil
// }

// TODO: MoveDepartment 移动部门
// func (d *DepartmentService) MoveDepartment(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMessage, error) {
// 	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
// 	parentDepartmentCode := encrypts.DecryptNoErr(msg.ParentDepartmentCode)

// 	err := d.departmentDomain.MoveDepartment(departmentCode, parentDepartmentCode)
// 	if err != nil {
// 		return nil, errs.GrpcError(err)
// 	}

// 	return &department.DepartmentMessage{}, nil
// }

func (d *DepartmentService) DepartmentMemberList(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMemberListResponse, error) {
	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
	zap.L().Info("DepartmentMemberList", zap.Int64("departmentCode", departmentCode))
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	members, err := d.departmentRepo.FindDepartmentMember(c, departmentCode)
	if err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	zap.L().Info("DepartmentMemberList members", zap.Int("count", len(members)))
	if len(members) == 0 {
		return &department.DepartmentMemberListResponse{List: []*department.DepartmentMemberMessage{}, Total: 0}, nil
	}

	var accountIds []int64
	for _, m := range members {
		accountIds = append(accountIds, m.AccountCode)
	}
	zap.L().Info("DepartmentMemberList accountIds", zap.Any("ids", accountIds))
	memberAccounts, err := d.departmentRepo.FindMemberAccountsByIds(c, accountIds)
	if err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	zap.L().Info("DepartmentMemberList memberAccounts", zap.Int("count", len(memberAccounts)))
	maMap := make(map[int64]*data.MemberAccount)
	var memberCodes []int64
	for _, ma := range memberAccounts {
		maMap[ma.Id] = ma
		memberCodes = append(memberCodes, ma.MemberCode)
	}

	memberInfoResp, err2 := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: memberCodes})
	if err2 != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	memberInfoMap := make(map[int64]*login.MemberMessage)
	for _, mi := range memberInfoResp.List {
		memberInfoMap[mi.Id] = mi
	}

	var list []*department.DepartmentMemberMessage
	for _, m := range members {
		ma, ok := maMap[m.AccountCode]
		if !ok {
			continue
		}
		dm := &department.DepartmentMemberMessage{
			Id:          m.AccountCode,
			Code:        encrypts.EncryptNoErr(m.AccountCode),
			AccountCode: encrypts.EncryptNoErr(m.AccountCode),
			IsPrincipal: int32(m.IsPrincipal),
		}
		if info, ok2 := memberInfoMap[ma.MemberCode]; ok2 {
			dm.Name = info.Name
			dm.Avatar = info.Avatar
			dm.Email = info.Email
		}
		list = append(list, dm)
	}
	zap.L().Info("DepartmentMemberList result", zap.Int("count", len(list)))
	return &department.DepartmentMemberListResponse{List: list, Total: int64(len(list))}, nil
}

func (d *DepartmentService) SearchDepartmentMember(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMemberListResponse, error) {
	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
	organizationCode := encrypts.DecryptNoErr(msg.OrganizationCode)
	keyword := msg.Keyword

	// keyword 为空时不搜索，直接返回空列表
	if keyword == "" {
		return &department.DepartmentMemberListResponse{List: []*department.DepartmentMemberMessage{}, Total: 0}, nil
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 用 map 去重，key = accountCode 加密串
	tempList := make(map[string]*department.DepartmentMemberMessage)

	if departmentCode > 0 {
		// ===== 添加至部门 =====
		// 查已加入该部门的 account_code 集合
		existingMembers, err := d.departmentRepo.FindDepartmentMember(c, departmentCode)
		if err != nil {
			return nil, errs.GrpcError(model.DBError)
		}
		joinedSet := make(map[int64]bool)
		for _, m := range existingMembers {
			joinedSet[m.AccountCode] = true
		}

		// 查组织内按 name like keyword 的 member_account
		orgMembers, err2 := d.departmentRepo.FindMemberAccountsByOrgWithKeyword(c, organizationCode, keyword)
		if err2 != nil {
			return nil, errs.GrpcError(model.DBError)
		}
		for _, ma := range orgMembers {
			// 添加至部门时 accountCode = ms_member_account.id 的加密串
			acCode := encrypts.EncryptNoErr(ma.Id)
			tempList[acCode] = &department.DepartmentMemberMessage{
				Id:          ma.Id,
				Code:        acCode,
				AccountCode: acCode,
				Name:        ma.Name,
				Avatar:      ma.Avatar,
				Email:       ma.Email,
				Joined:      joinedSet[ma.Id],
			}
		}
	} else {
		// ===== 添加至组织 =====
		// 查已加入组织的 member_code 集合
		orgMembers, err := d.departmentRepo.FindMemberAccountsByOrg(c, organizationCode)
		if err != nil {
			return nil, errs.GrpcError(model.DBError)
		}
		joinedMemberCodes := make(map[int64]bool)
		for _, ma := range orgMembers {
			joinedMemberCodes[ma.MemberCode] = true
		}

		// 查组织内按 name like keyword 的 member_account（status=1）
		matchedMembers, err2 := d.departmentRepo.FindMemberAccountsByOrgWithKeyword(c, organizationCode, keyword)
		if err2 != nil {
			return nil, errs.GrpcError(model.DBError)
		}
		for _, ma := range matchedMembers {
			// 添加至组织时 accountCode = ms_member.id（即 member_code）的加密串
			acCode := encrypts.EncryptNoErr(ma.MemberCode)
			tempList[acCode] = &department.DepartmentMemberMessage{
				Id:          ma.MemberCode,
				Code:        acCode,
				AccountCode: acCode,
				Name:        ma.Name,
				Avatar:      ma.Avatar,
				Email:       ma.Email,
				Joined:      true, // 已在组织内
			}
		}

		// 从 ms_member 全平台按 email 精确匹配，补充未加入组织的人
		if keyword != "" {
			member, err3 := d.departmentRepo.FindMemberByEmail(c, keyword)
			if err3 != nil {
				return nil, errs.GrpcError(model.DBError)
			}
			if member != nil {
				acCode := encrypts.EncryptNoErr(member.Id)
				if _, exists := tempList[acCode]; !exists {
					tempList[acCode] = &department.DepartmentMemberMessage{
						Id:          member.Id,
						Code:        acCode,
						AccountCode: acCode,
						Name:        member.Name,
						Avatar:      member.Avatar,
						Email:       member.Email,
						Joined:      joinedMemberCodes[member.Id],
					}
				}
			}
		}
	}

	var list []*department.DepartmentMemberMessage
	for _, item := range tempList {
		list = append(list, item)
	}
	if list == nil {
		list = []*department.DepartmentMemberMessage{}
	}
	return &department.DepartmentMemberListResponse{List: list, Total: int64(len(list))}, nil
}

func (d *DepartmentService) InviteDepartmentMember(ctx context.Context, msg *department.DepartmentReqMessage) (*department.InviteDepartmentMemberResponse, error) {
	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
	organizationCode := encrypts.DecryptNoErr(msg.OrganizationCode)
	// accountCode 在两种场景含义不同：
	// 有 departmentCode：accountCode = ms_member_account.id
	// 无 departmentCode：accountCode = ms_member.id（member_code）
	accountCode := encrypts.DecryptNoErr(msg.AccountCode)

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if departmentCode > 0 {
		// ===== 添加至部门 =====
		existing, err := d.departmentRepo.FindDepartmentMemberByAccount(c, departmentCode, accountCode)
		if err != nil {
			return nil, errs.GrpcError(model.DBError)
		}
		if existing != nil {
			return &department.InviteDepartmentMemberResponse{}, nil
		}
		member := &data.DepartmentMember{
			DepartmentCode:   departmentCode,
			OrganizationCode: organizationCode,
			AccountCode:      accountCode,
			JoinTime:         time.Now().UnixMilli(),
			IsOwner:          0,
			IsPrincipal:      0,
		}
		if err2 := d.departmentRepo.SaveDepartmentMember(member); err2 != nil {
			return nil, errs.GrpcError(model.DBError)
		}
		// 同步更新 ms_member_account.department_code（accountCode = ms_member_account.id）
		if err3 := d.departmentRepo.UpdateMemberAccountDepartment(c, accountCode, departmentCode); err3 != nil {
			return nil, errs.GrpcError(model.DBError)
		}
	} else {
		// ===== 添加至组织 =====
		// accountCode 此时是 ms_member.id（member_code）
		memberCode := accountCode
		// 检查是否已加入组织
		existing, err := d.departmentRepo.FindMemberAccountByMemberAndOrg(c, memberCode, organizationCode)
		if err != nil {
			return nil, errs.GrpcError(model.DBError)
		}
		if existing != nil {
			return &department.InviteDepartmentMemberResponse{}, nil
		}
		// 直接写 ms_member_account，name/email 从 rpc 获取
		memberInfoResp, err3 := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: []int64{memberCode}})
		if err3 != nil || len(memberInfoResp.List) == 0 {
			return nil, errs.GrpcError(model.DBError)
		}
		info := memberInfoResp.List[0]
		ma := &data.MemberAccount{
			MemberCode:       memberCode,
			OrganizationCode: organizationCode,
			IsOwner:          0,
			Name:             info.Name,
			Email:            info.Email,
			Avatar:           info.Avatar,
			Status:           1,
			CreateTime:       time.Now().UnixMilli(),
		}
		if err4 := d.departmentRepo.SaveMemberAccount(ma); err4 != nil {
			return nil, errs.GrpcError(model.DBError)
		}
	}
	return &department.InviteDepartmentMemberResponse{}, nil
}

func (d *DepartmentService) RemoveDepartmentMember(ctx context.Context, msg *department.DepartmentReqMessage) (*department.RemoveDepartmentMemberResponse, error) {
	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
	accountCode := encrypts.DecryptNoErr(msg.AccountCode)
	if err := d.departmentRepo.DeleteDepartmentMember(departmentCode, accountCode); err != nil {
		return nil, errs.GrpcError(model.DBError)
	}
	return &department.RemoveDepartmentMemberResponse{}, nil
}

func (d *DepartmentService) DepartmentMemberDetail(ctx context.Context, msg *department.DepartmentReqMessage) (*department.DepartmentMemberDetailResponse, error) {
	// code 是 accountCode 的加密串，organizationCode 用于校验
	accountCode := encrypts.DecryptNoErr(msg.AccountCode)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberInfoResp, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: []int64{accountCode}})
	if err != nil || len(memberInfoResp.List) == 0 {
		return &department.DepartmentMemberDetailResponse{}, nil
	}
	info := memberInfoResp.List[0]

	// 查询部门成员关系（可选，获取 isPrincipal）
	departmentCode := encrypts.DecryptNoErr(msg.DepartmentCode)
	var isPrincipal int32
	if departmentCode > 0 {
		m, _ := d.departmentRepo.FindDepartmentMemberByAccount(c, departmentCode, accountCode)
		if m != nil {
			isPrincipal = int32(m.IsPrincipal)
		}
	}
	return &department.DepartmentMemberDetailResponse{
		Member: &department.DepartmentMemberMessage{
			Id:          info.Id,
			Code:        info.Code,
			AccountCode: info.Code,
			Name:        info.Name,
			Avatar:      info.Avatar,
			Email:       info.Email,
			IsPrincipal: isPrincipal,
		},
	}, nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

func init() {
	_ = model.DepartmentNotExist
	_ = time.Now()
}
