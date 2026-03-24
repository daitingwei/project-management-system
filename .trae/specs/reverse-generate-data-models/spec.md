# 反向生成 Data 数据模型文件规格说明

## Why
由于 `/Users/daitingwei/Desktop/第二版/web` 目录下 `project-user` 和 `project-project` 的 `internal/data` 文件夹丢失，需要根据现有的 proto 文件、API 模型、DAO 层代码和数据库使用情况反向生成数据模型文件。

## What Changes
- 创建 `project-user/internal/data/member/member.go` - 会员数据模型
- 创建 `project-user/internal/data/organization/organization.go` - 组织数据模型
- 创建 `project-project/internal/data/` 目录下的所有数据模型文件：
  - `project.go` - 项目模型
  - `project_member.go` - 项目成员模型
  - `project_collection.go` - 项目收藏模型
  - `project_template.go` - 项目模板模型
  - `task.go` - 任务模型
  - `task_member.go` - 任务成员模型
  - `task_stages.go` - 任务阶段模型
  - `task_stages_template.go` - 任务阶段模板模型
  - `task_tag.go` - 任务标签模型
  - `task_work_time.go` - 任务工时模型
  - `task_log.go` - 任务日志模型
  - `department.go` - 部门模型
  - `member_account.go` - 会员账户模型
  - `menu.go` - 菜单模型
  - `project_auth.go` - 项目权限模型
  - `project_auth_node.go` - 项目权限节点模型
  - `project_node.go` - 项目节点模型
  - `project_log.go` - 项目日志模型
  - `file.go` - 文件模型
  - `source_link.go` - 源链接模型

## Impact
- **Affected specs**: 数据访问层、领域层、服务层
- **Affected code**:
  - `project-user/internal/dao/` - DAO 层依赖数据模型
  - `project-user/internal/repo/` - Repository 层依赖数据模型
  - `project-user/pkg/service/` - 服务层依赖数据模型
  - `project-project/internal/dao/` - DAO 层依赖数据模型
  - `project-project/internal/repo/` - Repository 层依赖数据模型
  - `project-project/pkg/service/` - 服务层依赖数据模型

## ADDED Requirements

### Requirement: 数据模型结构定义
系统 SHALL 提供完整的数据模型定义，包括：
1. GORM 模型结构体定义
2. 表名映射
3. 字段标签（数据库字段、JSON 字段）
4. 时间字段格式化
5. 显示模型转换方法（ToDisplay）

#### Scenario: 会员数据模型
- **WHEN** 创建会员数据模型时
- **THEN** 应包含字段：ID、Account、Password、Name、Mobile、Email、CreateTime、LastLoginTime、Status、Realname、Address、Province、City、Area、Code、OrganizationCode、Avatar
- **AND** 表名应为 `ms_member`

#### Scenario: 组织数据模型
- **WHEN** 创建组织数据模型时
- **THEN** 应包含字段：ID、Name、Avatar、Description、MemberId、CreateTime、Personal、Address、Province、City、Area、Code、OwnerCode
- **AND** 表名应为 `ms_organization`

#### Scenario: 项目数据模型
- **WHEN** 创建项目数据模型时
- **THEN** 应包含字段：ID、Cover、Name、Description、AccessControlType、WhiteList、Order、Deleted、TemplateCode、Schedule、CreateTime、OrganizationCode、DeletedTime、Private、Prefix、OpenPrefix、Archive、ArchiveTime、OpenBeginTime、OpenTaskPrivate、TaskBoardTheme、BeginTime、EndTime、AutoUpdateSchedule
- **AND** 表名应为 `ms_project`

#### Scenario: 任务数据模型
- **WHEN** 创建任务数据模型时
- **THEN** 应包含字段：ID、ProjectCode、Name、Pri、ExecuteStatus、Description、CreateBy、DoneBy、DoneTime、CreateTime、AssignTo、Deleted、StageCode、TaskTag、Done、BeginTime、EndTime、RemindTime、Pcode、Sort、Like、Star、DeletedTime、Private、IdNum、Path、Schedule、VersionCode、FeaturesCode、WorkTime、Status、Code
- **AND** 表名应为 `ms_task`

### Requirement: 数据模型关联关系
系统 SHALL 正确定义数据模型之间的关联关系：
1. 一对多关系（如：项目 -> 任务）
2. 多对多关系（如：项目 <-> 成员）
3. 外键关联

#### Scenario: 项目成员关联
- **WHEN** 查询项目成员时
- **THEN** 应能通过 ProjectCode 和 MemberCode 关联项目和会员

#### Scenario: 任务成员关联
- **WHEN** 查询任务成员时
- **THEN** 应能通过 TaskCode 和 MemberCode 关联任务和会员

### Requirement: 数据模型方法
系统 SHALL 为数据模型提供必要的转换方法：
1. `ToDisplay()` - 转换为展示模型
2. `TableName()` - 返回表名
3. 时间戳格式化方法

#### Scenario: 时间字段格式化
- **WHEN** 调用 ToDisplay 方法时
- **THEN** 时间戳字段应转换为可读的日期时间字符串

## MODIFIED Requirements
无

## REMOVED Requirements
无

## 数据表清单

根据代码分析，需要创建以下数据表对应的模型：

### project-user 服务
| 表名 | 模型文件 | 说明 |
|------|---------|------|
| ms_member | member/member.go | 会员表 |
| ms_organization | organization/organization.go | 组织表 |

### project-project 服务
| 表名 | 模型文件 | 说明 |
|------|---------|------|
| ms_project | project.go | 项目表 |
| ms_project_member | project_member.go | 项目成员表 |
| ms_project_collection | project_collection.go | 项目收藏表 |
| ms_project_template | project_template.go | 项目模板表 |
| ms_task | task.go | 任务表 |
| ms_task_member | task_member.go | 任务成员表 |
| ms_task_stages | task_stages.go | 任务阶段表 |
| ms_task_stages_template | task_stages_template.go | 任务阶段模板表 |
| ms_task_tag | task_tag.go | 任务标签表 |
| ms_task_work_time | task_work_time.go | 任务工时表 |
| ms_task_log | task_log.go | 任务日志表 |
| ms_department | department.go | 部门表 |
| ms_member_account | member_account.go | 会员账户表 |
| ms_menu | menu.go | 菜单表 |
| ms_project_auth | project_auth.go | 项目权限表 |
| ms_project_auth_node | project_auth_node.go | 项目权限节点表 |
| ms_project_node | project_node.go | 项目节点表 |
| ms_project_log | project_log.go | 项目日志表 |
| ms_file | file.go | 文件表 |
| ms_source_link | source_link.go | 源链接表 |

## 字段映射规则

### Proto 类型到 Go 类型映射
- `int64` -> `int64`
- `int32` -> `int32`
- `string` -> `string`
- `double` -> `float64`
- `bool` -> `bool`
- `repeated` -> `[]Type`

### 数据库字段命名
- 使用下划线命名法（snake_case）
- 例如：`CreateTime` -> `create_time`

### JSON 字段命名
- 使用下划线命名法（snake_case）
- 例如：`CreateTime` -> `create_time`

### GORM 标签
- `gorm:"column:字段名"`
- `gorm:"type:字段类型"`
- `gorm:"primarykey"` - 主键
- `gorm:"foreignKey:外键字段"` - 外键

## 实现优先级

1. **高优先级**（核心业务模型）：
   - member.go
   - organization.go
   - project.go
   - project_member.go
   - task.go
   - task_stages.go

2. **中优先级**（功能支撑模型）：
   - department.go
   - member_account.go
   - menu.go
   - project_auth.go
   - project_node.go

3. **低优先级**（辅助功能模型）：
   - project_collection.go
   - project_template.go
   - task_member.go
   - task_tag.go
   - task_work_time.go
   - task_log.go
   - project_log.go
   - file.go
   - source_link.go
