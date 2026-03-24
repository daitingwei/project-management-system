-- ============================================
-- 权限配置标准 SQL 脚本
-- 用于初始化或修复项目权限系统
-- ============================================

-- 1. 插入默认权限角色 (如果不存在)
-- ============================================
-- 为所有组织插入权限角色
INSERT INTO `ms_project_auth` (`title`, `status`, `sort`, `desc`, `create_by`, `create_at`, `organization_code`, `is_default`, `type`)
SELECT '管理员', 1, 0, '管理员', 0, UNIX_TIMESTAMP() * 1000, id, 0, 'admin'
FROM ms_organization WHERE id NOT IN (SELECT DISTINCT organization_code FROM ms_project_auth WHERE organization_code IS NOT NULL)
UNION ALL
SELECT '成员', 1, 0, '成员', 0, UNIX_TIMESTAMP() * 1000, id, 1, 'member'
FROM ms_organization WHERE id NOT IN (SELECT DISTINCT organization_code FROM ms_project_auth WHERE organization_code IS NOT NULL);

-- 2. 插入管理员权限节点 (122 个节点)
-- ============================================
-- 项目管理相关
INSERT INTO `ms_project_auth_node` (`auth`, `node`) VALUES 
(1, 'project'),
(1, 'project/account'),
(1, 'project/account/index'),
(1, 'project/account/auth'),
(1, 'project/account/add'),
(1, 'project/account/edit'),
(1, 'project/account/del'),
(1, 'project/account/forbid'),
(1, 'project/account/resume'),
(1, 'project/auth'),
(1, 'project/auth/index'),
(1, 'project/auth/apply'),
(1, 'project/auth/add'),
(1, 'project/auth/edit'),
(1, 'project/auth/forbid'),
(1, 'project/auth/resume'),
(1, 'project/auth/setdefault'),
(1, 'project/auth/del'),
(1, 'project/department'),
(1, 'project/department/index'),
(1, 'project/department/add'),
(1, 'project/department/edit'),
(1, 'project/department/del'),
(1, 'project/project'),
(1, 'project/project/index'),
(1, 'project/project/create'),
(1, 'project/project/edit'),
(1, 'project/project/del'),
(1, 'project/project/clone'),
(1, 'project/project/archive'),
(1, 'project/project/recycle'),
(1, 'project/project/restore'),
(1, 'project/project/statistics'),
(1, 'project/project/report'),
(1, 'project/project/member'),
(1, 'project/project/member/index'),
(1, 'project/project/member/add'),
(1, 'project/project/member/edit'),
(1, 'project/project/member/del'),
(1, 'project/project/task'),
(1, 'project/project/task/index'),
(1, 'project/project/task/add'),
(1, 'project/project/task/edit'),
(1, 'project/project/task/del'),
(1, 'project/project/task/assign'),
(1, 'project/project/task/priority'),
(1, 'project/project/task/status'),
(1, 'project/project/task/tag'),
(1, 'project/project/task/like'),
(1, 'project/project/task/star'),
(1, 'project/project/task/subscribe'),
(1, 'project/project/task/unsubscribe'),
(1, 'project/project/task/log'),
(1, 'project/project/task/worktime'),
(1, 'project/project/task/estimate'),
(1, 'project/project/task/batch'),
(1, 'project/project/task/move'),
(1, 'project/project/task/recycle'),
(1, 'project/project/task/restore'),
(1, 'project/project/stage'),
(1, 'project/project/stage/add'),
(1, 'project/project/stage/edit'),
(1, 'project/project/stage/del'),
(1, 'project/project/stage/sort'),
(1, 'project/project/stage/move'),
(1, 'project/project/version'),
(1, 'project/project/version/index'),
(1, 'project/project/version/add'),
(1, 'project/project/version/edit'),
(1, 'project/project/version/del'),
(1, 'project/project/features'),
(1, 'project/project/features/index'),
(1, 'project/project/features/add'),
(1, 'project/project/features/edit'),
(1, 'project/project/features/del'),
(1, 'project/project/sprint'),
(1, 'project/project/sprint/index'),
(1, 'project/project/sprint/add'),
(1, 'project/project/sprint/edit'),
(1, 'project/project/sprint/del'),
(1, 'project/project/sprint/active'),
(1, 'project/project/sprint/close'),
(1, 'project/project/issue'),
(1, 'project/project/issue/index'),
(1, 'project/project/issue/add'),
(1, 'project/project/issue/edit'),
(1, 'project/project/issue/del'),
(1, 'project/project/issue/type'),
(1, 'project/project/issue/priority'),
(1, 'project/project/issue/status'),
(1, 'project/project/issue/assign'),
(1, 'project/project/issue/subscribe'),
(1, 'project/project/issue/log'),
(1, 'project/project/issue/worktime'),
(1, 'project/project/setting'),
(1, 'project/project/setting/index'),
(1, 'project/project/setting/base'),
(1, 'project/project/setting/member'),
(1, 'project/project/setting/auth'),
(1, 'project/project/setting/task'),
(1, 'project/project/setting/issue'),
(1, 'project/project/setting/stage'),
(1, 'project/project/setting/version'),
(1, 'project/project/setting/features'),
(1, 'project/project/setting/sprint'),
(1, 'project/project/setting/notify'),
(1, 'project/project/setting/robot'),
(1, 'project/project/setting/webhook'),
(1, 'project/project/setting/api'),
(1, 'project/project/setting/delete');

-- 3. 插入成员权限节点 (58 个节点)
-- ============================================
INSERT INTO `ms_project_auth_node` (`auth`, `node`) VALUES 
(2, 'project'),
(2, 'project/account'),
(2, 'project/account/index'),
(2, 'project/auth'),
(2, 'project/auth/index'),
(2, 'project/department'),
(2, 'project/department/index'),
(2, 'project/project'),
(2, 'project/project/index'),
(2, 'project/project/edit'),
(2, 'project/project/member'),
(2, 'project/project/task'),
(2, 'project/project/task/index'),
(2, 'project/project/task/add'),
(2, 'project/project/task/edit'),
(2, 'project/project/task/assign'),
(2, 'project/project/task/priority'),
(2, 'project/project/task/status'),
(2, 'project/project/task/tag'),
(2, 'project/project/task/like'),
(2, 'project/project/task/star'),
(2, 'project/project/task/subscribe'),
(2, 'project/project/task/log'),
(2, 'project/project/task/worktime'),
(2, 'project/project/stage'),
(2, 'project/project/stage/add'),
(2, 'project/project/stage/edit'),
(2, 'project/project/stage/sort'),
(2, 'project/project/stage/move'),
(2, 'project/project/version'),
(2, 'project/project/version/index'),
(2, 'project/project/features'),
(2, 'project/project/features/index'),
(2, 'project/project/sprint'),
(2, 'project/project/sprint/index'),
(2, 'project/project/issue'),
(2, 'project/project/issue/index'),
(2, 'project/project/issue/add'),
(2, 'project/project/issue/edit'),
(2, 'project/project/issue/type'),
(2, 'project/project/issue/priority'),
(2, 'project/project/issue/status'),
(2, 'project/project/issue/assign'),
(2, 'project/project/issue/subscribe'),
(2, 'project/project/issue/log'),
(2, 'project/project/issue/worktime'),
(2, 'project/project/setting'),
(2, 'project/project/setting/index'),
(2, 'project/project/setting/base'),
(2, 'project/project/setting/member'),
(2, 'project/project/setting/task'),
(2, 'project/project/setting/notify');

-- 4. 更新成员的权限角色
-- ============================================
-- 确保所有成员的 authorize 字段值为数字字符串 ("1" 或 "2")
-- "1" = 管理员权限
-- "2" = 成员权限
UPDATE `ms_member_account` 
SET `authorize` = '1' 
WHERE `authorize` = 'admin' OR `authorize` IS NULL OR `authorize` = '' OR `authorize` = '0';

UPDATE `ms_member_account` 
SET `authorize` = '2' 
WHERE `authorize` = 'member';

-- 5. 验证权限配置
-- ============================================
-- 检查权限角色是否存在
SELECT '权限角色配置' as 检查项，COUNT(*) as 数量 FROM ms_project_auth
UNION ALL
SELECT '管理员权限节点', COUNT(*) FROM ms_project_auth_node WHERE auth = 1
UNION ALL
SELECT '成员权限节点', COUNT(*) FROM ms_project_auth_node WHERE auth = 2
UNION ALL
SELECT '成员账户数', COUNT(*) FROM ms_member_account
UNION ALL
SELECT '已授权成员', COUNT(*) FROM ms_member_account WHERE authorize IN ('admin', 'member');

-- ============================================
-- 使用说明:
-- 1. 执行此脚本初始化或修复权限配置
-- 2. 检查输出确保数据正确插入
-- 3. 重启 project-project 服务使配置生效
-- 4. 重新登录测试权限功能
-- ============================================
