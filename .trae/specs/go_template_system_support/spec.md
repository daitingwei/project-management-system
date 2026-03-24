# Go项目模板系统模板支持规格说明

## 问题
当前Go版本的项目模板查询逻辑与Java不一致：
- **Java设计**：`is_system=1` 表示系统模板，对所有组织可见，查询时使用 `organization_code = 当前组织 OR is_system = 1`
- **Go当前设计**：`FindProjectTemplateAll` 只按 `organization_code` 查询，没有包含系统模板

根据数据库查询结果，现有系统模板（is_system=1）的 organization_code 被错误设置为 17，这不符合系统模板的设计。

## 修改内容
- 修改 `FindProjectTemplateAll` 方法，使其查询逻辑与Java一致：返回**组织私有模板 + 系统模板**
- 确保系统模板（is_system=1）对所有组织可见

## 影响范围
- 关键文件：
  - `project-project/internal/dao/project_template.go` - 修改 `FindProjectTemplateAll` 查询逻辑
  - `project-project/internal/repo/project.go` - 接口层
  - `project-project/pkg/service/project.service.v1/project_service.go` - 服务层

## 需求说明
### 需求：项目模板查询支持系统模板
系统必须支持两种类型的模板：
1. **系统模板**（is_system=1）：对所有组织可见
2. **组织私有模板**（is_system=0）：只属于特定组织

#### 场景：查询所有模板（viewType=-1）
- **当** 用户请求查看所有项目模板（viewType=-1）
- **则** 返回当前组织的私有模板 + 系统模板（is_system=1）

#### 场景：查询系统模板（viewType=1）
- **当** 用户请求只查看系统模板（viewType=1）
- **则** 只返回 is_system=1 的模板

#### 场景：查询组织私有模板（viewType=0）
- **当** 用户请求只查看组织私有模板（viewType=0）
- **则** 只返回当前组织的私有模板（is_system=0 且 organization_code=当前组织）
