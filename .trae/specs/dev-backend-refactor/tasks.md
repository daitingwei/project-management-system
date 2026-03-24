# Tasks - Dev-Backend 分支重置与细腻提交

## 任务列表

- [x] Task 1: 添加后端目录到 .gitignore
  - [x] SubTask 1.1: 添加 "后端/" 到 .gitignore
  - [x] SubTask 1.2: 提交 .gitignore 变更

- [ ] Task 2: 重置 dev-backend 分支
  - [ ] SubTask 2.1: 查看当前提交历史
  - [ ] SubTask 2.2: 重置到仅有 .gitignore 的状态
  - [ ] SubTask 2.3: 确认分支状态正确

- [ ] Task 3: 前端作为一个提交
  - [ ] SubTask 3.1: 添加 page/ 到暂存区
  - [ ] SubTask 3.2: 提交前端代码

- [ ] Task 4: 后端 - project-common 公共库 (按模块)
  - [ ] SubTask 4.1: 对比 model.go, validate.go 并提交
  - [ ] SubTask 4.2: 对比 errs/ 并提交
  - [ ] SubTask 4.3: 对比 encrypts/ 并提交
  - [ ] SubTask 4.4: 对比 jwts/ 并提交
  - [ ] SubTask 4.5: 对比 kk/ 并提交
  - [ ] SubTask 4.6: 对比 discovery/ 并提交
  - [ ] SubTask 4.7: 对比 nacos/, min/ 并提交
  - [ ] SubTask 4.8: 对比 tms/, fs/, run.go 并提交

- [ ] Task 5: 后端 - project-grpc 服务定义
  - [ ] SubTask 5.1: 对比 proto 文件并提交
  - [ ] SubTask 5.2: 对比 pb.go 生成文件并提交

- [ ] Task 6: 后端 - project-user 用户服务
  - [ ] SubTask 6.1: 对比 config 并提交
  - [ ] SubTask 6.2: 对比 dao, database 并提交
  - [ ] SubTask 6.3: 对比 repo, interceptor 并提交
  - [ ] SubTask 6.4: 对比 pkg, router, main.go 并提交

- [ ] Task 7: 后端 - project-project 项目服务
  - [ ] SubTask 7.1: 对比 config 并提交
  - [ ] SubTask 7.2: 对比 dao 并提交
  - [ ] SubTask 7.3: 对比 domain, database, interceptor 并提交
  - [ ] SubTask 7.4: 对比 repo 并提交
  - [ ] SubTask 7.5: 对比 pkg/service 并提交
  - [ ] SubTask 7.6: 对比 rpc, router, main.go 并提交

- [ ] Task 8: 后端 - project-api 网关
  - [ ] SubTask 8.1: 对比 config, middleware, tracing 并提交
  - [ ] SubTask 8.2: 对比 pkg/model 并提交
  - [ ] SubTask 8.3: 对比 api/, router 并提交
  - [ ] SubTask 8.4: 对比 main.go, go.mod, Dockerfile 并提交

- [ ] Task 9: 部署配置提交
  - [ ] SubTask 9.1: 对比 docker-compose.yaml 并提交
  - [ ] SubTask 9.2: 对比 go.work, 脚本文件并提交

- [ ] Task 10: 验证提交历史
  - [ ] SubTask 10.1: 检查 git log
  - [ ] SubTask 10.2: 确认未跟踪文件正确
  - [ ] SubTask 10.3: 对比 master 和 dev-backend 分支

## Task Dependencies
- Task 2 依赖 Task 1 完成
- Task 3 依赖 Task 2 完成
- Task 4-9 依赖 Task 3 完成
- Task 10 依赖 Task 4-9 全部完成
