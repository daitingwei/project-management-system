# 任务列表

## Task 1: 检查前端任务排序和缓存逻辑
- [x] Task 1.1: 检查 `task.vue` 中的任务列表数据获取和缓存机制
- [x] Task 1.2: 检查 `taskDetail.vue` 关闭时是否触发了列表刷新
- [x] Task 1.3: 检查前端是否有全局缓存或状态管理导致数据不更新

## Task 2: 检查后端任务排序接口
- [x] Task 2.1: 检查 `TaskSort` 方法的实现逻辑
- [x] Task 2.2: 检查排序更新是否正确提交到数据库
- [x] Task 2.3: 验证排序逻辑与数据库操作的原子性

## Task 3: 修复缓存问题
- [x] Task 3.1: 已修复后端 bug：TaskSort 中比较加密字符串的问题
- [ ] Task 3.2: 待验证修复效果

## 已修复的 Bug
### Bug 1: TaskSort 比较错误
- **位置**: `project-project/pkg/service/task.service.v1/task_service.go` 第323行
- **问题**: `if msg.PreTaskCode == msg.NextTaskCode` 比较的是加密字符串，应该比较解密后的数字
- **修复**: 先解密再比较
