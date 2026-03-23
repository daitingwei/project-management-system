package dao

import (
	"context"
	"gorm.io/gorm"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type DepartmentDao struct {
	conn *gorms.GormConn
}

func (d *DepartmentDao) Save(dpm *data.Department) error {
	err := d.conn.Session(context.Background()).Save(&dpm).Error
	return err
}

func (d *DepartmentDao) FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*data.Department, error) {
	session := d.conn.Session(ctx)
	session = session.Model(&data.Department{}).Where("organization_code=? AND name=?", organizationCode, name)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	var dp *data.Department
	err := session.Limit(1).Take(&dp).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return dp, err
}

func (d *DepartmentDao) ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*data.Department, total int64, err error) {
	session := d.conn.Session(context.Background())
	session = session.Model(&data.Department{})
	session = session.Where("organization_code=?", organizationCode)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	err = session.Count(&total).Error
	err = session.Limit(int(size)).Offset(int((page - 1) * size)).Find(&list).Error
	return
}

func (d *DepartmentDao) FindDepartmentById(ctx context.Context, id int64) (dt *data.Department, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("id=?", id).Find(&dt).Error
	return
}

func (d *DepartmentDao) Update(dpm *data.Department) error {
	err := d.conn.Session(context.Background()).Save(dpm).Error
	return err
}

func (d *DepartmentDao) Delete(id int64) error {
	err := d.conn.Session(context.Background()).Where("id=?", id).Delete(&data.Department{}).Error
	return err
}

func (d *DepartmentDao) FindChildDepartments(parentCode int64) (list []*data.Department, err error) {
	session := d.conn.Session(context.Background())
	err = session.Where("pcode=?", parentCode).Find(&list).Error
	return
}

func (d *DepartmentDao) FindAllDepartmentByOrg(organizationCode int64) (list []*data.Department, err error) {
	session := d.conn.Session(context.Background())
	err = session.Where("organization_code=?", organizationCode).Find(&list).Error
	return
}

func (d *DepartmentDao) SaveDepartmentMember(dpm *data.DepartmentMember) error {
	err := d.conn.Session(context.Background()).Save(&dpm).Error
	return err
}

func (d *DepartmentDao) FindDepartmentMember(ctx context.Context, departmentCode int64) (list []*data.DepartmentMember, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("department_code=?", departmentCode).Find(&list).Error
	return
}

func (d *DepartmentDao) FindDepartmentMemberByAccount(ctx context.Context, departmentCode int64, accountCode int64) (*data.DepartmentMember, error) {
	session := d.conn.Session(ctx)
	var dm *data.DepartmentMember
	err := session.Where("department_code=? AND account_code=?", departmentCode, accountCode).Limit(1).Take(&dm).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return dm, err
}

func (d *DepartmentDao) DeleteDepartmentMember(departmentCode int64, accountCode int64) error {
	err := d.conn.Session(context.Background()).Where("department_code=? AND account_code=?", departmentCode, accountCode).Delete(&data.DepartmentMember{}).Error
	return err
}

func (d *DepartmentDao) UpdateDepartmentMember(dpm *data.DepartmentMember) error {
	err := d.conn.Session(context.Background()).Save(dpm).Error
	return err
}

func (d *DepartmentDao) FindDepartmentMembersByOrg(ctx context.Context, organizationCode int64) (list []*data.DepartmentMember, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("organization_code=?", organizationCode).Find(&list).Error
	return
}

func (d *DepartmentDao) FindMemberAccountsByOrg(ctx context.Context, organizationCode int64) (list []*data.MemberAccount, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("organization_code=? AND status=1", organizationCode).Find(&list).Error
	return
}

func (d *DepartmentDao) FindMemberAccountsByOrgWithKeyword(ctx context.Context, organizationCode int64, keyword string) (list []*data.MemberAccount, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("organization_code=? AND status=1 AND (name LIKE ? OR email LIKE ?)", organizationCode, "%"+keyword+"%", "%"+keyword+"%").Find(&list).Error
	return
}

func (d *DepartmentDao) FindMemberByEmail(ctx context.Context, email string) (*data.Member, error) {
	session := d.conn.Session(ctx)
	var m data.Member
	err := session.Where("email=?", email).Limit(1).Take(&m).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &m, err
}

func (d *DepartmentDao) FindMemberAccountByMemberAndOrg(ctx context.Context, memberCode int64, organizationCode int64) (*data.MemberAccount, error) {
	session := d.conn.Session(ctx)
	var ma data.MemberAccount
	err := session.Where("member_code=? AND organization_code=?", memberCode, organizationCode).Limit(1).Take(&ma).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ma, err
}

func (d *DepartmentDao) SaveMemberAccount(ma *data.MemberAccount) error {
	return d.conn.Session(context.Background()).Save(ma).Error
}

func (d *DepartmentDao) FindMemberAccountsByIds(ctx context.Context, ids []int64) (list []*data.MemberAccount, err error) {
	if len(ids) == 0 {
		return nil, nil
	}
	err = d.conn.Session(ctx).Where("member_code IN (?)", ids).Find(&list).Error
	return
}

func (d *DepartmentDao) UpdateMemberAccountDepartment(ctx context.Context, accountId int64, departmentCode int64) error {
	return d.conn.Session(ctx).Model(&data.MemberAccount{}).
		Where("id = ?", accountId).
		Update("department_code", departmentCode).Error
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{
		conn: gorms.New(),
	}
}
