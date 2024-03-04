# 项目管理系统

一个基于微服务架构的项目管理系统，支持项目管理、任务管理、成员管理、权限管理等功能。

## 技术栈

### 后端
- **语言**：Go 1.19+
- **框架**：Gin
- **数据库**：MySQL
- **缓存**：Redis
- **消息队列**：Kafka
- **配置中心**：Nacos
- **服务发现**：Etcd
- **链路追踪**：Jaeger
- **对象存储**：MinIO

### 前端
- **框架**：Vue.js
- **构建工具**：Webpack
- **UI 组件**：Element UI

## 项目结构

```
第二版/
├── project-api/          # API 网关服务
├── project-common/       # 公共组件库
├── project-grpc/        # gRPC 服务定义
├── project-project/     # 项目管理服务
├── project-user/       # 用户服务
├── page/               # 前端页面
├── docker-compose.yaml   # Docker 编排文件
├── start-all.sh       # 启动脚本
├── stop-all.sh        # 停止脚本
└── all-services.yml    # 服务配置文件
```

## 功能模块

### 1. 项目管理
- 项目列表
- 项目详情
- 项目成员管理
- 项目模板
- 项目归档
- 项目统计分析

### 2. 任务管理
- 任务列表
- 任务详情
- 任务阶段管理
- 任务标签管理
- 任务成员管理
- 任务文件管理
- 任务工时管理
- 任务回收站

### 3. 成员管理
- 成员列表
- 成员详情
- 成员邀请
- 成员权限管理

### 4. 部门管理
- 部门列表
- 部门详情
- 部门成员管理

### 5. 权限管理
- 权限列表
- 权限详情
- 权限分配
- 权限模板

### 6. 用户管理
- 用户登录
- 用户注册
- 用户信息管理
- 密码重置

## 快速开始

### 前置要求
- Docker 和 Docker Compose
- Go 1.19+
- MySQL 5.7+
- Redis 6.0+
- Nacos 2.0+
- Kafka 2.0+

### 启动服务
```bash
# 启动所有服务
./start-all.sh

# 或者单独启动某个服务
docker-compose up -d project-api
docker-compose up -d project-project
docker-compose up -d project-user
```

### 停止服务
```bash
# 停止所有服务
./stop-all.sh

# 或者单独停止某个服务
docker-compose down
```

## 开发指南

### 添加新功能
1. 在对应的服务模块中添加 API 接口
2. 在 `project-grpc` 中定义 gRPC 接口
3. 重新生成 gRPC 代码：`cd project-grpc/task && protoc --go_out=. --go-grpc_opt=paths=source_relative *.proto`
4. 更新路由注册
5. 测试 API 接口

### 提交代码
1. 使用功能分支进行开发
2. 提交信息格式：`feat: 添加功能描述`
3. 提交前进行代码审查

## 配置说明

### 环境变量
主要配置文件：
- `project-api/config/config.yaml` - API 服务配置
- `project-project/config/config.yaml` - 项目服务配置
- `project-user/config/config.yaml` - 用户服务配置
- `docker-compose.yaml` - Docker 编排配置

### Nacos 配置
- 配置文件：`project-api/config/bootstrap.yaml`
- 配置中心：`http://localhost:8848/nacos`
- 命名空间：`project-management`

### 数据库
- 主机：`localhost`
- 端口：`3306`
- 数据库：`project_management`

## 部署

### Docker 部署
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 生产环境部署
1. 配置环境变量
2. 修改数据库连接信息
3. 使用 Docker Compose 编排
4. 配置 Nacos 地址
5. 启动服务

## 许可证

Apache License 2.0
