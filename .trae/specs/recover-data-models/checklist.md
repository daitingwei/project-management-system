# Checklist

## 数据模型完整性检查

- [ ] project-user/internal/data/member/member.go 文件存在且包含完整的 Member 结构体定义
- [ ] project-user/internal/data/organization/organization.go 文件存在且包含完整的 Organization 结构体定义
- [ ] project-project/internal/data/project.go 文件存在且包含完整的 Project 结构体定义
- [ ] project-project/internal/data/project_member.go 文件存在且包含完整的 ProjectMember 结构体定义
- [ ] project-project/internal/data/project_collection.go 文件存在且包含完整的 ProjectCollection 结构体定义
- [ ] project-project/internal/data/project_template.go 文件存在且包含完整的 ProjectTemplate 结构体定义
- [ ] project-project/internal/data/task.go 文件存在且包含完整的 Task 结构体定义
- [ ] project-project/internal/data/task_member.go 文件存在且包含完整的 TaskMember 结构体定义
- [ ] project-project/internal/data/task_stages.go 文件存在且包含完整的 TaskStages 结构体定义
- [ ] project-project/internal/data/task_stages_template.go 文件存在且包含完整的 TaskStagesTemplate 结构体定义
- [ ] project-project/internal/data/task_tag.go 文件存在且包含完整的 TaskTag 结构体定义
- [ ] project-project/internal/data/task_work_time.go 文件存在且包含完整的 TaskWorkTime 结构体定义
- [ ] project-project/internal/data/task_log.go 文件存在且包含完整的 TaskLog 结构体定义
- [ ] project-project/internal/data/department.go 文件存在且包含完整的 Department 结构体定义
- [ ] project-project/internal/data/member_account.go 文件存在且包含完整的 MemberAccount 结构体定义
- [ ] project-project/internal/data/menu.go 文件存在且包含完整的 Menu 结构体定义
- [ ] project-project/internal/data/project_auth.go 文件存在且包含完整的 ProjectAuth 结构体定义
- [ ] project-project/internal/data/project_auth_node.go 文件存在且包含完整的 ProjectAuthNode 结构体定义
- [ ] project-project/internal/data/project_node.go 文件存在且包含完整的 ProjectNode 结构体定义
- [ ] project-project/internal/data/project_log.go 文件存在且包含完整的 ProjectLog 结构体定义
- [ ] project-project/internal/data/file.go 文件存在且包含完整的 File 结构体定义
- [ ] project-project/internal/data/source_link.go 文件存在且包含完整的 SourceLink 结构体定义

## GORM 标签检查

- [ ] 所有模型都包含 `gorm:"column:字段名"` 标签
- [ ] 主键字段包含 `gorm:"primarykey"` 标签
- [ ] 表名通过 TableName() 方法正确映射

## JSON 标签检查

- [ ] 所有模型字段都包含正确的 JSON 标签
- [ ] JSON 标签使用下划线命名法（snake_case）

## 方法实现检查

- [ ] Member 模型实现了必要的转换方法
- [ ] Organization 模型实现了必要的转换方法
- [ ] Project 模型实现了 ToDisplay() 方法
- [ ] Task 模型实现了 ToDisplay() 方法
- [ ] Department 模型实现了 ToDisplay() 方法
- [ ] MemberAccount 模型实现了 ToDisplay() 方法
- [ ] Menu 模型实现了 CovertChild() 方法

## 代码编译检查

- [ ] project-user 服务可以成功编译
- [ ] project-project 服务可以成功编译
- [ ] 没有未使用的导入
- [ ] 没有类型错误

## 字段完整性检查

### Member 模型字段
- [ ] ID 字段存在
- [ ] Account 字段存在
- [ ] Password 字段存在
- [ ] Name 字段存在
- [ ] Mobile 字段存在
- [ ] Email 字段存在
- [ ] CreateTime 字段存在
- [ ] LastLoginTime 字段存在
- [ ] Status 字段存在
- [ ] Realname 字段存在
- [ ] Address 字段存在
- [ ] Province 字段存在
- [ ] City 字段存在
- [ ] Area 字段存在
- [ ] Code 字段存在
- [ ] OrganizationCode 字段存在
- [ ] Avatar 字段存在

### Organization 模型字段
- [ ] ID 字段存在
- [ ] Name 字段存在
- [ ] Avatar 字段存在
- [ ] Description 字段存在
- [ ] MemberId 字段存在
- [ ] CreateTime 字段存在
- [ ] Personal 字段存在
- [ ] Address 字段存在
- [ ] Province 字段存在
- [ ] City 字段存在
- [ ] Area 字段存在
- [ ] Code 字段存在
- [ ] OwnerCode 字段存在

### Project 模型字段
- [ ] ID 字段存在
- [ ] Cover 字段存在
- [ ] Name 字段存在
- [ ] Description 字段存在
- [ ] AccessControlType 字段存在
- [ ] WhiteList 字段存在
- [ ] Order 字段存在
- [ ] Deleted 字段存在
- [ ] TemplateCode 字段存在
- [ ] Schedule 字段存在
- [ ] CreateTime 字段存在
- [ ] OrganizationCode 字段存在
- [ ] DeletedTime 字段存在
- [ ] Private 字段存在
- [ ] Prefix 字段存在
- [ ] OpenPrefix 字段存在
- [ ] Archive 字段存在
- [ ] ArchiveTime 字段存在
- [ ] OpenBeginTime 字段存在
- [ ] OpenTaskPrivate 字段存在
- [ ] TaskBoardTheme 字段存在
- [ ] BeginTime 字段存在
- [ ] EndTime 字段存在
- [ ] AutoUpdateSchedule 字段存在

### Task 模型字段
- [ ] ID 字段存在
- [ ] ProjectCode 字段存在
- [ ] Name 字段存在
- [ ] Pri 字段存在
- [ ] ExecuteStatus 字段存在
- [ ] Description 字段存在
- [ ] CreateBy 字段存在
- [ ] DoneBy 字段存在
- [ ] DoneTime 字段存在
- [ ] CreateTime 字段存在
- [ ] AssignTo 字段存在
- [ ] Deleted 字段存在
- [ ] StageCode 字段存在
- [ ] TaskTag 字段存在
- [ ] Done 字段存在
- [ ] BeginTime 字段存在
- [ ] EndTime 字段存在
- [ ] RemindTime 字段存在
- [ ] Pcode 字段存在
- [ ] Sort 字段存在
- [ ] Like 字段存在
- [ ] Star 字段存在
- [ ] DeletedTime 字段存在
- [ ] Private 字段存在
- [ ] IdNum 字段存在
- [ ] Path 字段存在
- [ ] Schedule 字段存在
- [ ] VersionCode 字段存在
- [ ] FeaturesCode 字段存在
- [ ] WorkTime 字段存在
- [ ] Status 字段存在
