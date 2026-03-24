# Tasks - 完善数据库与测试前后端连接

## 阶段一：完善数据库测试数据

### 数据现状（20张表分析）
- ✅ 已有数据：ms_organization(1), ms_project_template(3), ms_task_stages_template(10), ms_project_node(33)
- ⚠️ 不足：ms_member(1)
- ❌ 需要添加：ms_department, ms_project, ms_project_member, ms_task, ms_task_stages, ms_member_account, ms_project_auth, ms_project_auth_node, ms_project_collection, ms_project_log, ms_project_menu, ms_source_link, ms_task_member, ms_task_work_time, ms_file

- [ ] Task 1: 添加测试用户数据
  - [ ] 1.1 插入2-4个测试用户到ms_member表
  - [ ] 1.2 插入用户账号关联到ms_member_account表
  - [ ] 1.3 验证用户数据（至少3条）

- [ ] Task 2: 添加组织架构数据
  - [ ] 2.1 插入部门数据到ms_department表（至少3个部门）
  - [ ] 2.2 验证部门数据

- [ ] Task 3: 添加项目数据
  - [ ] 3.1 插入项目数据到ms_project表（2-3个项目）
  - [ ] 3.2 插入项目成员到ms_project_member表
  - [ ] 3.3 插入项目收藏到ms_project_collection表
  - [ ] 3.4 插入项目权限到ms_project_auth表
  - [ ] 3.5 插入权限节点到ms_project_auth_node表
  - [ ] 3.6 插入项目菜单到ms_project_menu表
  - [ ] 3.7 插入项目日志到ms_project_log表

- [ ] Task 4: 添加任务数据
  - [ ] 4.1 插入任务阶段到ms_task_stages表（每个项目2-3个阶段）
  - [ ] 4.2 插入任务到ms_task表（每个项目3-5个任务）
  - [ ] 4.3 插入任务成员到ms_task_member表
  - [ ] 4.4 插入任务工时到ms_task_work_time表

- [ ] Task 5: 添加其他辅助数据
  - [ ] 5.1 插入外链数据到ms_source_link表（可选）
  - [ ] 5.2 插入文件数据到ms_file表（可选）

## 阶段二：启动必要服务

- [ ] Task 6: 启动Redis服务
  - [ ] 6.1 创建Redis数据目录
  - [ ] 6.2 创建Redis配置文件
  - [ ] 6.3 启动Redis容器
  - [ ] 6.4 验证Redis连接

- [ ] Task 7: 启动Etcd服务
  - [ ] 7.1 创建Etcd数据目录
  - [ ] 7.2 启动Etcd容器
  - [ ] 7.3 验证Etcd连接

## 阶段三：测试前后端连接

- [ ] Task 8: 验证后端服务
  - [ ] 8.1 启动project-api后端服务
  - [ ] 8.2 检查后端服务运行状态
  - [ ] 8.3 测试API端点可访问性

- [ ] Task 9: 验证前后端连接
  - [ ] 9.1 启动前端开发服务器
  - [ ] 9.2 测试登录功能
  - [ ] 9.3 测试数据加载（组织、部门、项目、任务）

# Task Dependencies
- Task 2 依赖于 Task 1（部门需要关联用户）
- Task 3 依赖于 Task 2（项目需要组织和部门）
- Task 4 依赖于 Task 3（任务需要项目）
- Task 5 可在Task 4后执行（辅助数据）
- Task 6, 7 可与 Task 1-5 并行执行
- Task 8, 9 依赖于 Task 6, 7 和 Task 1-5
