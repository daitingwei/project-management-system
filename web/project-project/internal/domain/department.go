package domain

import (
	"context"
	"test.com/project-common/errs"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
	"time"
)

type DepartmentDomain struct {
	departmentRepo repo.DepartmentRepo
}

func (d *DepartmentDomain) FindDepartmentById(id int64) (*data.Department, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dp, err := d.departmentRepo.FindDepartmentById(c, id)
	if err != nil {
		return nil, model.DBError
	}
	return dp, nil
}

func (d *DepartmentDomain) List(organizationCode int64, parentDepartmentCode int64, page int64, size int64) ([]*data.DepartmentDisplay, int64, *errs.BError) {
	list, total, err := d.departmentRepo.ListDepartment(organizationCode, parentDepartmentCode, page, size)
	if err != nil {
		return nil, 0, model.DBError
	}
	var dList []*data.DepartmentDisplay
	for _, v := range list {
		dList = append(dList, v.ToDisplay())
	}
	return dList, total, nil
}

func (d *DepartmentDomain) Save(
	organizationCode int64,
	departmentCode int64,
	parentDepartmentCode int64,
	name string) (*data.DepartmentDisplay, *errs.BError) {

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dpm, err := d.departmentRepo.FindDepartment(c, organizationCode, parentDepartmentCode, name)
	if err != nil {
		return nil, model.DBError
	}
	if dpm == nil {
		dpm = &data.Department{
			Name:             name,
			OrganizationCode: organizationCode,
			CreateTime:       time.Now().UnixMilli(),
		}
		if parentDepartmentCode > 0 {
			dpm.Pcode = parentDepartmentCode
		}
		err := d.departmentRepo.Save(dpm)
		if err != nil {
			return nil, model.DBError
		}
		return dpm.ToDisplay(), nil
	}
	return dpm.ToDisplay(), nil
}

// TODO: Update 更新部门信息
// func (d *DepartmentDomain) Update(
// 	departmentCode int64,
// 	name string,
// 	parentDepartmentCode int64,
// 	icon string,
// 	sort int) (*data.DepartmentDisplay, *errs.BError) {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	dpm, err := d.departmentRepo.FindDepartmentById(c, departmentCode)
// 	if err != nil {
// 		return nil, model.DBError
// 	}
// 	if dpm == nil {
// 		return nil, model.DepartmentNotExist
// 	}

// 	if name != "" {
// 		dpm.Name = name
// 	}
// 	if parentDepartmentCode >= 0 {
// 		dpm.Pcode = parentDepartmentCode
// 	}
// 	if icon != "" {
// 		dpm.Icon = icon
// 	}
// 	if sort > 0 {
// 		dpm.Sort = sort
// 	}

// 	err = d.departmentRepo.Update(dpm)
// 	if err != nil {
// 		return nil, model.DBError
// 	}
// 	return dpm.ToDisplay(), nil
// }

// TODO: Delete 删除部门
// func (d *DepartmentDomain) Delete(departmentCode int64) *errs.BError {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	dpm, err := d.departmentRepo.FindDepartmentById(c, departmentCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	if dpm == nil {
// 		return model.DepartmentNotExist
// 	}

// 	children, err := d.departmentRepo.FindChildDepartments(departmentCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	if len(children) > 0 {
// 		return &errs.BError{Code: 400, Message: "该部门下存在子部门，无法删除"}
// 	}

// 	members, err := d.departmentRepo.FindDepartmentMember(c, departmentCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	if len(members) > 0 {
// 		return &errs.BError{Code: 400, Message: "该部门下存在成员，无法删除"}
// 	}

// 	err = d.departmentRepo.Delete(departmentCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	return nil
// }

// TODO: AddMember 添加部门成员
// func (d *DepartmentDomain) AddMember(
// 	organizationCode int64,
// 	departmentCode int64,
// 	accountCode int64,
// 	isOwner int) (*data.DepartmentMemberDisplay, *errs.BError) {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	dpm, err := d.departmentRepo.FindDepartmentById(c, departmentCode)
// 	if err != nil {
// 		return nil, model.DBError
// 	}
// 	if dpm == nil {
// 		return nil, model.DepartmentNotExist
// 	}

// 	existingMember, err := d.departmentRepo.FindDepartmentMemberByAccount(c, departmentCode, accountCode)
// 	if err != nil {
// 		return nil, model.DBError
// 	}
// 	if existingMember != nil {
// 		return nil, &errs.BError{Code: 400, Message: "该成员已在部门中"}
// 	}

// 	member := &data.DepartmentMember{
// 		DepartmentCode:   departmentCode,
// 		OrganizationCode: organizationCode,
// 		AccountCode:      accountCode,
// 		JoinTime:         time.Now().UnixMilli(),
// 		IsOwner:          isOwner,
// 	}

// 	err = d.departmentRepo.SaveDepartmentMember(member)
// 	if err != nil {
// 		return nil, model.DBError
// 	}

// 	return member.ToDisplay(), nil
// }

// TODO: RemoveMember 移除部门成员
// func (d *DepartmentDomain) RemoveMember(departmentCode int64, accountCode int64) *errs.BError {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	dpm, err := d.departmentRepo.FindDepartmentById(c, departmentCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	if dpm == nil {
// 		return model.DepartmentNotExist
// 	}

// 	member, err := d.departmentRepo.FindDepartmentMemberByAccount(c, departmentCode, accountCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	if member == nil {
// 		return &errs.BError{Code: 400, Message: "该成员不在部门中"}
// 	}

// 	err = d.departmentRepo.DeleteDepartmentMember(departmentCode, accountCode)
// 	if err != nil {
// 		return model.DBError
// 	}

// 	return nil
// }

// TODO: ListMembers 获取部门成员列表
// func (d *DepartmentDomain) ListMembers(departmentCode int64) ([]*data.DepartmentMemberDisplay, *errs.BError) {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	members, err := d.departmentRepo.FindDepartmentMember(c, departmentCode)
// 	if err != nil {
// 		return nil, model.DBError
// 	}

// 	var result []*data.DepartmentMemberDisplay
// 	for _, m := range members {
// 		result = append(result, m.ToDisplay())
// 	}

// 	return result, nil
// }

// TODO: SetPrincipal 设置部门负责人
// func (d *DepartmentDomain) SetPrincipal(departmentCode int64, accountCode int64, isPrincipal int) *errs.BError {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	member, err := d.departmentRepo.FindDepartmentMemberByAccount(c, departmentCode, accountCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	if member == nil {
// 		return &errs.BError{Code: 400, Message: "该成员不在部门中"}
// 	}

// 	member.IsPrincipal = isPrincipal
// 	err = d.departmentRepo.UpdateDepartmentMember(member)
// 	if err != nil {
// 		return model.DBError
// 	}

// 	return nil
// }

// TODO: GetDepartmentTree 获取部门树形结构
// type DepartmentTreeNode struct {
// 	*data.DepartmentDisplay
// 	Children []*DepartmentTreeNode
// }

// func (d *DepartmentDomain) GetDepartmentTree(organizationCode int64) ([]*DepartmentTreeNode, *errs.BError) {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	allDepts, err := d.departmentRepo.FindAllDepartmentByOrg(organizationCode)
// 	if err != nil {
// 		return nil, model.DBError
// 	}

// 	deptMap := make(map[int64]*data.Department)
// 	for _, dept := range allDepts {
// 		deptMap[dept.Id] = dept
// 	}

// 	var buildTree func(parentCode int64) []*DepartmentTreeNode
// 	buildTree = func(parentCode int64) []*DepartmentTreeNode {
// 		var result []*DepartmentTreeNode
// 		for _, dept := range allDepts {
// 			if dept.Pcode == parentCode {
// 				node := &DepartmentTreeNode{
// 					DepartmentDisplay: dept.ToDisplay(),
// 					Children:          buildTree(dept.Id),
// 				}
// 				result = append(result, node)
// 			}
// 		}
// 		return result
// 	}

// 	tree := buildTree(0)
// 	return tree, nil
// }

// TODO: MoveDepartment 移动部门（调整层级）
// func (d *DepartmentDomain) MoveDepartment(departmentCode int64, newParentCode int64) *errs.BError {
// 	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()

// 	dpm, err := d.departmentRepo.FindDepartmentById(c, departmentCode)
// 	if err != nil {
// 		return model.DBError
// 	}
// 	if dpm == nil {
// 		return model.DepartmentNotExist
// 	}

// 	if newParentCode > 0 {
// 		parentDept, err := d.departmentRepo.FindDepartmentById(c, newParentCode)
// 		if err != nil {
// 			return model.DBError
// 		}
// 		if parentDept == nil {
// 			return &errs.BError{Code: 400, Message: "目标父部门不存在"}
// 		}
// 		if newParentCode == departmentCode {
// 			return &errs.BError{Code: 400, Message: "不能将部门移动到自身"}
// 		}

// 		isDescendant := func(deptCode int64, ancestorCode int64) bool {
// 			current := deptCode
// 			for current != 0 {
// 				if current == ancestorCode {
// 					return true
// 				}
// 				parent, _ := deptMap[current]
// 				if parent == nil {
// 					break
// 				}
// 				current = parent.Pcode
// 			}
// 			return false
// 		}

// 		deptMap := make(map[int64]*data.Department)
// 		allDepts, _ := d.departmentRepo.FindAllDepartmentByOrg(dpm.OrganizationCode)
// 		for _, dept := range allDepts {
// 			deptMap[dept.Id] = dept
// 		}

// 		if isDescendant(newParentCode, departmentCode) {
// 			return &errs.BError{Code: 400, Message: "不能将部门移动到其子部门"}
// 		}
// 	}

// 	dpm.Pcode = newParentCode
// 	err = d.departmentRepo.Update(dpm)
// 	if err != nil {
// 		return model.DBError
// 	}

// 	return nil
// }

func NewDepartmentDomain() *DepartmentDomain {
	return &DepartmentDomain{
		departmentRepo: dao.NewDepartmentDao(),
	}
}

func (d *DepartmentDomain) Delete(departmentCode int64) *errs.BError {
	if err := d.departmentRepo.Delete(departmentCode); err != nil {
		return model.DBError
	}
	return nil
}
