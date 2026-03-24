# Develop 分支 Git 历史重构 Spec

## Why
当前 develop 分支的提交历史不够清晰，需要按 Conventional Commits 规范重组，使其更易读、易追溯。

## What Changes
- 使用 `git reset --soft` 重置到基准提交
- 分三次提交，每次只包含特定模块的改动
- 使用规范的 commit message 格式

## Impact
- 只影响 Git 提交历史，不影响代码本身
- 分支：`develop`
- 提交数量：从 3 个重组成 3 个更清晰的提交

## 当前提交历史
```
c86bff83 feat: 添加前端Vue项目和后端Go微服务项目
bcbfd125 chore: 清理重复的Java框架代码
9999cfb6 feat: 隐藏未实现功能的前端入口
```

## 目标提交历史
```
<hash1> feat(frontend): 添加Vue前端项目
<hash2> feat(backend): 添加Go微服务后端
<hash3> chore: 清理旧代码并更新配置
```

## ADDED Requirements
### Requirement: 前端提交
**WHEN** 提交前端代码时
**THEN** 只包含 page/ 目录，使用 feat(frontend): 格式

### Requirement: 后端提交
**WHEN** 提交后端代码时
**THEN** 只包含 web/ 目录，使用 feat(backend): 格式

### Requirement: 配置提交
**WHEN** 提交配置和清理时
**THEN** 包含 .gitignore 等配置文件，使用 chore: 格式

## MODIFIED Requirements
无

## REMOVED Requirements
无
