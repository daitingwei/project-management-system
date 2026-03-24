# 任务列表

## Task 1: 修改 FindProjectTemplateAll 查询逻辑
- [x] Task 1.1: 修改 `project_project/internal/dao/project_template.go` 中的 `FindProjectTemplateAll` 方法
  - 将查询条件从 `organization_code = ?` 修改为 `organization_code = ? OR is_system = 1`
  - 这样可以返回当前组织的私有模板 + 系统模板

## Task 2: 修正数据库中系统模板的 organization_code
- [x] Task 2.1: 将系统模板（is_system=1）的 organization_code 设置为 0
  - 已更新3条系统模板数据
