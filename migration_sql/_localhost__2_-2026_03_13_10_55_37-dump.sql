-- MySQL dump 10.13  Distrib 8.0.45, for macos26.2 (arm64)
--
-- Host: 127.0.0.1    Database: msproject
-- ------------------------------------------------------
-- Server version	8.0.35

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ 'fa02004c-1dd1-11f1-ac75-56130ae051ba:1-5420';

--
-- Table structure for table `ms_department`
--

DROP TABLE IF EXISTS `ms_department`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_department` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `organization_code` bigint DEFAULT NULL COMMENT '组织编号',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '名称',
  `sort` int DEFAULT '0' COMMENT '排序',
  `pcode` bigint DEFAULT NULL COMMENT '上级编号',
  `icon` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '图标',
  `create_time` bigint DEFAULT NULL COMMENT '创建时间',
  `path` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '上级路径',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='部门表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_department`
--

LOCK TABLES `ms_department` WRITE;
/*!40000 ALTER TABLE `ms_department` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_department` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_department_member`
--

DROP TABLE IF EXISTS `ms_department_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_department_member` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `department_code` bigint DEFAULT NULL COMMENT '部门id',
  `organization_code` bigint DEFAULT NULL COMMENT '组织id',
  `account_code` bigint DEFAULT NULL COMMENT '成员id',
  `join_time` bigint DEFAULT NULL COMMENT '加入时间',
  `is_principal` tinyint(1) DEFAULT NULL COMMENT '是否负责人',
  `is_owner` tinyint(1) DEFAULT '0' COMMENT '拥有者',
  `authorize` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '角色',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='部门-成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_department_member`
--

LOCK TABLES `ms_department_member` WRITE;
/*!40000 ALTER TABLE `ms_department_member` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_department_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_file`
--

DROP TABLE IF EXISTS `ms_file`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_file` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `path_name` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '相对路径',
  `title` char(90) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '名称',
  `extension` char(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '扩展名',
  `size` int unsigned DEFAULT '0' COMMENT '文件大小',
  `object_type` char(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '对象类型',
  `organization_code` bigint DEFAULT NULL COMMENT '组织编码',
  `task_code` bigint DEFAULT NULL COMMENT '任务编码',
  `project_code` bigint DEFAULT NULL COMMENT '项目编码',
  `create_by` bigint DEFAULT NULL COMMENT '上传人',
  `create_time` bigint DEFAULT NULL COMMENT '创建时间',
  `downloads` mediumint unsigned DEFAULT '0' COMMENT '下载次数',
  `extra` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '额外信息',
  `deleted` tinyint(1) DEFAULT '0' COMMENT '删除标记',
  `file_url` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci COMMENT '完整地址',
  `file_type` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '文件类型',
  `deleted_time` bigint DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_file`
--

LOCK TABLES `ms_file` WRITE;
/*!40000 ALTER TABLE `ms_file` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_file` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_member`
--

DROP TABLE IF EXISTS `ms_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_member` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ç³»ç»Ÿå‰å°ç”¨æˆ·è¡¨',
  `account` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT 'ç”¨æˆ·ç™»é™†è´¦å·',
  `password` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT 'ç™»é™†å¯†ç ',
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT 'ç”¨æˆ·æ˜µç§°',
  `mobile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'æ‰‹æœº',
  `realname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'çœŸå®žå§“å',
  `create_time` bigint DEFAULT NULL COMMENT 'åˆ›å»ºæ—¶é—´',
  `status` tinyint(1) DEFAULT '0' COMMENT 'çŠ¶æ€',
  `last_login_time` bigint DEFAULT NULL COMMENT 'ä¸Šæ¬¡ç™»å½•æ—¶é—´',
  `sex` tinyint DEFAULT '0' COMMENT 'æ€§åˆ«',
  `avatar` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT 'å¤´åƒ',
  `idcard` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT 'èº«ä»½è¯',
  `province` int DEFAULT '0' COMMENT 'çœ',
  `city` int DEFAULT '0' COMMENT 'å¸‚',
  `area` int DEFAULT '0' COMMENT 'åŒº',
  `address` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT 'æ‰€åœ¨åœ°å€',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT 'å¤‡æ³¨',
  `email` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT 'é‚®ç®±',
  `dingtalk_openid` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT 'é’‰é’‰openid',
  `dingtalk_unionid` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT 'é’‰é’‰unionid',
  `dingtalk_userid` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT 'é’‰é’‰ç”¨æˆ·id',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1002 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='ç”¨æˆ·è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_member`
--

LOCK TABLES `ms_member` WRITE;
/*!40000 ALTER TABLE `ms_member` DISABLE KEYS */;
INSERT INTO `ms_member` VALUES (1000,'代庭玮','62ded160b46d1940d384b200e1acba62','代庭玮','17606254275','',1773301495938,1,1773301495938,0,'','',0,0,0,'','','dai2004105@163.com','','',''),(1001,'qry','20e0f295238e7083a9e9bc598f071252','qry','17606254276','',1773316226133,1,1773316226133,0,'','',0,0,0,'','','805@qry.com','','','');
/*!40000 ALTER TABLE `ms_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_member_account`
--

DROP TABLE IF EXISTS `ms_member_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_member_account` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `member_code` bigint DEFAULT NULL COMMENT '所属账号id',
  `organization_code` bigint DEFAULT NULL COMMENT '所属组织',
  `department_code` bigint DEFAULT NULL COMMENT '部门编号',
  `authorize` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '角色',
  `is_owner` tinyint(1) DEFAULT '0' COMMENT '是否主账号',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '姓名',
  `mobile` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '手机号码',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮件',
  `create_time` bigint DEFAULT NULL COMMENT '创建时间',
  `last_login_time` bigint DEFAULT NULL COMMENT '上次登录时间',
  `status` tinyint(1) DEFAULT '0' COMMENT '状态0禁用 1使用中',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '头像',
  `position` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '职位',
  `department` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部门',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='组织账号表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_member_account`
--

LOCK TABLES `ms_member_account` WRITE;
/*!40000 ALTER TABLE `ms_member_account` DISABLE KEYS */;
INSERT INTO `ms_member_account` VALUES (35,1000,1,NULL,'1',1,'ä»£åº­çŽ®',NULL,NULL,1773311753000,NULL,1,NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `ms_member_account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_organization`
--

DROP TABLE IF EXISTS `ms_organization`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_organization` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'åç§°',
  `avatar` varchar(511) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'å¤´åƒ',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'æè¿°',
  `member_id` bigint DEFAULT NULL COMMENT 'æ‹¥æœ‰è€…',
  `create_time` bigint DEFAULT NULL COMMENT 'åˆ›å»ºæ—¶é—´',
  `personal` tinyint(1) DEFAULT '0' COMMENT 'æ˜¯å¦ä¸ªäººé¡¹ç›®',
  `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'åœ°å€',
  `province` int DEFAULT '0' COMMENT 'çœ',
  `city` int DEFAULT '0' COMMENT 'å¸‚',
  `area` int DEFAULT '0' COMMENT 'åŒº',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='ç»„ç»‡è¡¨';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_organization`
--

LOCK TABLES `ms_organization` WRITE;
/*!40000 ALTER TABLE `ms_organization` DISABLE KEYS */;
INSERT INTO `ms_organization` VALUES (8,'代庭玮个人组织','https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5','',1000,1773301495941,1,'',0,0,0),(9,'qry个人组织','https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5','',1001,1773316226136,1,'',0,0,0);
/*!40000 ALTER TABLE `ms_organization` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project`
--

DROP TABLE IF EXISTS `ms_project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `cover` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '封面',
  `name` varchar(90) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '名称',
  `description` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci COMMENT '描述',
  `access_control_type` tinyint DEFAULT '0' COMMENT '访问控制l类型',
  `white_list` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '可以访问项目的权限组（白名单）',
  `sort` int unsigned DEFAULT '0' COMMENT '排序',
  `deleted` tinyint(1) DEFAULT '0' COMMENT '删除标记',
  `template_code` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT '项目类型',
  `schedule` double(5,2) DEFAULT '0.00' COMMENT '进度',
  `create_time` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '创建时间',
  `organization_code` bigint DEFAULT NULL COMMENT '组织id',
  `deleted_time` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '删除时间',
  `private` tinyint(1) DEFAULT '1' COMMENT '是否私有',
  `prefix` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '项目前缀',
  `open_prefix` tinyint(1) DEFAULT '0' COMMENT '是否开启项目前缀',
  `archive` tinyint(1) DEFAULT '0' COMMENT '是否归档',
  `archive_time` bigint DEFAULT NULL COMMENT '归档时间',
  `open_begin_time` tinyint(1) DEFAULT '0' COMMENT '是否开启任务开始时间',
  `open_task_private` tinyint(1) DEFAULT '0' COMMENT '是否开启新任务默认开启隐私模式',
  `task_board_theme` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT 'default' COMMENT '看板风格',
  `begin_time` bigint DEFAULT NULL COMMENT '项目开始日期',
  `end_time` bigint DEFAULT NULL COMMENT '项目截止日期',
  `auto_update_schedule` tinyint(1) DEFAULT '0' COMMENT '自动更新项目进度',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `project` (`sort`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13043 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='项目表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project`
--

LOCK TABLES `ms_project` WRITE;
/*!40000 ALTER TABLE `ms_project` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_project` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_auth`
--

DROP TABLE IF EXISTS `ms_project_auth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_auth` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '权限名称',
  `status` tinyint unsigned DEFAULT '1' COMMENT '状态(0:禁用,1:启用)',
  `sort` smallint unsigned DEFAULT '0' COMMENT '排序权重',
  `desc` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '备注说明',
  `create_by` bigint unsigned DEFAULT '0' COMMENT '创建人',
  `create_at` bigint DEFAULT NULL COMMENT '创建时间',
  `organization_code` bigint DEFAULT NULL COMMENT '所属组织',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认',
  `type` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '权限类型',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='项目权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_auth`
--

LOCK TABLES `ms_project_auth` WRITE;
/*!40000 ALTER TABLE `ms_project_auth` DISABLE KEYS */;
INSERT INTO `ms_project_auth` VALUES (1,'管理员',1,0,'管理员',0,1773296332341,NULL,0,'admin'),(2,'成员',1,0,'成员',0,1773296332343,NULL,1,'member');
/*!40000 ALTER TABLE `ms_project_auth` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_auth_node`
--

DROP TABLE IF EXISTS `ms_project_auth_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_auth_node` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `auth` bigint unsigned DEFAULT NULL COMMENT '角色ID',
  `node` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '节点路径',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `index_system_auth_auth` (`auth`) USING BTREE,
  KEY `index_system_auth_node` (`node`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5280 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='项目角色与节点绑定';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_auth_node`
--

LOCK TABLES `ms_project_auth_node` WRITE;
/*!40000 ALTER TABLE `ms_project_auth_node` DISABLE KEYS */;
INSERT INTO `ms_project_auth_node` VALUES (3097,1,'project'),(3098,1,'project/account'),(3099,1,'project/account/index'),(3100,1,'project/account/auth'),(3101,1,'project/account/add'),(3102,1,'project/account/edit'),(3103,1,'project/account/del'),(3104,1,'project/account/forbid'),(3105,1,'project/account/resume'),(3106,1,'project/auth'),(3107,1,'project/auth/index'),(3108,1,'project/auth/apply'),(3109,1,'project/auth/add'),(3110,1,'project/auth/edit'),(3111,1,'project/auth/forbid'),(3112,1,'project/auth/resume'),(3113,1,'project/auth/setdefault'),(3114,1,'project/auth/del'),(3115,1,'project/department'),(3116,1,'project/department/index'),(3117,1,'project/department/read'),(3118,1,'project/department/save'),(3119,1,'project/department/edit'),(3120,1,'project/department/delete'),(3121,1,'project/department_member'),(3122,1,'project/department_member/index'),(3123,1,'project/department_member/searchinvitemember'),(3124,1,'project/department_member/invitemember'),(3125,1,'project/department_member/removemember'),(3126,1,'project/index'),(3127,1,'project/index/index'),(3128,1,'project/index/changecurrentorganization'),(3129,1,'project/index/systemconfig'),(3130,1,'project/index/info'),(3131,1,'project/index/editpersonal'),(3132,1,'project/index/editpassword'),(3133,1,'project/index/uploadimg'),(3134,1,'project/index/uploadavatar'),(3135,1,'project/menu'),(3136,1,'project/menu/menu'),(3137,1,'project/menu/menuadd'),(3138,1,'project/menu/menuedit'),(3139,1,'project/menu/menuforbid'),(3140,1,'project/menu/menuresume'),(3141,1,'project/menu/menudel'),(3142,1,'project/node'),(3143,1,'project/node/index'),(3144,1,'project/node/alllist'),(3145,1,'project/node/clear'),(3146,1,'project/node/save'),(3147,1,'project/notify'),(3148,1,'project/notify/index'),(3149,1,'project/notify/noreads'),(3150,1,'project/notify/setreadied'),(3151,1,'project/notify/batchdel'),(3152,1,'project/notify/read'),(3153,1,'project/notify/delete'),(3154,1,'project/organization'),(3155,1,'project/organization/index'),(3156,1,'project/organization/save'),(3157,1,'project/organization/read'),(3158,1,'project/organization/edit'),(3159,1,'project/organization/delete'),(3160,1,'project/project'),(3161,1,'project/project/index'),(3162,1,'project/project/selflist'),(3163,1,'project/project/save'),(3164,1,'project/project/read'),(3165,1,'project/project/edit'),(3166,1,'project/project/uploadcover'),(3167,1,'project/project/recycle'),(3168,1,'project/project/recovery'),(3169,1,'project/project/archive'),(3170,1,'project/project/recoveryarchive'),(3171,1,'project/project/quit'),(3172,1,'project/project_collect'),(3173,1,'project/project_collect/collect'),(3174,1,'project/project_member'),(3175,1,'project/project_member/index'),(3176,1,'project/project_member/searchinvitemember'),(3177,1,'project/project_member/invitemember'),(3178,1,'project/project_template'),(3179,1,'project/project_template/index'),(3180,1,'project/project_template/save'),(3181,1,'project/project_template/uploadcover'),(3182,1,'project/project_template/edit'),(3183,1,'project/project_template/delete'),(3184,1,'project/task'),(3185,1,'project/task/index'),(3186,1,'project/task/selflist'),(3187,1,'project/task/read'),(3188,1,'project/task/save'),(3189,1,'project/task/taskdone'),(3190,1,'project/task/assigntask'),(3191,1,'project/task/sort'),(3192,1,'project/task/createcomment'),(3193,1,'project/task/edit'),(3194,1,'project/task/like'),(3195,1,'project/task/star'),(3196,1,'project/task/recycle'),(3197,1,'project/task/recovery'),(3198,1,'project/task/delete'),(3199,1,'project/task_log'),(3200,1,'project/task_log/index'),(3201,1,'project/task_log/getlistbyselfproject'),(3202,1,'project/task_member'),(3203,1,'project/task_member/index'),(3204,1,'project/task_member/searchinvitemember'),(3205,1,'project/task_member/invitemember'),(3206,1,'project/task_member/invitememberbatch'),(3207,1,'project/task_stages'),(3208,1,'project/task_stages/index'),(3209,1,'project/task_stages/tasks'),(3210,1,'project/task_stages/sort'),(3211,1,'project/task_stages/save'),(3212,1,'project/task_stages/edit'),(3213,1,'project/task_stages/delete'),(3214,1,'project/task_stages_template'),(3215,1,'project/task_stages_template/index'),(3216,1,'project/task_stages_template/save'),(3217,1,'project/task_stages_template/edit'),(3218,1,'project/task_stages_template/delete'),(3219,2,'project/account/index'),(3220,2,'project/auth/index'),(3221,2,'project/index/index'),(3222,2,'project/index'),(3223,2,'project/index/changecurrentorganization'),(3224,2,'project/index/systemconfig'),(3225,2,'project/index/info'),(3226,2,'project/index/editpersonal'),(3227,2,'project/index/editpassword'),(3228,2,'project/index/uploadimg'),(3229,2,'project/index/uploadavatar'),(3230,2,'project/menu/menu'),(3231,2,'project/node/index'),(3232,2,'project/node/alllist'),(3233,2,'project/notify/index'),(3234,2,'project/notify'),(3235,2,'project/notify/noreads'),(3236,2,'project/notify/setreadied'),(3237,2,'project/notify/batchdel'),(3238,2,'project/notify/read'),(3239,2,'project/notify/delete'),(3240,2,'project/organization/index'),(3241,2,'project/organization'),(3242,2,'project/organization/save'),(3243,2,'project/organization/read'),(3244,2,'project/organization/edit'),(3245,2,'project/organization/delete'),(3246,2,'project/project/index'),(3247,2,'project/project/read'),(3248,2,'project/project_collect/collect'),(3249,2,'project/project_collect'),(3250,2,'project/project_member/index'),(3251,2,'project/project_template/index'),(3252,2,'project/task/index'),(3253,2,'project/task/read'),(3254,2,'project/task/save'),(3255,2,'project/task/taskdone'),(3256,2,'project/task/assigntask'),(3257,2,'project/task/sort'),(3258,2,'project/task/createcomment'),(3259,2,'project/task/like'),(3260,2,'project/task/star'),(3261,2,'project/task_log/index'),(3262,2,'project/task_log'),(3263,2,'project/task_log/getlistbyselfproject'),(3264,2,'project/task_member/index'),(3265,2,'project/task_member/searchinvitemember'),(3266,2,'project/task_stages/index'),(3267,2,'project/task_stages/tasks'),(3268,2,'project/task_stages/sort'),(3269,2,'project/task_stages_template/index'),(3270,2,'project/department/index'),(3271,2,'project/department/read'),(3272,2,'project/department_member/index'),(3273,2,'project/department_member/searchinvitemember'),(3274,2,'project/project/selflist'),(3275,2,'project/project/save'),(3276,2,'project/task/selflist');
/*!40000 ALTER TABLE `ms_project_auth_node` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_collection`
--

DROP TABLE IF EXISTS `ms_project_collection`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_collection` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_code` bigint DEFAULT '0' COMMENT '项目id',
  `member_code` bigint DEFAULT '0' COMMENT '成员id',
  `create_time` bigint DEFAULT '0' COMMENT '加入时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='项目-收藏表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_collection`
--

LOCK TABLES `ms_project_collection` WRITE;
/*!40000 ALTER TABLE `ms_project_collection` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_project_collection` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_log`
--

DROP TABLE IF EXISTS `ms_project_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `member_code` bigint DEFAULT '0' COMMENT '操作人id',
  `content` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci COMMENT '操作内容',
  `remark` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `type` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT 'create' COMMENT '操作类型',
  `create_time` bigint DEFAULT NULL COMMENT '添加时间',
  `source_code` bigint DEFAULT '0' COMMENT '任务id',
  `action_type` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '场景类型',
  `to_member_code` bigint DEFAULT '0',
  `is_comment` tinyint(1) DEFAULT '0' COMMENT '是否评论，0：否',
  `project_code` bigint DEFAULT NULL,
  `icon` varchar(20) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `is_robot` tinyint(1) DEFAULT '0' COMMENT '是否机器人',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `member_code` (`member_code`) USING BTREE,
  KEY `source_code` (`source_code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5086 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='项目日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_log`
--

LOCK TABLES `ms_project_log` WRITE;
/*!40000 ALTER TABLE `ms_project_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_project_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_member`
--

DROP TABLE IF EXISTS `ms_project_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_member` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_code` bigint DEFAULT NULL COMMENT '项目id',
  `member_code` bigint DEFAULT NULL COMMENT '成员id',
  `join_time` bigint DEFAULT NULL COMMENT '加入时间',
  `is_owner` bigint DEFAULT '0' COMMENT '拥有者',
  `authorize` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '角色',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique` (`project_code`,`member_code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='项目-成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_member`
--

LOCK TABLES `ms_project_member` WRITE;
/*!40000 ALTER TABLE `ms_project_member` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_project_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_menu`
--

DROP TABLE IF EXISTS `ms_project_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `pid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '父id',
  `title` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `icon` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '菜单图标',
  `url` varchar(400) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL DEFAULT '' COMMENT '链接',
  `file_path` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '文件路径',
  `params` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '' COMMENT '链接参数',
  `node` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT '#' COMMENT '权限节点',
  `sort` int unsigned DEFAULT '0' COMMENT '菜单排序',
  `status` tinyint unsigned DEFAULT '1' COMMENT '状态(0:禁用,1:启用)',
  `create_by` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建人',
  `is_inner` tinyint(1) DEFAULT '0' COMMENT '是否内页',
  `values` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '参数默认值',
  `show_slider` tinyint(1) DEFAULT '1' COMMENT '是否显示侧栏',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=176 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='项目菜单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_menu`
--

LOCK TABLES `ms_project_menu` WRITE;
/*!40000 ALTER TABLE `ms_project_menu` DISABLE KEYS */;
INSERT INTO `ms_project_menu` VALUES (120,0,'工作台','appstore-o','home','home',':org','#',0,1,0,0,'',0),(121,0,'项目管理','project','#','#','','#',0,1,0,0,'',1),(122,121,'项目列表','branches','#','#','','#',0,1,0,0,'',1),(124,0,'系统设置','setting','#','#','','#',100,1,0,0,'',1),(125,124,'成员管理','unlock','#','#','','#',10,1,0,0,'',1),(126,125,'账号列表','','system/account','system/account','','project/account/index',10,1,0,0,'',1),(127,122,'我的组织','','organization','organization','','project/organization/index',30,1,0,0,'',1),(130,125,'访问授权','','system/account/auth','system/account/auth','','project/auth/index',20,1,0,0,'',1),(131,125,'授权页面','','system/account/apply','system/account/apply',':id','project/auth/apply',30,1,0,1,'',1),(138,121,'消息提醒','info-circle-o','#','#','','#',30,1,0,0,'',1),(139,138,'站内消息','','notify/notice','notify/notice','','project/notify/index',0,1,0,0,'',1),(140,138,'系统公告','','notify/system','notify/system','','project/notify/index',10,1,0,0,'',1),(143,124,'系统管理','appstore','#','#','','#',0,1,0,0,'',1),(144,143,'菜单路由','','system/config/menu','system/config/menu','','project/menu/menuadd',0,1,0,0,'',1),(145,143,'访问节点','','system/config/node','system/config/node','','project/node/save',0,1,0,0,'',1),(148,124,'个人管理','user','#','#','','#',0,1,0,0,'',1),(149,148,'个人设置','','account/setting/base','account/setting/base','','project/index/editpersonal',0,1,0,0,'',1),(150,148,'安全设置','','account/setting/security','account/setting/security','','project/index/editpersonal',0,1,0,1,'',1),(151,122,'我的项目','','project/list','project/list',':type','project/project/index',0,1,0,0,'my',1),(152,122,'回收站','','project/recycle','project/recycle','','project/project/index',20,1,0,0,'',1),(153,121,'项目空间','heat-map','project/space/task','project/space/task',':code','#',20,1,0,1,'',1),(154,153,'任务详情','','project/space/task/:code/detail','project/space/taskdetail',':code','project/task/read',0,1,0,1,'',0),(155,122,'我的收藏','','project/list','project/list',':type','project/project/index',10,1,0,0,'collect',1),(156,121,'基础设置','experiment','#','#','','#',0,1,0,0,'',1),(157,156,'项目模板','','project/template','project/template','','project/project_template/index',0,1,0,0,'',1),(158,156,'项目列表模板','','project/template/taskStages','project/template/taskStages',':code','project/task_stages_template/index',0,1,0,1,'',0),(159,122,'已归档项目','','project/archive','project/archive','','project/project/index',10,1,0,0,'',1),(160,0,'团队成员','team','#','#','','#',0,1,0,1,'',0),(161,153,'项目概况','','project/space/overview','project/space/overview',':code','project/index/info',20,1,0,1,'',0),(162,153,'项目文件','','project/space/files','project/space/files',':code','project/index/info',10,1,0,1,'',0),(163,122,'项目分析','','project/analysis','project/analysis','','project/index/info',5,1,0,0,'',1),(164,160,'团队成员','','#','#','','#',0,1,0,1,'',0),(166,164,'团队成员','','members','members','','project/department/index',0,1,0,1,'',0),(167,164,'成员信息','','members/profile','members/profile',':code','project/department/read',0,1,0,1,'',0),(168,153,'版本管理','','project/space/features','project/space/features',':code','project/index/info',20,1,0,1,'',0);
/*!40000 ALTER TABLE `ms_project_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_node`
--

DROP TABLE IF EXISTS `ms_project_node`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_node` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `node` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '节点代码',
  `title` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '节点标题',
  `is_menu` tinyint unsigned DEFAULT '0' COMMENT '是否可设置为菜单',
  `is_auth` tinyint unsigned DEFAULT '1' COMMENT '是否启动RBAC权限控制',
  `is_login` tinyint unsigned DEFAULT '1' COMMENT '是否启动登录控制',
  `create_at` bigint DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `index_system_node_node` (`node`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=641 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='项目端节点表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_node`
--

LOCK TABLES `ms_project_node` WRITE;
/*!40000 ALTER TABLE `ms_project_node` DISABLE KEYS */;
INSERT INTO `ms_project_node` VALUES (360,'project','项目管理模块',0,1,1,1673277965322),(361,'project/index/info','详情',0,0,1,1673277965322),(362,'project/index','基础版块',0,1,1,1673277965322),(363,'project/index/index','框架布局',0,0,1,1673277965322),(364,'project/index/systemconfig','系统信息',0,0,0,1673277965322),(365,'project/index/editpersonal','修改个人资料',0,0,1,1673277965322),(366,'project/index/uploadavatar','上传头像',0,0,1,1673277965322),(370,'project/account','账号管理',0,1,1,1673277965322),(371,'project/account/index','账号列表',0,0,1,1673277965322),(372,'project/organization/index','组织列表',0,0,1,1673277965322),(373,'project/organization/save','创建组织',0,0,1,1673277965322),(374,'project/organization/read','组织信息',0,0,1,1673277965322),(375,'project/organization/edit','编辑组织',0,1,1,1673277965322),(376,'project/organization/delete','删除组织',0,1,1,1673277965322),(377,'project/organization','组织管理',0,1,1,1673277965322),(388,'project/auth/index','权限列表',0,0,1,1673277965322),(389,'project/auth/add','添加权限角色',0,1,1,1673277965322),(390,'project/auth/edit','编辑权限',0,1,1,1673277965322),(391,'project/auth/forbid','禁用权限',0,1,1,1673277965322),(392,'project/auth/resume','启用权限',0,1,1,1673277965322),(393,'project/auth/del','删除权限',0,1,1,1673277965322),(394,'project/auth','访问授权',0,1,1,1673277965322),(395,'project/auth/apply','应用权限',0,1,1,1673277965322),(396,'project/notify/index','通知列表',0,0,1,1673277965322),(397,'project/notify/noreads','未读通知',0,0,1,1673277965322),(399,'project/notify/read','通知信息',0,1,1,1673277965322),(401,'project/notify/delete','删除通知',0,1,1,1673277965322),(402,'project/notify','通知管理',0,1,1,1673277965322),(434,'project/account/auth','授权管理',0,1,1,1673277965322),(435,'project/account/add','添加账号',0,1,1,1673277965322),(436,'project/account/edit','编辑账号',0,1,1,1673277965322),(437,'project/account/del','删除账号',0,1,1,1673277965322),(438,'project/account/forbid','禁用账号',0,1,1,1673277965322),(439,'project/account/resume','启用账号',0,1,1,1673277965322),(498,'project/notify/setreadied','设置已读',0,1,1,1673277965322),(499,'project/notify/batchdel','批量删除',0,1,1,1673277965322),(500,'project/auth/setdefault','设置默认权限',0,1,1,1673277965322),(501,'project/department','部门管理',0,1,1,1673277965322),(502,'project/department/index','部门列表',0,0,1,1673277965322),(503,'project/department/read','部门信息',0,0,1,1673277965322),(504,'project/department/save','创建部门',0,1,1,1673277965322),(505,'project/department/edit','编辑部门',0,1,1,1673277965322),(506,'project/department/delete','删除部门',0,1,1,1673277965322),(507,'project/department_member','部门成员管理',0,1,1,1673277965322),(508,'project/department_member/index','部门成员列表',0,0,1,1673277965322),(509,'project/department_member/searchinvitemember','搜索部门成员',0,0,1,1673277965322),(510,'project/department_member/invitemember','添加部门成员',0,1,1,1673277965322),(511,'project/department_member/removemember','移除部门成员',0,1,1,1673277965322),(512,'project/index/changecurrentorganization','切换当前组织',0,0,1,1673277965322),(513,'project/index/editpassword','修改密码',0,1,1,1673277965322),(514,'project/index/uploadimg','上传图片',0,0,1,1673277965322),(515,'project/menu','菜单管理',0,1,1,1673277965322),(516,'project/menu/menu','菜单列表',0,0,0,1673277965322),(517,'project/menu/menuadd','添加菜单',0,1,1,1673277965322),(518,'project/menu/menuedit','编辑菜单',0,1,1,1673277965322),(519,'project/menu/menuforbid','禁用菜单',0,1,1,1673277965322),(520,'project/menu/menuresume','启用菜单',0,1,1,1673277965322),(521,'project/menu/menudel','删除菜单',0,1,1,1673277965322),(522,'project/node','节点管理',0,1,1,1673277965322),(523,'project/node/index','节点列表',0,1,1,1673277965322),(524,'project/node/alllist','全部节点列表',0,1,1,1673277965322),(525,'project/node/clear','清理节点',0,1,1,1673277965322),(526,'project/node/save','编辑节点',0,1,1,1673277965322),(527,'project/project','项目管理',0,1,1,1673277965322),(528,'project/project/index','项目列表',0,0,1,1673277965322),(529,'project/project/selflist','个人项目列表',0,0,1,1673277965322),(530,'project/project/save','创建项目',0,1,1,1673277965322),(531,'project/project/read','项目信息',0,0,1,1673277965322),(532,'project/project/edit','编辑项目',0,1,1,1673277965322),(533,'project/project/uploadcover','上传项目封面',0,0,1,1673277965322),(534,'project/project/recycle','项目放入回收站',0,1,1,1673277965322),(535,'project/project/recovery','恢复项目',0,1,1,1673277965322),(536,'project/project/archive','归档项目',0,1,1,1673277965322),(537,'project/project/recoveryarchive','取消归档项目',0,1,1,1673277965322),(538,'project/project/quit','退出项目',0,1,1,1673277965322),(539,'project/project_collect','项目收藏管理',0,0,1,1673277965322),(540,'project/project_collect/collect','收藏项目',0,1,1,1673277965322),(541,'project/project_member','项目成员管理',0,1,1,1673277965322),(542,'project/project_member/index','项目成员列表',0,0,1,1673277965322),(543,'project/project_member/searchinvitemember','搜索项目成员',0,0,1,1673277965322),(544,'project/project_member/invitemember','邀请项目成员',0,1,1,1673277965322),(545,'project/project_template','项目模板管理',0,1,1,1673277965322),(546,'project/project_template/index','项目模板列表',0,0,1,1673277965322),(547,'project/project_template/save','创建项目模板',0,1,1,1673277965322),(548,'project/project_template/uploadcover','上传项目模板封面',0,1,1,1673277965322),(549,'project/project_template/edit','编辑项目模板',0,1,1,1673277965322),(550,'project/project_template/delete','删除项目模板',0,1,1,1673277965322),(551,'project/task/index','任务列表',0,0,1,1673277965322),(552,'project/task/selflist','个人任务列表',0,0,1,1673277965322),(553,'project/task/read','任务信息',0,0,1,1673277965322),(554,'project/task/save','创建任务',0,1,1,1673277965322),(555,'project/task/taskdone','更改任务状态',0,0,1,1673277965322),(556,'project/task/assigntask','指派任务执行者',0,1,1,1673277965322),(557,'project/task/sort','任务排序',0,1,1,1673277965322),(558,'project/task/createcomment','发表任务评论',0,1,1,1673277965322),(559,'project/task/edit','编辑任务',0,1,1,1673277965322),(560,'project/task/like','点赞任务',0,0,1,1673277965322),(561,'project/task/star','收藏任务',0,0,1,1673277965322),(562,'project/task/recycle','移动任务到回收站',0,1,1,1673277965322),(563,'project/task/recovery','恢复任务',0,1,1,1673277965322),(564,'project/task/delete','删除任务',0,1,1,1673277965322),(565,'project/task','任务管理',0,1,1,1673277965322),(569,'project/task_member/index','任务成员列表',0,0,1,1673277965322),(570,'project/task_member/searchinvitemember','搜索任务成员',0,0,1,1673277965322),(571,'project/task_member/invitemember','添加任务成员',0,1,1,1673277965322),(572,'project/task_member/invitememberbatch','批量添加任务成员',0,1,1,1673277965322),(573,'project/task_member','任务成员管理',0,1,1,1673277965322),(574,'project/task_stages','任务分组管理',0,1,1,1673277965322),(575,'project/task_stages/index','任务分组列表',0,0,1,1673277965322),(576,'project/task_stages/tasks','任务分组任务列表',0,0,1,1673277965322),(577,'project/task_stages/sort','任务分组排序',0,1,1,1673277965322),(578,'project/task_stages/save','添加任务分组',0,1,1,1673277965322),(579,'project/task_stages/edit','编辑任务分组',0,1,1,1673277965322),(580,'project/task_stages/delete','删除任务分组',0,1,1,1673277965322),(581,'project/task_stages_template/index','任务分组模板列表',0,0,1,1673277965322),(582,'project/task_stages_template/save','创建任务分组模板',0,1,1,1673277965322),(583,'project/task_stages_template/edit','编辑任务分组模板',0,1,1,1673277965322),(584,'project/task_stages_template/delete','删除任务分组模板',0,1,1,1673277965322),(585,'project/task_stages_template','任务分组模板管理',0,1,1,1673277965322),(587,'project/project_member/removemember','移除项目成员',0,1,1,1673277965322),(588,'project/task/datetotalforproject','任务统计',0,0,1,1673277965322),(589,'project/task/tasksources','任务资源列表',0,0,1,1673277965322),(590,'project/file','文件管理',0,1,1,1673277965322),(591,'project/file/index','文件列表',0,0,1,1673277965322),(592,'project/file/read','文件详情',0,0,1,1673277965322),(593,'project/file/uploadfiles','上传文件',0,1,1,1673277965322),(594,'project/file/edit','编辑文件',0,1,1,1673277965322),(595,'project/file/recycle','文件移至回收站',0,1,1,1673277965322),(596,'project/file/recovery','恢复文件',0,1,1,1673277965322),(597,'project/file/delete','删除文件',0,1,1,1673277965322),(598,'project/project/getlogbyselfproject','项目概况',0,1,1,1673277965322),(599,'project/source_link','资源关联管理',0,1,1,1673277965322),(600,'project/source_link/delete','取消关联',0,1,1,1673277965322),(601,'project/task/tasklog','任务动态',0,1,1,1673277965322),(602,'project/task/recyclebatch','批量移动任务到回收站',0,1,1,1673277965322),(603,'project/invite_link','邀请链接管理',0,1,1,1673277965322),(604,'project/invite_link/save','创建邀请链接',0,1,1,1673277965322),(605,'project/task/setprivate','设置任务隐私模式',0,1,1,1673277965322),(606,'project/account/read','账号信息',0,1,1,1673277965322),(607,'project/task/batchassigntask','批量指派任务',0,1,1,1673277965322),(608,'project/task/tasktotags','任务标签',0,1,1,1673277965322),(609,'project/task/settag','设置任务标签',0,1,1,1673277965322),(610,'project/task_tag','任务标签管理',0,1,1,1673277965322),(611,'project/task_tag/index','任务标签列表',0,1,1,1673277965322),(612,'project/task_tag/save','创建任务标签',0,1,1,1673277965322),(613,'project/task_tag/edit','编辑任务标签',0,1,1,1673277965322),(614,'project/task_tag/delete','删除任务标签',0,1,1,1673277965322),(615,'project/project_features','项目版本库管理',0,1,1,1673277965322),(616,'project/project_features/index','版本库列表',0,1,1,1673277965322),(617,'project/project_features/save','添加版本库',0,1,1,1673277965322),(618,'project/project_features/edit','编辑版本库',0,1,1,1673277965322),(619,'project/project_features/delete','删除版本库',0,1,1,1673277965322),(620,'project/project_version','项目版本管理',0,1,1,1673277965322),(621,'project/project_version/index','项目版本列表',0,1,1,1673277965322),(622,'project/project_version/save','添加项目版本',0,1,1,1673277965322),(623,'project/project_version/edit','编辑项目版本',0,1,1,1673277965322),(624,'project/project_version/changestatus','更改项目版本状态',0,1,1,1673277965322),(625,'project/project_version/read','项目版本详情',0,1,1,1673277965322),(626,'project/project_version/addversiontask','关联项目版本任务',0,1,1,1673277965322),(627,'project/project_version/removeversiontask','移除项目版本任务',0,1,1,1673277965322),(628,'project/project_version/delete','删除项目版本',0,1,1,1673277965322),(629,'project/task/getlistbytasktag','标签任务列表',0,1,1,1673277965322),(630,'project/task_workflow','任务流转管理',0,1,1,1673277965322),(631,'project/task_workflow/index','任务流转列表',0,1,1,1673277965322),(632,'project/task_workflow/save','添加任务流转',0,1,1,1673277965322),(633,'project/task_workflow/edit','编辑任务流转',0,1,1,1673277965322),(634,'project/task_workflow/delete','删除任务流转',0,1,1,1673277965322),(635,'project/department_member/detail','部门成员详情',0,1,1,1673277965322),(636,'project/department_member/uploadfile','上传头像',0,1,1,1673277965322),(637,'project/task/savetaskworktime','保存任务流转',0,1,1,1673277965322),(638,'project/task/edittaskworktime','编辑任务流转',0,1,1,1673277965322),(639,'project/task/deltaskworktime','删除任务流转',0,1,1,1673277965322),(640,'project/task/uploadfile','上传文件',0,1,1,1673277965322);
/*!40000 ALTER TABLE `ms_project_node` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_project_template`
--

DROP TABLE IF EXISTS `ms_project_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_project_template` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '类型名称',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注',
  `sort` tinyint DEFAULT '0',
  `create_time` bigint DEFAULT '0',
  `organization_code` bigint DEFAULT NULL COMMENT '组织id',
  `cover` varchar(511) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '封面',
  `member_code` bigint DEFAULT NULL COMMENT '创建人',
  `is_system` tinyint(1) DEFAULT '0' COMMENT '系统默认',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='项目类型表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_project_template`
--

LOCK TABLES `ms_project_template` WRITE;
/*!40000 ALTER TABLE `ms_project_template` DISABLE KEYS */;
INSERT INTO `ms_project_template` VALUES (11,'产品进展','适用于互联网产品人员对产品计划、跟进及发布管理',0,1670904236057,17,'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fbpic.51yuansu.com%2Fpic3%2Fcover%2F01%2F91%2F92%2F5982adf6c88ea_610.jpg&refer=http%3A%2F%2Fbpic.51yuansu.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673496114&t=956c5614481fedea97794e161deddb00',NULL,1),(12,'需求管理','适用于产品部门对需求的收集、评估及反馈管理',0,1670904236057,17,'https://img0.baidu.com/it/u=437485064,4277010738&fm=253&fmt=auto&app=138&f=JPEG?w=610&h=491',NULL,1),(13,'机械制造','适用于制造商对图纸设计及制造安装的工作流程管理',0,1670904236057,17,'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fbpic.51yuansu.com%2Fpic2%2Fcover%2F00%2F38%2F93%2F5812ca7a24020_610.jpg&refer=http%3A%2F%2Fbpic.51yuansu.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673496114&t=6d03fb91b230058fc43f1b7ae00f73e3',NULL,1),(19,'OKR 管理','适用于团队的 OKR 管理',0,1670904236057,17,'https://img2.baidu.com/it/u=2241642503,1613686234&fm=253&fmt=auto&app=138&f=JPEG?w=603&h=500',1015,0);
/*!40000 ALTER TABLE `ms_project_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_source_link`
--

DROP TABLE IF EXISTS `ms_source_link`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_source_link` (
  `id` int NOT NULL AUTO_INCREMENT,
  `source_type` char(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '资源类型',
  `source_code` bigint DEFAULT NULL COMMENT '资源编号',
  `link_type` char(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '关联类型',
  `link_code` bigint DEFAULT NULL COMMENT '关联编号',
  `organization_code` bigint DEFAULT NULL COMMENT '组织编码',
  `create_by` bigint DEFAULT NULL COMMENT '创建人',
  `create_time` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '创建时间',
  `sort` int DEFAULT '0' COMMENT '排序',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='资源关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_source_link`
--

LOCK TABLES `ms_source_link` WRITE;
/*!40000 ALTER TABLE `ms_source_link` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_source_link` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_task`
--

DROP TABLE IF EXISTS `ms_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `project_code` bigint NOT NULL DEFAULT '0' COMMENT '项目编号',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `pri` tinyint unsigned DEFAULT '0' COMMENT '紧急程度',
  `execute_status` tinyint DEFAULT NULL COMMENT '执行状态',
  `description` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci COMMENT '详情',
  `create_by` bigint DEFAULT NULL COMMENT '创建人',
  `done_by` bigint DEFAULT NULL COMMENT '完成人',
  `done_time` bigint DEFAULT NULL COMMENT '完成时间',
  `create_time` bigint DEFAULT NULL COMMENT '创建日期',
  `assign_to` bigint DEFAULT NULL COMMENT '指派给谁',
  `deleted` tinyint(1) DEFAULT '0' COMMENT '回收站',
  `stage_code` int DEFAULT NULL COMMENT '任务列表',
  `task_tag` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '任务标签',
  `done` tinyint DEFAULT '0' COMMENT '是否完成',
  `begin_time` bigint DEFAULT NULL COMMENT '开始时间',
  `end_time` bigint DEFAULT NULL COMMENT '截止时间',
  `remind_time` bigint DEFAULT NULL COMMENT '提醒时间',
  `pcode` bigint DEFAULT NULL COMMENT '父任务id',
  `sort` int DEFAULT '0' COMMENT '排序',
  `like` int DEFAULT '0' COMMENT '点赞数',
  `star` int DEFAULT '0' COMMENT '收藏数',
  `deleted_time` bigint DEFAULT NULL COMMENT '删除时间',
  `private` tinyint(1) DEFAULT '0' COMMENT '是否隐私模式',
  `id_num` int DEFAULT '1' COMMENT '任务id编号',
  `path` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci COMMENT '上级任务路径',
  `schedule` int DEFAULT '0' COMMENT '进度百分比',
  `version_code` bigint DEFAULT '0' COMMENT '版本id',
  `features_code` bigint DEFAULT '0' COMMENT '版本库id',
  `work_time` int DEFAULT '0' COMMENT '预估工时',
  `status` tinyint DEFAULT '0' COMMENT '执行状态。0：未开始，1：已完成，2：进行中，3：挂起，4：测试中',
  PRIMARY KEY (`id`,`project_code`) USING BTREE,
  KEY `stage_code` (`stage_code`) USING BTREE,
  KEY `project_code` (`project_code`) USING BTREE,
  KEY `pcode` (`pcode`) USING BTREE,
  KEY `sort` (`sort`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12363 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='任务表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_task`
--

LOCK TABLES `ms_task` WRITE;
/*!40000 ALTER TABLE `ms_task` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_task_member`
--

DROP TABLE IF EXISTS `ms_task_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_task_member` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_code` bigint DEFAULT '0' COMMENT '任务ID',
  `is_executor` tinyint(1) DEFAULT '0' COMMENT '执行者',
  `member_code` bigint DEFAULT NULL COMMENT '成员id',
  `join_time` bigint DEFAULT NULL,
  `is_owner` tinyint(1) DEFAULT '0' COMMENT '是否创建人',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=273 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='任务-成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_task_member`
--

LOCK TABLES `ms_task_member` WRITE;
/*!40000 ALTER TABLE `ms_task_member` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_task_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_task_stages`
--

DROP TABLE IF EXISTS `ms_task_stages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_task_stages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '类型名称',
  `project_code` bigint DEFAULT NULL COMMENT '项目id',
  `sort` int DEFAULT '0' COMMENT '排序',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注',
  `create_time` bigint DEFAULT NULL COMMENT '创建时间',
  `deleted` tinyint(1) DEFAULT '0' COMMENT '删除标记',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='任务列表表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_task_stages`
--

LOCK TABLES `ms_task_stages` WRITE;
/*!40000 ALTER TABLE `ms_task_stages` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_task_stages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_task_stages_template`
--

DROP TABLE IF EXISTS `ms_task_stages_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_task_stages_template` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '类型名称',
  `project_template_code` int DEFAULT '0' COMMENT '项目id',
  `create_time` bigint DEFAULT NULL,
  `sort` int DEFAULT '0',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=84 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=COMPACT COMMENT='任务列表模板表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_task_stages_template`
--

LOCK TABLES `ms_task_stages_template` WRITE;
/*!40000 ALTER TABLE `ms_task_stages_template` DISABLE KEYS */;
INSERT INTO `ms_task_stages_template` VALUES (61,'待处理',19,1670904236057,1),(62,'进行中',19,1670904236057,0),(63,'已完成',19,1670904236057,0),(65,'协议签订',13,1670904236057,0),(66,'图纸设计',13,1670904236057,0),(67,'评审及打样',13,1670904236057,0),(68,'构件采购',13,1670904236057,0),(69,'制造安装',13,1670904236057,0),(70,'内部检验',13,1670904236057,0),(71,'验收',13,1670904236057,0),(72,'需求收集',12,1670904236057,0),(73,'评估确认',12,1670904236057,0),(74,'需求暂缓',12,1670904236057,0),(75,'研发中',12,1670904236057,0),(76,'内测中',12,1670904236057,0),(77,'通知用户',12,1670904236057,0),(78,'已完成&归档',12,1670904236057,0),(79,'产品计划',11,1670904236057,0),(80,'即将发布',11,1670904236057,0),(81,'测试',11,1670904236057,0),(82,'准备发布',11,1670904236057,0),(83,'发布成功',11,1670904236057,0);
/*!40000 ALTER TABLE `ms_task_stages_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ms_task_work_time`
--

DROP TABLE IF EXISTS `ms_task_work_time`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ms_task_work_time` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_code` bigint DEFAULT '0' COMMENT '任务ID',
  `member_code` bigint DEFAULT NULL COMMENT '成员id',
  `create_time` bigint DEFAULT NULL,
  `content` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '描述',
  `begin_time` bigint DEFAULT NULL COMMENT '开始时间',
  `num` int DEFAULT '0' COMMENT '工时',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=COMPACT COMMENT='任务工时表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ms_task_work_time`
--

LOCK TABLES `ms_task_work_time` WRITE;
/*!40000 ALTER TABLE `ms_task_work_time` DISABLE KEYS */;
/*!40000 ALTER TABLE `ms_task_work_time` ENABLE KEYS */;
UNLOCK TABLES;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-03-13 10:55:38
