# 任务列表

本文档记录从 `后端` 目录到 `web` 目录的代码演进任务，按照正常开发流程细致划分，并通过 Git 提交。

## 阶段一：基础设施增强

### Task 1: 新增 Nacos 配置中心支持
- [ ] SubTask 1.1: 在 project-common 中创建 nacos 包
- [ ] SubTask 1.2: 创建 nacos/bootstrap.go
- [ ] SubTask 1.3: 创建 nacos/nacos.go
- [ ] SubTask 1.4: 提交：feat(common): 添加 Nacos 配置中心支持

### Task 2: 新增 Kafka 消费者功能
- [ ] SubTask 2.1: 在 project-common/kk 中创建 consumer.go
- [ ] SubTask 2.2: 提交：feat(common): 添加 Kafka 消费者功能

### Task 3: 新增服务发现插件
- [ ] SubTask 3.1: 编译生成 discovery.so 插件文件
- [ ] SubTask 3.2: 提交：feat(common): 添加服务发现插件

## 阶段二：配置系统升级

### Task 4: 升级 project-api 配置系统
- [ ] SubTask 4.1: 在 config.go 中添加 MinioConfig 结构体定义
- [ ] SubTask 4.2: 在 config.go 中导入 nacos 包
- [ ] SubTask 4.3: 重构 InitConfig 函数，支持 Nacos 配置读取
- [ ] SubTask 4.4: 实现 ReLoadAllConfig 方法
- [ ] SubTask 4.5: 实现 ReadMinioConfig 方法
- [ ] SubTask 4.6: 创建 bootstrap.yaml 配置文件
- [ ] SubTask 4.7: 提交：feat(api): 升级配置系统支持 Nacos 和 Minio

## 阶段三：中间件重构

### Task 5: 创建独立的认证中间件包
- [ ] SubTask 5.1: 创建 project-api/middleware 目录
- [ ] SubTask 5.2: 创建 middleware/auth.go
- [ ] SubTask 5.3: 更新 router.go 使用 middleware.Auth()
- [ ] SubTask 5.4: 提交：refactor(api): 将认证中间件重构到独立包

## 阶段四：RPC 客户端重构

### Task 6: 重构 RPC 客户端调用方式
- [ ] SubTask 6.1: 创建 project-api/api/rpc/project_rpc.go
- [ ] SubTask 6.2: 将 RPC 客户端初始化逻辑迁移
- [ ] SubTask 6.3: 更新 task.go 中的 RPC 调用
- [ ] SubTask 6.4: 提交：refactor(api): 重构 RPC 客户端调用方式

## 阶段五：API 路由和处理器扩展

### Task 7: 新增 index.go 处理器
- [ ] SubTask 7.1: 创建 project-api/api/project/index.go
- [ ] SubTask 7.2: 实现 uploadAvatar 方法
- [ ] SubTask 7.3: 实现 uploadImg 方法
- [ ] SubTask 7.4: 实现 editPersonal 方法
- [ ] SubTask 7.5: 提交：feat(api): 添加个人中心 API

### Task 8: 新增 department_member.go 处理器
- [ ] SubTask 8.1: 创建 project-api/api/project/department_member.go
- [ ] SubTask 8.2: 实现部门成员相关方法
- [ ] SubTask 8.3: 提交：feat(api): 添加部门成员管理 API

### Task 9: 新增 task_tag.go 处理器
- [ ] SubTask 9.1: 创建 project-api/api/project/task_tag.go
- [ ] SubTask 9.2: 实现任务标签相关方法
- [ ] SubTask 9.3: 提交：feat(api): 添加任务标签管理 API

### Task 10: 扩展 task.go 处理器
- [ ] SubTask 10.1: 添加任务阶段管理方法
- [ ] SubTask 10.2: 添加文件管理方法
- [ ] SubTask 10.3: 添加任务回收站方法
- [ ] SubTask 10.4: 添加其他扩展方法
- [ ] SubTask 10.5: 提交：feat(api): 扩展任务管理 API

### Task 11: 扩展 project.go 处理器
- [ ] SubTask 11.1: 添加项目模板管理方法
- [ ] SubTask 11.2: 添加项目成员管理方法
- [ ] SubTask 11.3: 添加其他扩展方法
- [ ] SubTask 11.4: 提交：feat(api): 扩展项目管理 API

### Task 12: 扩展 account.go 处理器
- [ ] SubTask 12.1: 添加账号管理扩展方法
- [ ] SubTask 12.2: 提交：feat(api): 扩展账号管理 API

### Task 13: 扩展 department.go 处理器
- [ ] SubTask 13.1: 添加部门管理扩展方法
- [ ] SubTask 13.2: 提交：feat(api): 扩展部门管理 API

### Task 14: 扩展 auth.go 处理器
- [ ] SubTask 14.1: 添加权限管理扩展方法
- [ ] SubTask 14.2: 提交：feat(api): 扩展权限管理 API

### Task 15: 更新路由注册
- [ ] SubTask 15.1: 在 router.go 中注册所有新增路由
- [ ] SubTask 15.2: 提交：feat(api): 注册所有新增 API 路由

## 阶段六：业务逻辑层扩展

### Task 16: 扩展 task.go DAO
- [ ] SubTask 16.1: 添加任务 DAO 扩展方法
- [ ] SubTask 16.2: 提交：feat(project): 扩展任务 DAO

### Task 17: 新增 task_tag.go DAO
- [ ] SubTask 17.1: 创建 task_tag.go DAO
- [ ] SubTask 17.2: 提交：feat(project): 添加任务标签 DAO

### Task 18: 新增 task_tag 数据模型
- [ ] SubTask 18.1: 添加 TaskTag 和 TaskToTag 结构体
- [ ] SubTask 18.2: 提交：feat(project): 添加任务标签数据模型

### Task 19: 新增 task_tag 仓库
- [ ] SubTask 19.1: 创建 task_tag.go 仓库
- [ ] SubTask 19.2: 提交：feat(project): 添加任务标签仓库

### Task 20: 新增 task_tag_service 服务
- [ ] SubTask 20.1: 创建 task_tag_service.go
- [ ] SubTask 20.2: 提交：feat(project): 添加任务标签服务

## 阶段七：gRPC 接口更新

### Task 21: 更新 proto 定义文件
- [ ] SubTask 21.1: 更新 task_service.proto
- [ ] SubTask 21.2: 提交：feat(grpc): 更新任务服务 proto 定义

### Task 22: 重新生成 gRPC 代码
- [ ] SubTask 22.1: 运行 protoc 生成代码
- [ ] SubTask 22.2: 提交：feat(grpc): 重新生成 gRPC 代码

## 阶段八：其他文件更新

### Task 23: 更新 go.mod 依赖
- [ ] SubTask 23.1: 更新各模块 go.mod
- [ ] SubTask 23.2: 提交：chore: 更新依赖

### Task 24: 更新配置文件
- [ ] SubTask 24.1: 更新 config.yaml
- [ ] SubTask 24.2: 提交：chore: 更新配置文件

### Task 25: 新增工具脚本
- [ ] SubTask 25.1: 创建 start-all.sh
- [ ] SubTask 25.2: 创建 stop-all.sh
- [ ] SubTask 25.3: 提交：chore: 添加启动停止脚本

## 任务依赖关系

- Task 4 依赖于 Task 1
- Task 6 可以独立进行
- Task 7-14 依赖于 Task 6
- Task 16-20 可以与 Task 7-14 并行进行
- Task 21-22 应在服务实现之前完成
- Task 23-25 可以在最后进行
