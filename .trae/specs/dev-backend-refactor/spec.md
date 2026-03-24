# Dev-Backend 分支重置与细腻提交规范

## Why
用户希望重新整理 dev-backend 分支的提交历史：
- 原始版本在 `/Users/daitingwei/Desktop/第二版/后端/`
- 当前修改后的版本在 `/Users/daitingwei/Desktop/第二版/web/`
- 前端 page 保持一个大提交
- 后端按模块对比原始版本，细腻提交

## What Changes

### 1. .gitignore 配置
- 添加 `后端/` 到 .gitignore，确保原始版本不被 git 管理

### 2. 分支重置
- 重置 dev-backend 到初始状态（仅有 .gitignore 提交）
- 删除当前所有细腻提交（40个），保留 .gitignore

### 3. 前端一个大提交
- page/ 整个目录作为一个提交
- commit message: "feat(page): 添加 Vue 前端页面"

### 4. 后端按模块细腻提交
对比 `后端/` (原始) 和 `web/` (当前) 的区别，按模块提交：

#### project-common 公共库
- model.go, validate.go - 基础模型
- errs/ - 错误处理
- encrypts/ - 加密工具
- jwts/ - JWT 认证
- kk/ - Kafka 消息队列
- discovery/ - 服务发现
- nacos/ - Nacos 配置
- min/ - MinIO 存储
- tms/ - 链路追踪
- fs/ - 文件系统

#### project-grpc
- proto 文件定义
- 生成的 pb.go 文件

#### project-user
- config - 配置
- dao - 数据访问
- database - 数据库
- repo - 仓储
- interceptor - 拦截器
- pkg/model - 业务模型
- pkg/service - 服务
- router - 路由
- main.go - 入口

#### project-project
- config - 配置
- dao - 数据访问
- database - 数据库
- domain - 领域模型
- repo - 仓储
- interceptor - 拦截器
- pkg/model - 业务模型
- pkg/service - 服务层
- rpc - RPC 客户端
- router - 路由
- main.go - 入口

#### project-api
- config - 配置
- middleware - 中间件
- pkg/model - 数据模型
- api/ - 业务处理器
- router - 路由
- main.go - 入口

#### 部署配置
- docker-compose.yaml
- go.work
- 各种 shell 脚本

## Impact
- 分支: dev-backend 将完全重建
- 原有 40 个提交将被替换为新的提交结构

## Requirements

### Requirement: .gitignore 配置
系统 SHALL 确保 `后端/` 目录不被 git 追踪

#### Scenario: .gitignore 配置正确
- WHEN 检查 .gitignore 文件
- THEN 包含 `后端/` 规则

### Requirement: 前端单提交
系统 SHALL 将 page/ 目录作为单一提交

#### Scenario: 前端提交
- WHEN 执行前端提交
- THEN 仅包含 page/ 目录的变更

### Requirement: 后端模块化提交
系统 SHALL 按模块对比并提交后端代码

#### Scenario: 后端提交
- WHEN 提交后端代码
- THEN 每个模块独立提交，可清晰追溯
