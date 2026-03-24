# Tasks - Go版本对标Java版本完整任务清单

## 阶段一：部门管理模块（8个任务）
- [ ] Task 1: 部门删除功能 - 参考Java OrgController.departmentDelete()
- [ ] Task 2: 部门编辑功能 - 参考Java OrgController.departmentEdit()
- [ ] Task 3: 部门树形结构 - 获取部门层级关系
- [ ] Task 4: 添加部门成员 - 参考Java OrgService.inviteMember()
- [ ] Task 5: 移除部门成员 - 参考Java OrgService.removeMember()
- [ ] Task 6: 获取部门成员列表 - 参考Java OrgService.searchInviteMember()
- [ ] Task 7: 设置部门负责人
- [ ] Task 8: 移动部门（调整层级）

## 阶段二：项目管理模块（5个任务）
- [ ] Task 9: 项目归档功能
- [ ] Task 10: 取消项目归档
- [ ] Task 11: 退出项目
- [ ] Task 12: 项目统计功能
- [ ] Task 13: 项目报告功能

## 阶段三：任务管理模块（10个任务）
- [ ] Task 14: 任务删除功能
- [ ] Task 15: 任务恢复功能
- [ ] Task 16: 任务完成标记
- [ ] Task 17: 任务回收站
- [ ] Task 18: 批量任务分配
- [ ] Task 19: 任务标签管理
- [ ] Task 20: 任务点赞功能
- [ ] Task 21: 任务收藏功能
- [ ] Task 22: 任务评论优化
- [ ] Task 23: 任务日志完善

## 阶段四：成员管理模块（4个任务）
- [ ] Task 24: 搜索成员
- [ ] Task 25: 邀请成员
- [ ] Task 26: 移除成员
- [ ] Task 27: 批量导入成员

## 阶段五：权限管理模块（4个任务）
- [ ] Task 28: 权限编辑
- [ ] Task 29: 权限添加
- [ ] Task 30: 权限删除
- [ ] Task 31: 设置默认角色

## 阶段六：文件上传模块（7个任务）
- [x] Task 32: 通用文件下载
- [x] Task 33: 文件删除
- [x] Task 51: 文件列表 - 对标前端 /project/file
- [x] Task 52: 文件详情 - 对标前端 /project/file/read
- [x] Task 53: 文件编辑 - 对标前端 /project/file/edit
- [x] Task 54: 文件回收 - 对标前端 /project/file/recycle
- [x] Task 55: 文件恢复 - 对标前端 /project/file/recovery

## 阶段七：登录认证模块（5个任务）
- [ ] Task 35: 验证码功能
- [ ] Task 36: Token刷新机制
- [ ] Task 37: 修改密码
- [ ] Task 38: 绑定手机号
- [ ] Task 39: 绑定邮箱

## 阶段八：组织管理模块（4个任务）
- [ ] Task 40: 组织列表
- [ ] Task 41: 组织编辑
- [ ] Task 42: 切换组织
- [ ] Task 43: 组织树形

## 阶段九：菜单管理模块（4个任务）
- [ ] Task 44: 菜单添加
- [ ] Task 45: 菜单编辑
- [ ] Task 46: 菜单删除
- [ ] Task 47: 菜单禁用/启用

## 阶段十：其他功能（3个任务）
- [ ] Task 48: 邀请链接功能
- [ ] Task 49: WebSocket通知（可选）
- [ ] Task 50: 系统配置功能

---

## Task Dependencies
- Task 1-8 部门管理模块：内部有依赖，部门树需在列表基础上实现
- Task 14-23 任务管理：任务删除需先实现回收站
- Task 35-39 登录认证：验证码是其他功能的基础
- Task 40-43 组织管理：组织列表是其他功能的基础
