# 完善数据库与测试前后端连接 Spec

## Why
当前系统数据库仅有部分测试数据，需要添加完整的测试数据才能进行功能验证。需要分析所有表的数据情况并补充必要数据。

## What Changes
- 分析所有20张表的数据情况
- 添加缺失的测试用户数据（当前仅1条）
- 添加组织架构数据（部门表为空）
- 添加项目相关数据（项目、成员、收藏等）
- 添加任务数据（任务、任务阶段、任务成员）
- 启动必要的基础服务（Redis、Etcd）
- 验证后端服务运行状态
- 验证前后端API连接

## Impact
- 系统具备完整测试数据
- 前后端可正常通信
- 可进行完整功能验证

## 当前数据库表数据统计

| 表名 | 当前数据量 | 状态 | 需要添加 |
|------|------------|------|----------|
| ms_member | 1 | ⚠️ 不足 | 是（需3-5条） |
| ms_member_account | 0 | ❌ 空 | 是 |
| ms_organization | 1 | ✅ 已有 | 否 |
| ms_department | 0 | ❌ 空 | **是（关键）** |
| ms_project | 0 | ❌ 空 | **是（关键）** |
| ms_project_member | 0 | ❌ 空 | **是（关键）** |
| ms_project_template | 3 | ✅ 已有 | 否 |
| ms_task_stages_template | 10 | ✅ 已有 | 否 |
| ms_project_node | 33 | ✅ 已有 | 否 |
| ms_project_auth | 0 | ❌ 空 | 是 |
| ms_project_auth_node | 0 | ❌ 空 | 是 |
| ms_project_collection | 0 | ❌ 空 | 是 |
| ms_project_log | 0 | ❌ 空 | 是 |
| ms_project_menu | 0 | ❌ 空 | 是 |
| ms_source_link | 0 | ❌ 空 | 是 |
| ms_task | 0 | ❌ 空 | **是（关键）** |
| ms_task_member | 0 | ❌ 空 | 是 |
| ms_task_stages | 0 | ❌ 空 | **是（关键）** |
| ms_task_work_time | 0 | ❌ 空 | 是 |
| ms_file | 0 | ❌ 空 | 是 |

## 阶段一：完善数据库测试数据

### Requirement: 添加测试用户数据
系统应包含至少3个测试用户用于功能验证

#### Scenario: 测试用户数据
- **GIVEN** 数据库msproject存在
- **WHEN** 插入测试用户数据到ms_member表
- **WHEN** 插入用户账号关联到ms_member_account表
- **THEN** 用户表至少包含3条记录
- **THEN** 账号关联表有对应记录

### Requirement: 添加组织架构数据
系统应包含公司和部门数据

#### Scenario: 组织架构数据
- **GIVEN** 用户数据存在
- **WHEN** 插入部门数据到ms_department表
- **THEN** 部门表有记录（当前为空）
- **THEN** 部门与组织关联正确

### Requirement: 添加项目数据
系统应包含测试项目及成员

#### Scenario: 项目数据
- **GIVEN** 组织和部门数据存在
- **WHEN** 插入项目数据到ms_project表
- **WHEN** 插入项目成员到ms_project_member表
- **WHEN** 插入项目收藏到ms_project_collection表
- **THEN** 项目表有记录
- **THEN** 项目成员表有记录

### Requirement: 添加任务数据
系统应包含测试任务及阶段

#### Scenario: 任务数据
- **GIVEN** 项目数据存在
- **WHEN** 插入任务阶段到ms_task_stages表
- **WHEN** 插入任务到ms_task表
- **WHEN** 插入任务成员到ms_task_member表
- **THEN** 任务阶段表有记录
- **THEN** 任务表有记录
- **THEN** 任务成员表有记录

### Requirement: 添加其他必要数据
系统应包含权限、菜单等辅助数据

#### Scenario: 辅助数据
- **WHEN** 插入项目权限到ms_project_auth表
- **WHEN** 插入权限节点到ms_project_auth_node表
- **WHEN** 插入项目菜单到ms_project_menu表
- **WHEN** 插入项目日志到ms_project_log表
- **THEN** 各辅助表有记录

## 阶段二：启动必要服务

### Requirement: Redis服务
后端服务需要Redis运行

#### Scenario: Redis连接
- **WHEN** Redis容器启动
- **THEN** 可以通过localhost:6379连接

### Requirement: Etcd服务
gRPC服务注册发现需要Etcd

#### Scenario: Etcd连接
- **WHEN** Etcd容器启动
- **THEN** 可以通过localhost:2379连接

## 阶段三：测试前后端连接

### Requirement: 后端API可访问
后端服务应正常运行并响应请求

#### Scenario: 后端服务运行
- **WHEN** 后端服务启动
- **THEN** 可以访问API端点

### Requirement: 前端可连接后端
前端应能成功调用后端API

#### Scenario: 前后端通信
- **WHEN** 前端发起请求
- **THEN** 收到正确的响应
