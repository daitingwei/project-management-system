# Tasks

- [ ] Task 1: 重置到基准提交
  - 找到 develop 分支的基准提交（e737e96e）
  - 执行 `git reset --soft e737e96e` 保留所有更改在暂存区

- [ ] Task 2: 提交前端代码
  - 分离 page/ 目录的文件
  - 执行 `git commit -m "feat(frontend): 添加Vue前端项目"`
  - 验证提交成功

- [ ] Task 3: 提交后端代码
  - 分离 web/ 目录的文件
  - 执行 `git commit -m "feat(backend): 添加Go微服务后端"`
  - 验证提交成功

- [ ] Task 4: 提交配置和清理
  - 提交剩余文件（.gitignore 等）
  - 执行 `git commit -m "chore: 清理旧代码并更新配置"`
  - 验证提交成功

- [ ] Task 5: 验证最终结果
  - 执行 `git log --oneline` 确认提交历史
  - 确认只有 3 个清晰的提交
  - 确认工作区干净

# Task Dependencies
- Task 1 必须先完成
- Task 2, 3, 4 依赖 Task 1
- Task 5 依赖 Task 2, 3, 4
