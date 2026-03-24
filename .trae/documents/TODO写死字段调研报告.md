# 全项目 TODO 调研报告

## 一、Go版本 TODO 清单（共52处）

### 1. 部门模块（Department）- 8处 ❌ 待实现
| 文件 | 行号 | 内容 | 对比Java |
|------|------|------|----------|
| department_service.go | 91 | TODO: Delete 删除部门 | ✅ Java有实现 |
| department_service.go | 101 | TODO: Update 更新部门 | ✅ Java有实现 |
| department_service.go | 122 | TODO: AddDepartmentMember 添加部门成员 | ✅ Java有实现 |
| department_service.go | 143 | TODO: RemoveDepartmentMember 移除部门成员 | ✅ Java有实现 |
| department_service.go | 156 | TODO: ListDepartmentMembers 获取部门成员列表 | ✅ Java有实现 |
| department_service.go | 171 | TODO: SetDepartmentPrincipal 设置部门负责人 | ✅ Java有实现 |
| department_service.go | 184 | TODO: GetDepartmentTree 获取部门树形结构 | ✅ Java有实现 |
| department_service.go | 203 | TODO: MoveDepartment 移动部门 | ✅ Java有实现 |

### 2. 任务模块（Task）- 9处 ⭐ 与Java一致
| 文件 | 行号 | 内容 | 对比Java |
|------|------|------|----------|
| task_service.go | 280 | actionType写死为"task" | ✅ Java也是写死的 |
| task_service.go | 286 | 缓存key写死为"task" | ⚠️ Java可能简化了 |
| task_service.go | 658 | ObjectType写死为空 | ✅ Java也是空的 |
| task_service.go | 675 | LinkType写死为"task" | ✅ Java也是写死的 |
| task_service.go | 748 | 日志Type写死为"createComment" | ⚠️ 待确认 |
| task_service.go | 753 | ActionType写死为"task" | ✅ Java也是写死的 |
| source_link.go | 26 | 查询条件link_type写死 | ✅ Java也是写死的 |
| cache.go | 96 | 缓存hash key写死为"task" | ⚠️ Java可能简化了 |
| kafka.go | 45 | 缓存清除判断写死为"task" | ⚠️ Java可能简化了 |

### 3. 登录/用户模块（Login）- 3处 ⚠️ 待确认
| 文件 | 行号 | 内容 | 对比Java |
|------|------|------|----------|
| login_service.go | 134 | Personal个人组织标识写死为1 | ❓ |
| login_service.go | 136 | Avatar头像URL写死 | ❓ |
| login_service.go | 146 | 账户授权角色生成逻辑 | ❓ |

### 4. 常量定义（biz.go）- 9处 ⚠️ 可清理
| 文件 | 行号 | 内容 | 状态 |
|------|------|------|------|
| biz.go | 13 | Deleted常量未被使用 | 可删除 |
| biz.go | 19 | Archive常量未被使用 | 可删除 |
| biz.go | 23 | Open访问控制类型未使用 | 可删除 |
| biz.go | 25 | Custom未被使用 | 可删除 |
| biz.go | 34 | NoCollected未被使用 | 可删除 |
| biz.go | 39 | NoOwner未被使用 | 可删除 |
| biz.go | 44 | NoExecutor未被使用 | 可删除 |
| biz.go | 48 | NoCanRead/CanRead未被使用 | 可删除 |
| biz.go | 53 | UnDone未被使用 | 可删除 |

### 5. 缓存优化 - 5处
| 文件 | 行号 | 内容 |
|------|------|------|
| project-user/cache.go | 15 | TODO[缓存优化]: Token验证缓存 |
| project-project/cache.go | 16 | TODO[缓存优化]: 缓存问题修复 |
| project-api/midd.go | 3 | TODO[Token验证缓存优化] |
| project-user/router.go | 3 | TODO[缓存优化]: 启用缓存拦截器 |
| project-user/router.go | 68 | TODO[缓存优化]: 已启用缓存拦截器 |

### 6. 其他 - 8处
| 文件 | 行号 | 内容 | 优先级 |
|------|------|------|--------|
| member_account.go | 39 | SQL注入风险⚠️ | P0 |
| jwts.go | 3 | JWT刷新机制 | P2 |
| 前端index.js | 112 | 刷新Token机制 | P2 |
| run.go | 21 | HTTPS启用 | P3 |
| router.go | 87 | resolver合并 | P3 |
| project_service.go | 270 | 收藏功能优化 | P3 |
| main.go | 15 | 微服务拆分计划 | P3 |
| project.go | 131 | 字符串类型疑惑 | P3 |

---

## 二、Java版本现状

**Java版本中没有TODO注释！** 共52处TODO，分布如下：

| 类别 | Go版 | Java版 |
|------|------|--------|
| 部门模块 | 8处待实现 | ✅ 完整实现 |
| 任务模块 | 9处 | ✅ 一致（写死是正常的） |
| 登录模块 | 3处 | 待确认 |
| 常量 | 9处未使用 | 无 |
| 缓存 | 5处 | 可能简化了 |
| 其他 | 8处 | 无 |

---

## 三、分类总结

| 分类 | 数量 | 说明 |
|------|------|------|
| 待实现功能 | 8 | 部门模块CRUD - **需要实现** |
| 写死字段 | 6 | task相关 - **与Java一致，不需要改** |
| 未使用常量 | 9 | 可清理 |
| 优化项 | 5 | 缓存相关 |
| 风险 | 1 | SQL注入 |
| 其他 | 17 | 各种优化 |

---

## 四、优先级建议

### P0 - 必须处理（风险）
- [ ] SQL注入风险 (member_account.go:39)

### P1 - 重要（待实现）
- [ ] 部门模块CRUD（8处）- 参考Java实现

### P2 - 中等优化
- [ ] JWT刷新机制
- [ ] 前端刷新Token

### P3 - 低优先级
- [ ] 缓存优化（5处）
- [ ] 未使用常量清理（9处）
- [ ] 其他优化

---

## 五、结论

| 类型 | 数量 | 建议 |
|------|------|------|
| **需要实现** | 8 | 部门模块CRUD，参考Java |
| **不需要改** | 6 | 任务写死字段，与Java一致 |
| **需要清理** | 9 | 未使用常量 |
| **风险修复** | 1 | SQL注入 |
| **优化项** | 28 | 缓存、JWT等 |
