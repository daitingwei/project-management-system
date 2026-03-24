# 数据库表同步计划

## 问题分析

经过对比分析 `teamwork0203.sql` 和 `teamwork0309_init.sql` 两个文件，发现以下问题:

### 1. 数据库表结构差异

**teamwork0203.sql** (旧版本 - 2021/02/03):
- 包含基础表结构
- 字段相对较少
- 部分表字段类型不一致

**teamwork0309_init.sql** (新版本 - 2021/03/09):
- 包含更多表结构
- 字段更完善
- 添加了新表：`team_project_log`、`team_source_link`、`team_comment` 等
- 部分字段类型和长度有调整

### 2. Go 项目数据模型与数据库表的差异

**当前 Go 项目使用的表名:**
- `ms_member` (对应数据库：`team_member`)
- `ms_project` (对应数据库：`team_project`)
- `ms_department` (对应数据库：`team_department`)
- `ms_member_account` (对应数据库：`team_member_account`)
- `ms_task` (对应数据库：缺少此表)
- `ms_project_member` (对应数据库：缺少此表)
- `ms_project_collection` (对应数据库：缺少此表)
- `ms_project_template` (对应数据库：缺少此表)
- `ms_task_stages_template` (对应数据库：缺少此表)

**主要问题:**
1. **表名不匹配**: Go 项目使用 `ms_` 前缀，数据库使用 `team_` 前缀
2. **字段类型不匹配**: 
   - Go 项目使用 `int64` 时间戳，数据库使用 `varchar` 存储时间字符串
   - 部分字段类型和长度不一致
3. **缺少表**: Go 项目需要的部分表在数据库中不存在

## 解决方案

### 方案一：删除所有旧表，创建新表 (推荐)

按照用户要求，删除数据库中所有 `team_` 开头的表，创建 `ms_` 开头的新表，字段类型与 Go 项目保持一致。

### 实施步骤

#### 步骤 1: 分析 Go 项目所有数据模型

需要创建以下表:

1. **用户相关表**
   - `ms_member` - 用户表
   - `ms_member_account` - 组织账号表

2. **组织相关表**
   - `ms_organization` - 组织表

3. **部门相关表**
   - `ms_department` - 部门表
   - `ms_department_member` - 部门成员表

4. **项目相关表**
   - `ms_project` - 项目表
   - `ms_project_member` - 项目成员表
   - `ms_project_collection` - 项目收藏表
   - `ms_project_template` - 项目模板表
   - `ms_project_auth` - 项目权限表
   - `ms_project_auth_node` - 项目权限节点表
   - `ms_project_log` - 项目日志表

5. **任务相关表**
   - `ms_task` - 任务表
   - `ms_task_member` - 任务成员表
   - `ms_task_stages` - 任务阶段表
   - `ms_task_stages_template` - 任务阶段模板表
   - `ms_task_workflow` - 任务工作流表
   - `ms_task_tag` - 任务标签表
   - `ms_task_work_time` - 任务工时表
   - `ms_task_log` - 任务日志表

6. **文件相关表**
   - `ms_file` - 文件表
   - `ms_source_link` - 源链接表
   - `ms_comment` - 评论表

7. **系统相关表**
   - `ms_menu` - 菜单表
   - `ms_node` - 节点表
   - `ms_notify` - 通知表
   - `ms_invite_link` - 邀请链接表
   - `ms_lock` - 防灌水表
   - `ms_mailqueue` - 邮件队列表
   - `ms_collection` - 收藏表

#### 步骤 2: 创建新的 SQL 文件

创建 `ms_project_init.sql` 文件，包含:
- 所有表的 DROP TABLE 语句
- 所有表的 CREATE TABLE 语句 (使用 `ms_` 前缀)
- 字段类型与 Go 项目保持一致 (时间使用 int64/bigint)
- 必要的索引和约束
- 初始化数据

#### 步骤 3: 字段类型映射

| Go 类型 | MySQL 类型 | 说明 |
|---------|-----------|------|
| int64 (时间戳) | bigint | 时间字段统一使用 bigint 存储毫秒时间戳 |
| int | int | 普通整型 |
| int (布尔) | tinyint | 布尔值使用 tinyint(1) |
| string | varchar | 根据长度设置 varchar |
| string (长文本) | text | 长文本使用 text |

#### 步骤 4: 表名映射

| 原表名 (team_) | 新表名 (ms_) |
|--------------|------------|
| team_member | ms_member |
| team_member_account | ms_member_account |
| team_organization | ms_organization |
| team_department | ms_department |
| team_department_member | ms_department_member |
| team_project | ms_project |
| team_project_member | ms_project_member |
| team_project_collection | ms_project_collection |
| team_project_template | ms_project_template |
| team_project_auth | ms_project_auth |
| team_project_auth_node | ms_project_auth_node |
| team_project_log | ms_project_log |
| team_task | ms_task |
| team_task_member | ms_task_member |
| team_task_stages | ms_task_stages |
| team_task_stages_template | ms_task_stages_template |
| team_file | ms_file |
| team_source_link | ms_source_link |
| team_comment | ms_comment |
| team_menu | ms_menu |
| team_node | ms_node |
| team_notify | ms_notify |
| team_invite_link | ms_invite_link |
| team_lock | ms_lock |
| team_mailqueue | ms_mailqueue |
| team_collection | ms_collection |

#### 步骤 5: 验证

1. 检查所有表是否创建成功
2. 检查字段类型是否与 Go 项目匹配
3. 检查索引是否创建正确
4. 测试 Go 项目是否能正常访问数据库

## 预期结果

1. 数据库中所有 `team_` 开头的表被删除
2. 创建所有 `ms_` 开头的新表
3. 字段类型与 Go 项目完全匹配
4. Go 项目可以正常访问和操作数据库

## 注意事项

1. 删除表前需要确认是否有重要数据需要备份
2. 时间字段统一使用 bigint 存储毫秒时间戳
3. 保持必要的索引以提高查询性能
4. 保留必要的外键约束 (如果有)
5. 确保字符集使用 utf8mb4
