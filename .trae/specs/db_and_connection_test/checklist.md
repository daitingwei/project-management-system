# Checklist - 完善数据库与测试前后端连接验证

## 数据库测试数据验证（共20张表）

### 用户相关表
- [ ] ms_member表至少包含3条用户记录
- [ ] ms_member_account表有对应账号数据

### 组织架构表
- [ ] ms_organization表有组织数据（已有1条）
- [ ] ms_department表有部门数据（当前为0，需添加）

### 项目相关表
- [ ] ms_project表有项目数据（当前为0，需添加）
- [ ] ms_project_member表有项目成员数据
- [ ] ms_project_template表有模板数据（已有3条）
- [ ] ms_project_collection表有收藏数据
- [ ] ms_project_auth表有权限数据
- [ ] ms_project_auth_node表有权限节点数据
- [ ] ms_project_menu表有菜单数据
- [ ] ms_project_log表有日志数据

### 任务相关表
- [ ] ms_task_stages表有任务阶段数据（当前为0，需添加）
- [ ] ms_task表有任务数据（当前为0，需添加）
- [ ] ms_task_member表有任务成员数据
- [ ] ms_task_stages_template表有阶段模板数据（已有10条）
- [ ] ms_task_work_time表有工时数据

### 其他表
- [ ] ms_project_node表有节点数据（已有33条）
- [ ] ms_source_link表有外链数据（可选）
- [ ] ms_file表有文件数据（可选）

## 服务运行验证

- [ ] Redis容器运行在端口6379
- [ ] Redis可以正常连接（ping测试）
- [ ] Etcd容器运行在端口2379
- [ ] Etcd可以正常连接

## 前后端连接验证

- [ ] 后端API服务可访问
- [ ] 前端开发服务器可启动
- [ ] 登录页面可正常加载
- [ ] 登录请求可到达后端
- [ ] 登录成功可获取用户信息

## 功能验证

- [ ] 用户登录功能正常
- [ ] 组织架构数据加载正常
- [ ] 项目列表加载正常
- [ ] 项目成员显示正常
- [ ] 任务列表加载正常
- [ ] 任务阶段显示正常
