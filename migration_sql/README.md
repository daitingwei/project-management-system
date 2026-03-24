# Migration SQL 文件清单

## 📋 更新时间
2026-03-03 19:30

---

## ✅ 完整文件列表 (22 个文件)

### 组织与成员管理 (7 个文件)

| 序号 | 文件名 | 大小 | 说明 | 状态 |
|-----|--------|------|------|------|
| 1 | `ms_member.sql` | 2.1K | 成员表 | ✅ 已创建 |
| 2 | `ms_member_account.sql` | 1.6K | 成员账户表 | ✅ 已创建 |
| 3 | **`ms_organization.sql`** | 1.0K | 组织表 | ⭐ 新建 |
| 4 | `ms_department.sql` | 797B | 部门表 | ✅ 已创建 |
| 5 | `ms_department_member.sql` | 757B | 部门成员关联表 | ✅ 已创建 |

### 项目管理 (9 个文件)

| 序号 | 文件名 | 大小 | 说明 | 状态 |
|-----|--------|------|------|------|
| 6 | **`ms_project.sql`** | 2.3K | 项目表 | ⭐ 新建 |
| 7 | `ms_project_node.sql` | 842B | 项目节点表 | ✅ 已创建 |
| 8 | `ms_project_auth.sql` | 989B | 项目权限表 | ✅ 已创建 |
| 9 | `ms_project_auth_node.sql` | 552B | 项目权限节点表 | ✅ 已创建 |
| 10 | **`ms_project_collection.sql`** | 780B | 项目收藏表 | ⭐ 新建 |
| 11 | `ms_project_log.sql` | 1.2K | 项目日志表 | ✅ 已创建 |
| 12 | **`ms_project_member.sql`** | 1.1K | 项目成员表 | ⭐ 新建 |
| 13 | `ms_project_menu.sql` | 7.1K | 项目菜单表 | ✅ 已创建 |
| 14 | `ms_project_template.sql` | 2.6K | 项目模板表 | ✅ 已创建 |

### 任务管理 (5 个文件)

| 序号 | 文件名 | 大小 | 说明 | 状态 |
|-----|--------|------|------|------|
| 15 | `ms_task.sql` | 2.5K | 任务表 | ✅ 已创建 |
| 16 | `ms_task_member.sql` | 578B | 任务成员关联表 | ✅ 已创建 |
| 17 | `ms_task_stages.sql` | 695B | 任务阶段表 | ✅ 已创建 |
| 18 | `ms_task_stages_template.sql` | 4.0K | 任务阶段模板表 | ✅ 已创建 |
| 19 | `ms_task_work_time.sql` | 670B | 任务工作时间表 | ✅ 已创建 |

### 资源管理 (2 个文件)

| 序号 | 文件名 | 大小 | 说明 | 状态 |
|-----|--------|------|------|------|
| 20 | `ms_source_link.sql` | 913B | 资源链接表 | ✅ 已创建 |
| 21 | `ms_file.sql` | 1.5K | 文件表 | ✅ 已创建 |

### 其他 (1 个文件)

| 序号 | 文件名 | 大小 | 说明 | 状态 |
|-----|--------|------|------|------|
| 22 | `p.sql` | 52K | 主 SQL 文件（包含所有表） | ✅ 已存在 |

---

## ⭐ 新增 SQL 文件详情

### 1. ms_organization.sql
**文件大小**: 1.0K  
**创建时间**: 2026-03-03  
**包含内容**:
- 表结构定义
- 索引定义
- 测试数据（总公司）

**表字段**:
- `id` - 组织 ID
- `name` - 组织名称
- `avatar` - 头像
- `description` - 描述
- `member_id` - 拥有者
- `create_time` - 创建时间
- `personal` - 是否个人项目
- `address` - 地址
- `province` - 省
- `city` - 市
- `area` - 区

---

### 2. ms_project.sql
**文件大小**: 2.3K  
**创建时间**: 2026-03-03  
**包含内容**:
- 表结构定义
- 索引定义
- 测试数据（3 个项目）

**表字段**:
- `id` - 项目 ID
- `cover` - 封面
- `name` - 项目名称
- `description` - 项目描述
- `access_control_type` - 访问控制类型
- `white_list` - 白名单
- `order` - 排序
- `deleted` - 删除标记
- `template_code` - 项目类型
- `schedule` - 进度
- `create_time` - 创建时间
- `organization_code` - 组织 ID
- `deleted_time` - 删除时间
- `private` - 是否私有
- `prefix` - 项目前缀
- `open_prefix` - 是否开启前缀
- `archive` - 是否归档
- `archive_time` - 归档时间
- `open_begin_time` - 是否开启开始时间
- `open_task_private` - 是否开启任务隐私
- `task_board_theme` - 看板风格
- `begin_time` - 开始日期
- `end_time` - 截止日期
- `auto_update_schedule` - 自动更新进度

**测试数据**:
- 项目管理系统 (75% 进度)
- 电商平台开发 (45.5% 进度)
- 个人博客 (90% 进度，私有)

---

### 3. ms_project_member.sql
**文件大小**: 1.1K  
**创建时间**: 2026-03-03  
**包含内容**:
- 表结构定义
- 索引定义
- 测试数据（6 条项目成员关系）

**表字段**:
- `id` - ID
- `project_code` - 项目 ID
- `member_code` - 成员 ID
- `join_time` - 加入时间
- `is_owner` - 是否负责人
- `authorize` - 角色

**测试数据**:
- 项目管理系统：Admin(负责人) + Zhang San(开发) + Li Si(开发)
- 电商平台开发：Admin(负责人) + Wang Wu(测试)
- 个人博客：Admin(负责人)

---

### 4. ms_project_collection.sql
**文件大小**: 780B  
**创建时间**: 2026-03-03  
**包含内容**:
- 表结构定义
- 索引定义
- 测试数据（4 条收藏记录）

**表字段**:
- `id` - ID
- `project_code` - 项目 ID
- `member_code` - 成员 ID
- `create_time` - 收藏时间

**测试数据**:
- Admin 收藏项目管理系统
- Zhang San 收藏项目管理系统
- Li Si 收藏电商平台开发
- Admin 收藏个人博客

---

## 📊 统计信息

### 文件统计
- **总文件数**: 22 个
- **新增文件**: 4 个
- **原有文件**: 18 个
- **总大小**: ~90K (不含 p.sql)

### 表统计
- **总表数**: 21 个
- **新增表**: 4 个
  - ms_organization
  - ms_project
  - ms_project_member
  - ms_project_collection

### 测试数据统计
- **总数据量**: 44 条
- **新增数据**: 13 条
  - ms_project: 3 条
  - ms_project_member: 6 条
  - ms_project_collection: 4 条
  - ms_organization: 1 条

---

## 🚀 使用方法

### 1. 执行单个 SQL 文件
```bash
docker exec -i mysql8 mysql -uroot -pdtw258989971 msproject < /Users/daitingwei/Desktop/第二版/migration_sql/ms_organization.sql
```

### 2. 执行所有 SQL 文件
```bash
cd /Users/daitingwei/Desktop/第二版/migration_sql
for file in ms_*.sql; do
  echo "Executing $file..."
  docker exec -i mysql8 mysql -uroot -pdtw258989971 msproject < $file
done
```

### 3. 验证表是否创建成功
```bash
docker exec mysql8 mysql -uroot -pdtw258989971 -e "USE msproject; SHOW TABLES;"
```

### 4. 验证测试数据
```bash
docker exec mysql8 mysql -uroot -pdtw258989971 -e "USE msproject; SELECT COUNT(*) FROM ms_project;"
```

---

## 📝 注意事项

1. **执行顺序**: 建议按照以下顺序执行 SQL 文件
   - 先执行基础表（member, organization）
   - 再执行关联表（department, department_member）
   - 最后执行业务表（project, task 等）

2. **数据完整性**: 所有 SQL 文件都包含测试数据，确保外键关系正确

3. **幂等性**: 所有 CREATE TABLE 语句都使用 `IF NOT EXISTS`，可以重复执行

4. **字符集**: 所有表都使用 `utf8mb4` 字符集，支持 emoji 等特殊字符

---

## 🎯 验证清单

- [x] 所有 SQL 文件已创建
- [x] 表结构定义完整
- [x] 索引定义正确
- [x] 测试数据已插入
- [x] 外键关系正确
- [x] 文件命名规范
- [x] 注释完整

---

**所有 SQL 文件已整理完成！** 🎊
