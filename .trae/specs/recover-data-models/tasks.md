# Tasks

## 第一阶段：创建 project-user 数据模型

- [ ] Task 1: 创建 project-user/internal/data/member/member.go
  - [ ] SubTask 1.1: 定义 Member 结构体，包含所有字段（ID、Account、Password、Name、Mobile、Email、CreateTime、LastLoginTime、Status、Realname、Address、Province、City、Area、Code、OrganizationCode、Avatar）
  - [ ] SubTask 1.2: 实现 TableName() 方法，返回表名 ms_member
  - [ ] SubTask 1.3: 添加 GORM 标签和 JSON 标签
  - [ ] SubTask 1.4: 实现 ToDisplay() 方法（如果需要）

- [ ] Task 2: 创建 project-user/internal/data/organization/organization.go
  - [ ] SubTask 2.1: 定义 Organization 结构体，包含所有字段（ID、Name、Avatar、Description、MemberId、CreateTime、Personal、Address、Province、City、Area、Code、OwnerCode）
  - [ ] SubTask 2.2: 实现 TableName() 方法，返回表名 ms_organization
  - [ ] SubTask 2.3: 添加 GORM 标签和 JSON 标签
  - [ ] SubTask 2.4: 实现 ToDisplay() 方法（如果需要）

## 第二阶段：创建 project-project 核心数据模型

- [ ] Task 3: 创建 project-project/internal/data/project.go
  - [ ] SubTask 3.1: 定义 Project 结构体，包含所有字段
  - [ ] SubTask 3.2: 实现 TableName() 方法，返回表名 ms_project
  - [ ] SubTask 3.3: 添加 GORM 标签和 JSON 标签
  - [ ] SubTask 3.4: 实现 ToDisplay() 方法

- [ ] Task 4: 创建 project-project/internal/data/project_member.go
  - [ ] SubTask 4.1: 定义 ProjectMember 结构体
  - [ ] SubTask 4.2: 实现 TableName() 方法
  - [ ] SubTask 4.3: 添加 GORM 标签和 JSON 标签

- [ ] Task 5: 创建 project-project/internal/data/task.go
  - [ ] SubTask 5.1: 定义 Task 结构体，包含所有字段
  - [ ] SubTask 5.2: 实现 TableName() 方法，返回表名 ms_task
  - [ ] SubTask 5.3: 添加 GORM 标签和 JSON 标签
  - [ ] SubTask 5.4: 实现 ToDisplay() 方法

- [ ] Task 6: 创建 project-project/internal/data/task_stages.go
  - [ ] SubTask 6.1: 定义 TaskStages 结构体
  - [ ] SubTask 6.2: 实现 TableName() 方法
  - [ ] SubTask 6.3: 添加 GORM 标签和 JSON 标签

## 第三阶段：创建 project-project 功能支撑模型

- [ ] Task 7: 创建 project-project/internal/data/department.go
  - [ ] SubTask 7.1: 定义 Department 结构体
  - [ ] SubTask 7.2: 实现 TableName() 方法
  - [ ] SubTask 7.3: 添加 GORM 标签和 JSON 标签
  - [ ] SubTask 7.4: 实现 ToDisplay() 方法

- [ ] Task 8: 创建 project-project/internal/data/member_account.go
  - [ ] SubTask 8.1: 定义 MemberAccount 结构体
  - [ ] SubTask 8.2: 实现 TableName() 方法
  - [ ] SubTask 8.3: 添加 GORM 标签和 JSON 标签
  - [ ] SubTask 8.4: 实现 ToDisplay() 方法

- [ ] Task 9: 创建 project-project/internal/data/menu.go
  - [ ] SubTask 9.1: 定义 Menu 结构体
  - [ ] SubTask 9.2: 实现 TableName() 方法
  - [ ] SubTask 9.3: 添加 GORM 标签和 JSON 标签
  - [ ] SubTask 9.4: 实现 CovertChild() 方法（转换子菜单）

- [ ] Task 10: 创建 project-project/internal/data/project_auth.go
  - [ ] SubTask 10.1: 定义 ProjectAuth 结构体
  - [ ] SubTask 10.2: 实现 TableName() 方法
  - [ ] SubTask 10.3: 添加 GORM 标签和 JSON 标签

- [ ] Task 11: 创建 project-project/internal/data/project_node.go
  - [ ] SubTask 11.1: 定义 ProjectNode 结构体
  - [ ] SubTask 11.2: 实现 TableName() 方法
  - [ ] SubTask 11.3: 添加 GORM 标签和 JSON 标签

## 第四阶段：创建 project-project 辅助功能模型

- [ ] Task 12: 创建 project-project/internal/data/project_collection.go
  - [ ] SubTask 12.1: 定义 ProjectCollection 结构体
  - [ ] SubTask 12.2: 实现 TableName() 方法

- [ ] Task 13: 创建 project-project/internal/data/project_template.go
  - [ ] SubTask 13.1: 定义 ProjectTemplate 结构体
  - [ ] SubTask 13.2: 实现 TableName() 方法

- [ ] Task 14: 创建 project-project/internal/data/task_member.go
  - [ ] SubTask 14.1: 定义 TaskMember 结构体
  - [ ] SubTask 14.2: 实现 TableName() 方法

- [ ] Task 15: 创建 project-project/internal/data/task_stages_template.go
  - [ ] SubTask 15.1: 定义 TaskStagesTemplate 结构体
  - [ ] SubTask 15.2: 实现 TableName() 方法

- [ ] Task 16: 创建 project-project/internal/data/task_tag.go
  - [ ] SubTask 16.1: 定义 TaskTag 结构体
  - [ ] SubTask 16.2: 实现 TableName() 方法

- [ ] Task 17: 创建 project-project/internal/data/task_work_time.go
  - [ ] SubTask 17.1: 定义 TaskWorkTime 结构体
  - [ ] SubTask 17.2: 实现 TableName() 方法

- [ ] Task 18: 创建 project-project/internal/data/task_log.go
  - [ ] SubTask 18.1: 定义 TaskLog 结构体
  - [ ] SubTask 18.2: 实现 TableName() 方法

- [ ] Task 19: 创建 project-project/internal/data/project_auth_node.go
  - [ ] SubTask 19.1: 定义 ProjectAuthNode 结构体
  - [ ] SubTask 19.2: 实现 TableName() 方法

- [ ] Task 20: 创建 project-project/internal/data/project_log.go
  - [ ] SubTask 20.1: 定义 ProjectLog 结构体
  - [ ] SubTask 20.2: 实现 TableName() 方法

- [ ] Task 21: 创建 project-project/internal/data/file.go
  - [ ] SubTask 21.1: 定义 File 结构体
  - [ ] SubTask 21.2: 实现 TableName() 方法

- [ ] Task 22: 创建 project-project/internal/data/source_link.go
  - [ ] SubTask 22.1: 定义 SourceLink 结构体
  - [ ] SubTask 22.2: 实现 TableName() 方法

## 第五阶段：创建辅助方法和工具函数

- [ ] Task 23: 创建 project-project/internal/data/convert.go（如果需要）
  - [ ] SubTask 23.1: 实现 ToMap() 方法
  - [ ] SubTask 23.2: 实现 ToTaskStagesMap() 方法
  - [ ] SubTask 23.3: 实现其他辅助转换方法

## 第六阶段：验证和测试

- [ ] Task 24: 验证数据模型完整性
  - [ ] SubTask 24.1: 检查所有模型是否正确导入
  - [ ] SubTask 24.2: 检查所有字段是否完整
  - [ ] SubTask 24.3: 检查 GORM 标签是否正确

- [ ] Task 25: 编译验证
  - [ ] SubTask 25.1: 编译 project-user 服务
  - [ ] SubTask 25.2: 编译 project-project 服务
  - [ ] SubTask 25.3: 修复编译错误（如果有）

# Task Dependencies

- Task 1 和 Task 2 可以并行执行
- Task 3-6 可以并行执行（依赖 Task 1 完成）
- Task 7-11 可以并行执行（依赖 Task 3 完成）
- Task 12-22 可以并行执行（依赖 Task 7 完成）
- Task 23 依赖 Task 3-22 完成
- Task 24-25 依赖所有前置任务完成
