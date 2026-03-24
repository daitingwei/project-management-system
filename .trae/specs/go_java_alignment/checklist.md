# Checklist - Go版本对标Java版本验收清单

## 阶段一：部门管理模块（8个任务）
- [ ] Task 1: 部门删除功能 - 对标Java OrgController.departmentDelete()
- [ ] Task 2: 部门编辑功能 - 对标Java OrgController.departmentEdit()
- [ ] Task 3: 部门树形结构 - 获取部门层级关系
- [ ] Task 4: 添加部门成员 - 对标Java OrgService.inviteMember()
- [ ] Task 5: 移除部门成员 - 对标Java OrgService.removeMember()
- [ ] Task 6: 获取部门成员列表 - 对标Java OrgService.searchInviteMember()
- [ ] Task 7: 设置部门负责人
- [ ] Task 8: 移动部门（调整层级）

## 阶段二：项目管理模块（5个任务）
- [ ] Task 9: 项目归档功能 - 对标Java /project/archive
- [ ] Task 10: 取消项目归档 - 对标Java /project/recoveryArchive
- [ ] Task 11: 退出项目 - 对标Java /project/quit
- [ ] Task 12: 项目统计功能 - 对标Java /_projectStats
- [ ] Task 13: 项目报告功能 - 对标Java /_getProjectReport

## 阶段三：任务管理模块（10个任务）
- [ ] Task 14: 任务删除功能 - 对标Java /task/delete
- [ ] Task 15: 任务恢复功能 - 对标Java /task/recovery
- [ ] Task 16: 任务完成标记 - 对标Java /task/taskDone
- [ ] Task 17: 任务回收站 - 对标Java /task/recycle
- [ ] Task 18: 批量任务分配 - 对标Java /task/batchAssignTask
- [ ] Task 19: 任务标签管理 - 对标Java /task_tag/*
- [ ] Task 20: 任务点赞功能 - 对标Java /task/like
- [ ] Task 21: 任务收藏功能 - 对标Java /task/star
- [ ] Task 22: 任务评论优化 - 评论功能已实现
- [ ] Task 23: 任务日志完善 - 日志功能已实现

## 阶段四：成员管理模块（4个任务）
- [ ] Task 24: 搜索成员 - 对标Java /member/search
- [ ] Task 25: 邀请成员 - 对标Java /department_member/inviteMember
- [ ] Task 26: 移除成员 - 对标Java /department_member/removeMember
- [ ] Task 27: 批量导入成员 - 对标Java Excel导入

## 阶段五：权限管理模块（4个任务）
- [ ] Task 28: 权限编辑 - 对标Java /auth/edit
- [ ] Task 29: 权限添加 - 对标Java /auth/add
- [ ] Task 30: 权限删除 - 对标Java /auth/del
- [ ] Task 31: 设置默认角色 - 对标Java /auth/setDefault

## 阶段六：文件上传模块（7个任务）
- [x] Task 32: 通用文件下载 - 对标Java /common/download
- [x] Task 33: 文件删除 - 对标Java /source_link/delete
- [x] Task 51: 文件列表 - 对标前端 /project/file
- [x] Task 52: 文件详情 - 对标前端 /project/file/read
- [x] Task 53: 文件编辑 - 对标前端 /project/file/edit
- [x] Task 54: 文件回收 - 对标前端 /project/file/recycle
- [x] Task 55: 文件恢复 - 对标前端 /project/file/recovery
- [ ] Task 34: 头像上传 - 对标Java /index/uploadAvatar

## 阶段七：登录认证模块（5个任务）
- [ ] Task 35: 验证码功能 - 对标Java /project/login/getCaptcha
- [ ] Task 36: Token刷新机制 - 对标Java /index/index/refreshAccessToken
- [ ] Task 37: 修改密码 - 对标Java /project/index/editPassword
- [ ] Task 38: 绑定手机号 - 对标Java /project/login/_bindMobile
- [ ] Task 39: 绑定邮箱 - 对标Java /project/login/_bindMail

## 阶段八：组织管理模块（4个任务）
- [ ] Task 40: 组织列表 - 对标Java /organization/index
- [ ] Task 41: 组织编辑 - 对标Java /organization/edit
- [ ] Task 42: 切换组织 - 对标Java /index/changeCurrentOrganization
- [ ] Task 43: 组织树形 - 对标Java /_getOrgList

## 阶段九：菜单管理模块（4个任务）
- [ ] Task 44: 菜单添加 - 对标Java /menu/menuAdd
- [ ] Task 45: 菜单编辑 - 对标Java /menu/menuEdit
- [ ] Task 46: 菜单删除 - 对标Java /menu/menuDel
- [ ] Task 47: 菜单禁用/启用 - 对标Java /menu/menuResume

## 阶段十：其他功能（3个任务）
- [ ] Task 48: 邀请链接功能 - 对标Java /invite_link/save
- [ ] Task 49: WebSocket通知 - 对标Java Notify系统
- [ ] Task 50: 系统配置功能 - 对标Java SystemConfig

---

## 验收标准

### 通用标准
- [ ] 代码风格与现有Go代码一致
- [ ] 遵循GORM最佳实践
- [ ] 错误处理规范
- [ ] 日志记录完整

### 功能标准
- [ ] 所有API接口可正常调用
- [ ] 数据库操作正确
- [ ] 与Java版本功能一致

### 安全标准
- [ ] 无SQL注入风险
- [ ] 权限校验正确
- [ ] 敏感数据不暴露
