package model

import (
	"test.com/project-common/errs"
)

var (
	RedisError            = errs.NewError(999, "redis错误")
	DBError               = errs.NewError(998, "db错误")
	ParamsError           = errs.NewError(401, "参数错误")
	NoLegalMobile         = errs.NewError(10102001, "手机号不合法")
	CaptchaNotExist       = errs.NewError(10102002, "验证码不存在或者已过期")
	CaptchaError          = errs.NewError(10102003, "验证码错误")
	EmailExist            = errs.NewError(10102004, "邮箱已经存在了")
	AccountExist          = errs.NewError(10102005, "账号已经存在了")
	MobileExist           = errs.NewError(10102006, "手机号已经存在了")
	AccountAndPwdError    = errs.NewError(10102007, "账号密码不正确")
	TaskNameNotNull       = errs.NewError(20102001, "任务标题不能为空")
	TaskStagesNotNull     = errs.NewError(20102002, "任务步骤不存在")
	TaskNotFound          = errs.NewError(20102003, "任务不存在")
	ProjectAlreadyDeleted = errs.NewError(20102004, "项目已经删除了")
	DepartmentNotExist   = errs.NewError(20103001, "部门不存在")
	DepartmentExist      = errs.NewError(20103002, "部门名称已存在")
	DepartmentHasChild   = errs.NewError(20103003, "该部门下存在子部门，无法删除")
	DepartmentHasMember  = errs.NewError(20103004, "该部门下存在成员，无法删除")
	MemberAlreadyInDepartment = errs.NewError(20103005, "该成员已在部门中")
	MemberNotInDepartment    = errs.NewError(20103006, "该成员不在部门中")
	TaskTagNameExist         = errs.NewError(20102005, "标签名称已存在")
	TemplateNameNotNull      = errs.NewError(20102006, "模板名称不能为空")
	TemplateCodeNotNull     = errs.NewError(20102007, "模板编码不能为空")
	TemplateNotExist        = errs.NewError(20102008, "模板不存在")
)
