-- create user
-- CREATE USER 'username'@'host' IDENTIFIED BY 'password';
-- grant
-- GRANT privileges ON databasename.tablename TO 'username'@'host';
-- GRANT ALL ON `bailu-admin`.* TO 'username'@'%';


CREATE DATABASE  IF NOT EXISTS `bailu-admin` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `bailu-admin`;


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `file_category`
--

DROP TABLE IF EXISTS `file_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `file_category` (
                                 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                 `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文件名',
                                 `created_at` datetime DEFAULT NULL,
                                 `updated_at` datetime DEFAULT NULL,
                                 `deleted_at` datetime DEFAULT NULL,
                                 `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                                 `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `file_category`
--

LOCK TABLES `file_category` WRITE;
/*!40000 ALTER TABLE `file_category` DISABLE KEYS */;
/*!40000 ALTER TABLE `file_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `files` (
                         `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                         `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '地址',
                         `category_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '分类ID',
                         `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文件名',
                         `size` bigint DEFAULT NULL COMMENT '文件大小（KB）',
                         `origin_name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '原始文件名',
                         `mime` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文件类型',
                         `path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文件路径',
                         `tags` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '多tag以逗号分割',
                         `created_at` datetime DEFAULT NULL,
                         `updated_at` datetime DEFAULT NULL,
                         `deleted_at` datetime DEFAULT NULL,
                         `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                         `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `files`
--

LOCK TABLES `files` WRITE;
/*!40000 ALTER TABLE `files` DISABLE KEYS */;
/*!40000 ALTER TABLE `files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `msg_notice`
--

DROP TABLE IF EXISTS `msg_notice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `msg_notice` (
                              `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                              `type` tinyint DEFAULT NULL COMMENT ' 1通知，2公告(公告只能全体)',
                              `title` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '通知标题',
                              `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '内容',
                              `sender` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发布人',
                              `sender_id` bigint unsigned DEFAULT NULL COMMENT '发布人id,0表示是系统',
                              `receivers` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '接收者id(id类型由Group_type决定，多个以逗号分割)',
                              `send_scope` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '发送范围类型（all所有人、user指定用户、role角色、depart部门等），见字典message_scope_type',
                              `send_status` varchar(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '0' COMMENT '发布状态（0未发布，1已发布，2已撤销）',
                              `start_time` datetime DEFAULT NULL COMMENT '开始时间',
                              `end_time` datetime DEFAULT NULL COMMENT '结束时间',
                              `notify_channel` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'web' COMMENT '通知方式：web,app,mail,sms等',
                              `scheduled_time` datetime DEFAULT NULL COMMENT '指定发送时间。如果该通知/公告是特定时间发送的',
                              `send_time` datetime DEFAULT NULL COMMENT '发布时间',
                              `cancel_time` datetime DEFAULT NULL COMMENT '撤销时间',
                              `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                              `created_at` datetime DEFAULT NULL,
                              `deleted_at` datetime DEFAULT NULL,
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `msg_notice`
--

LOCK TABLES `msg_notice` WRITE;
/*!40000 ALTER TABLE `msg_notice` DISABLE KEYS */;
/*!40000 ALTER TABLE `msg_notice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `msg_notice_send`
--

DROP TABLE IF EXISTS `msg_notice_send`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `msg_notice_send` (
                                   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                   `msg_id` bigint unsigned DEFAULT NULL COMMENT '消息id',
                                   `receive_id` bigint unsigned DEFAULT NULL COMMENT '接收者id',
                                   `read_flag` tinyint DEFAULT '0' COMMENT '阅读状态（0未读，1已读, 2删除）',
                                   `read_time` datetime DEFAULT NULL COMMENT '查看时间或者删除时间',
                                   `created_at` datetime DEFAULT NULL,
                                   PRIMARY KEY (`id`),
                                   KEY `idx_msg_notice_send_msg_id` (`msg_id`),
                                   KEY `idx_msg_notice_send_receive_id` (`receive_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `msg_notice_send`
--

LOCK TABLES `msg_notice_send` WRITE;
/*!40000 ALTER TABLE `msg_notice_send` DISABLE KEYS */;
/*!40000 ALTER TABLE `msg_notice_send` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `msg_user_config`
--

DROP TABLE IF EXISTS `msg_user_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `msg_user_config` (
                                   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                   `switch_bit` bigint DEFAULT NULL COMMENT '位开关，第一位总开关，第二位：点赞，后面依次为评论、回复、@、关注',
                                   `notify_channel` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'web' COMMENT '通知方式：web,app,mail,sms等',
                                   `user_id` bigint unsigned DEFAULT NULL COMMENT '用户',
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `msg_user_config`
--

LOCK TABLES `msg_user_config` WRITE;
/*!40000 ALTER TABLE `msg_user_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `msg_user_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_config`
--

DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
                              `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                              `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '参数名称',
                              `key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '参数键名',
                              `value` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '参数值',
                              `type` tinyint(1) DEFAULT NULL COMMENT '1 系统类 2 业务类',
                              `status` tinyint DEFAULT '0' COMMENT '是否启用 (1:启用 2:禁用)',
                              `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
                              `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                              `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                              `created_at` datetime DEFAULT NULL,
                              `updated_at` datetime DEFAULT NULL,
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_config`
--

LOCK TABLES `sys_config` WRITE;
/*!40000 ALTER TABLE `sys_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dept`
--

DROP TABLE IF EXISTS `sys_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dept` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `pid` bigint unsigned DEFAULT '0' COMMENT '父部门ID',
                            `ancestors` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '上级部门列表，逗号隔开。dataScope会用到',
                            `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '部门名称',
                            `sort` bigint DEFAULT '0' COMMENT '排序',
                            `leader` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '部门领导',
                            `phone` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '联系电话',
                            `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮箱',
                            `status` tinyint DEFAULT NULL COMMENT '是否启用(1:启用 2:禁用)',
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                            `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dept`
--

LOCK TABLES `sys_dept` WRITE;
/*!40000 ALTER TABLE `sys_dept` DISABLE KEYS */;
INSERT INTO `sys_dept` VALUES (1,0,'','市场部',0,'张三','23456789','',1,'2025-11-14 12:51:22','2025-11-14 12:51:22',NULL,1,0),(2,0,'','研发部',0,'','123456','',1,'2026-01-17 17:05:19','2026-01-17 17:05:19',NULL,1,0),(3,0,'','产品部',0,'','123456','',1,'2026-01-17 17:06:15','2026-01-17 17:06:15',NULL,1,0),(4,0,'','人事部',0,'','123456','',1,'2026-01-17 17:06:39','2026-01-17 17:06:39',NULL,1,0);
/*!40000 ALTER TABLE `sys_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dict`
--

DROP TABLE IF EXISTS `sys_dict`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
                            `code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '字典类型/编码',
                            `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '字典描述',
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                            `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `uni_sys_dict_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict`
--

LOCK TABLES `sys_dict` WRITE;
/*!40000 ALTER TABLE `sys_dict` DISABLE KEYS */;
INSERT INTO `sys_dict` VALUES (1,'状态','status','状态','2025-12-13 06:46:06','2025-12-13 06:46:06',NULL,1,0),(2,'message_scope','message_scope_type','消息发送范围','2025-12-13 08:38:13','2025-12-13 08:38:13',NULL,1,0),(3,'notify_strategy','notify_strategy','通知策略，暂未实现','2026-01-18 06:40:57','2026-01-18 07:02:48',NULL,1,0);
/*!40000 ALTER TABLE `sys_dict` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dict_item`
--

DROP TABLE IF EXISTS `sys_dict_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict_item` (
                                 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                 `label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '字典标签',
                                 `value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '字典值',
                                 `code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '编码',
                                 `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
                                 `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认选中',
                                 `fixed` tinyint(1) DEFAULT '0' COMMENT '是否固定（固定的字典不提供编辑功能）',
                                 `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
                                 `status` tinyint DEFAULT '0' COMMENT '是否启用 (1:启用 2:禁用)',
                                 `created_at` datetime DEFAULT NULL,
                                 `updated_at` datetime DEFAULT NULL,
                                 `deleted_at` datetime DEFAULT NULL,
                                 `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                                 `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict_item`
--

LOCK TABLES `sys_dict_item` WRITE;
/*!40000 ALTER TABLE `sys_dict_item` DISABLE KEYS */;
INSERT INTO `sys_dict_item` VALUES (1,'role','role','message_scope_type','',0,0,1,1,'2025-12-13 08:38:48','2025-12-13 08:38:48',NULL,1,0),(2,'department','depart','message_scope_type','',0,0,1,1,'2025-12-13 10:09:25','2025-12-13 10:09:25',NULL,1,0),(3,'user','user','message_scope_type','',1,0,1,1,'2025-12-13 10:10:01','2025-12-13 10:37:42',NULL,1,0),(4,'不通知','1','notify_strategy','',0,0,1,1,'2026-01-18 06:58:19','2026-01-18 06:58:19',NULL,1,0),(5,'失败通知','2','notify_strategy','',0,0,1,1,'2026-01-18 06:59:01','2026-01-18 06:59:01',NULL,1,0),(6,'结束通知','3','notify_strategy','',0,0,1,1,'2026-01-18 07:00:24','2026-01-18 07:00:24',NULL,1,0),(7,'结果关键字匹配通知','4','notify_strategy','暂未实现',0,0,1,1,'2026-01-18 07:05:07','2026-01-18 07:05:07',NULL,1,0);
/*!40000 ALTER TABLE `sys_dict_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_log`
--

DROP TABLE IF EXISTS `sys_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_log` (
                           `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                           `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作模块',
                           `action` tinyint(1) DEFAULT NULL COMMENT '操作类型（业务类型（0其它 1新增 2修改 3删除））',
                           `method` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求方法',
                           `operator_type` tinyint(1) DEFAULT NULL COMMENT '操作类别（0其它 1后台用户 2手机端用户）',
                           `oper_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作人员姓名',
                           `dept_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作人员部门名称',
                           `oper_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作url',
                           `oper_ip` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作地址',
                           `oper_loc` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作地点',
                           `oper_param` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求参数',
                           `result` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '返回结果',
                           `status` tinyint(1) DEFAULT NULL COMMENT '操作状态（0正常 1异常）',
                           `error_msg` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '错误消息',
                           `oper_time` datetime DEFAULT NULL COMMENT '操作时间',
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_log`
--

LOCK TABLES `sys_log` WRITE;
/*!40000 ALTER TABLE `sys_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_login_info`
--

DROP TABLE IF EXISTS `sys_login_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_login_info` (
                                  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                  `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
                                  `ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '登录IP',
                                  `addr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '登录IP地址',
                                  `browser` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '浏览器',
                                  `os` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作系统',
                                  `status` tinyint DEFAULT NULL COMMENT '登录状态（0成功 1失败）',
                                  `msg` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '提示消息',
                                  `login_time` datetime DEFAULT NULL COMMENT '访问时间',
                                  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_login_info`
--

LOCK TABLES `sys_login_info` WRITE;
/*!40000 ALTER TABLE `sys_login_info` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_login_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_menu`
--

DROP TABLE IF EXISTS `sys_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menu` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `pid` bigint unsigned DEFAULT '0' COMMENT '父菜单ID',
                            `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单名称',
                            `i18n_key` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '国际化名称key',
                            `path` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '路由路径（链接地址）',
                            `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '组件路径',
                            `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '#' COMMENT '附加属性',
                            `query` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '附加属性',
                            `keep_alive` tinyint(1) DEFAULT '1' COMMENT '附加属性',
                            `is_frame` tinyint(1) DEFAULT '0' COMMENT '附加属性',
                            `hide` tinyint(1) DEFAULT '0' COMMENT '附加属性',
                            `permission` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '附加属性',
                            `type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '菜单类型(M目录 C菜单 F按钮）',
                            `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
                            `status` tinyint DEFAULT '0' COMMENT '是否启用 (1:启用 2:禁用)',
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                            `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menu`
--

LOCK TABLES `sys_menu` WRITE;
/*!40000 ALTER TABLE `sys_menu` DISABLE KEYS */;
INSERT INTO `sys_menu` VALUES (1,0,'仪表盘','','/dashboard','dashboard/index','carbon:dashboard',NULL,1,0,0,'','C',1,1,'2024-12-03 02:39:00','2024-12-16 09:52:14',NULL,0,0),(2,0,'系统管理','','/system','','ion:settings-outline',NULL,1,0,0,'','M',2,1,'2024-12-03 02:39:00','2026-01-21 14:08:28',NULL,0,0),(3,2,'用户管理','','user','sys/user/index','ic:round-manage-accounts',NULL,1,0,0,'sys:user:list','C',0,1,'2024-12-03 02:39:00','2024-12-03 02:39:00',NULL,0,0),(4,3,'查询','','','','#',NULL,1,0,0,'sys:user:query','F',1,1,'2024-12-03 02:39:00','2026-01-18 08:51:48',NULL,0,0),(5,3,'新增','','','','#',NULL,1,0,0,'sys:user:add','F',2,1,'2024-12-03 02:39:00','2026-01-21 14:11:12',NULL,0,0),(6,3,'修改','','','','#',NULL,1,0,0,'sys:user:edit','F',3,1,'2024-12-03 02:39:00','2026-01-21 14:11:23',NULL,0,0),(7,3,'删除','','','','#',NULL,1,0,0,'sys:user:remove','F',4,1,'2024-12-03 02:39:00','2026-01-21 14:11:29',NULL,0,0),(8,3,'导入','','','','#',NULL,1,0,0,'sys:user:import','F',6,1,'2024-12-03 02:39:00','2026-01-21 14:11:49',NULL,0,0),(9,3,'导出','','','','#',NULL,1,0,0,'sys:user:export','F',7,1,'2024-12-03 02:39:00','2026-01-21 14:11:57',NULL,0,0),(10,3,'重置密码','','','','#',NULL,1,0,0,'sys:user:resetPwd','F',5,1,'2024-12-03 02:39:00','2026-01-21 14:11:40',NULL,0,0),(11,2,'角色管理','','role','sys/role/index','carbon:user-role',NULL,1,0,0,'sys:role:list','C',0,1,'2024-12-03 02:39:00','2024-12-03 02:39:00',NULL,0,0),(12,11,'查询','','','','#',NULL,1,0,0,'sys:role:query','F',1,1,'2024-12-03 02:39:00','2026-01-17 18:08:09',NULL,0,0),(13,11,'新增','','','','#',NULL,1,0,0,'sys:role:add','F',1,1,'2024-12-03 02:39:00','2026-01-17 18:09:01',NULL,0,0),(14,11,'修改','','','','#',NULL,1,0,0,'sys:role:edit','F',1,1,'2024-12-03 02:39:00','2026-01-21 14:13:51',NULL,0,0),(15,11,'删除','','','','#',NULL,1,0,0,'sys:role:remove','F',1,1,'2024-12-03 02:39:00','2026-01-17 18:09:26',NULL,0,0),(16,11,'导出','','','','#',NULL,1,0,0,'sys:role:export','F',1,1,'2024-12-03 02:39:00','2026-01-21 14:14:04',NULL,0,0),(17,2,'菜单管理','','menu','sys/menu/index','clarity:tree-view-solid',NULL,1,0,0,'sys:menu:list','C',0,1,'2024-12-03 02:39:00','2024-12-03 02:39:00',NULL,0,0),(18,17,'查询','','',NULL,'#',NULL,1,0,0,'sys:menu:query','F',0,1,'2024-12-03 02:39:00','2024-12-03 02:39:00',NULL,0,0),(19,17,'新增','','','','#',NULL,1,0,0,'sys:menu:add','F',2,1,'2024-12-03 02:39:00','2026-01-21 14:12:37',NULL,0,0),(20,17,'修改','','','','#',NULL,1,0,0,'sys:menu:edit','F',3,1,'2024-12-03 02:39:00','2026-01-21 14:12:44',NULL,0,0),(21,17,'删除','','','','#',NULL,1,0,0,'sys:menu:remove','F',4,1,'2024-12-03 02:39:00','2026-01-21 14:12:50',NULL,0,0),(22,17,'导出','','','','#',NULL,1,0,0,'sys:menu:export','F',5,1,'2024-12-03 02:39:00','2026-01-21 14:13:00',NULL,0,0),(23,2,'部门管理','','dept','sys/dept/index','ant-design:apartment-outlined',NULL,1,0,0,'sys:dept:list','C',0,1,'2024-12-26 09:16:54','2024-12-26 09:58:27',NULL,0,0),(24,2,'岗位管理','','post','sys/post/index','lucide:id-card',NULL,1,0,0,'sys:post:list','C',0,1,'2024-12-26 09:56:45','2024-12-26 09:58:31',NULL,0,0),(25,2,'字典管理','','dict','sys/dict/index','lucide:book-a',NULL,1,0,0,'sys:dict:list','C',0,1,'2024-12-27 01:49:16','2024-12-27 08:36:00',NULL,0,0),(27,0,'系统监控','','monitor','','lucide:monitor',NULL,1,0,0,'','M',3,1,'2024-12-30 03:27:54','2026-01-21 14:08:35',NULL,0,0),(28,27,'在线用户','','online','monitor/online-user/index','bx:user-voice',NULL,1,0,0,'monitor:online:list','C',0,1,'2024-12-30 06:33:58','2024-12-30 06:33:58',NULL,0,0),(29,27,'定时任务','','crontab','monitor/crontab/index','carbon:event-schedule',NULL,1,0,0,'monitor:task:list','C',0,1,'2024-12-30 06:41:37','2024-12-30 06:41:37',NULL,0,0),(30,27,'服务监控','','server','monitor/server-info/index','lucide:server',NULL,1,0,0,'monitor:server:list','C',0,1,'2024-12-30 07:11:32','2024-12-30 07:13:29',NULL,0,0),(31,27,'登录日志','','login-log','monitor/login-log/index','ph:address-book-light',NULL,1,0,0,'monitor:log:list','C',0,1,'2024-12-30 07:52:58','2024-12-30 07:52:58',NULL,0,0),(32,27,'操作日志','','oper-log','monitor/oper-log/index','fluent-mdl2:text-document-edit',NULL,1,0,0,'monitor:operlog:list','C',0,1,'2024-12-30 09:23:29','2024-12-30 09:23:29',NULL,0,0),(35,23,'新增','','','','#',NULL,1,0,0,'sys:dept:add','F',1,1,'2026-01-17 17:23:01','2026-01-17 17:23:01',NULL,0,0),(36,23,'编辑','','','','#',NULL,1,0,0,'sys:dept:edit','F',1,1,'2026-01-17 17:49:49','2026-01-17 18:05:26',NULL,0,0),(37,23,'删除','','','','#',NULL,1,0,0,'sys:dept:remove','F',1,1,'2026-01-17 17:50:45','2026-01-17 17:51:30',NULL,0,0),(38,24,'新增','','','','#',NULL,1,0,0,'sys:post:add','F',1,1,'2026-01-17 17:52:27','2026-01-17 17:52:27',NULL,0,0),(39,24,'编辑','','','','#',NULL,1,0,0,'sys:post:edit','F',1,1,'2026-01-17 17:53:20','2026-01-17 18:05:15',NULL,0,0),(40,24,'删除','','','','#',NULL,1,0,0,'sys:post:remove','F',1,1,'2026-01-17 17:54:03','2026-01-17 17:54:57',NULL,0,0),(41,25,'新增','','','','#',NULL,1,0,0,'sys:dict:add','F',1,1,'2026-01-17 17:58:00','2026-01-17 17:58:00',NULL,0,0),(42,25,'配置','','','','#',NULL,1,0,0,'sys:dict:query','F',1,1,'2026-01-17 18:00:15','2026-01-17 18:00:15',NULL,0,0),(43,25,'编辑','','','','#',NULL,1,0,0,'sys:dict:edit','F',1,1,'2026-01-17 18:05:00','2026-01-17 18:05:00',NULL,0,0),(44,25,'删除','','','','#',NULL,1,0,0,'sys:dict:remove','F',1,1,'2026-01-17 18:07:20','2026-01-17 18:07:20',NULL,0,0),(45,11,'数据权限','','','','#',NULL,1,0,0,'sys:role:scope','F',1,1,'2026-01-17 18:10:28','2026-01-17 18:11:47',NULL,0,0),(46,28,'强退','','','','#',NULL,1,0,0,'sys:online:remove','F',1,1,'2026-01-17 18:14:20','2026-01-17 18:15:04',NULL,0,0),(47,31,'删除','','','','#',NULL,1,0,0,'monitor:loginlog:remove','F',1,1,'2026-01-17 18:21:29','2026-01-17 18:21:29',NULL,0,0),(48,31,'清空','','','','#',NULL,1,0,0,'monitor:loginlog:clear','F',1,1,'2026-01-17 18:22:41','2026-01-17 18:22:41',NULL,0,0),(49,29,'新增','','','','#',NULL,1,0,0,'monitor:task:add','F',1,1,'2026-01-18 07:07:15','2026-01-18 07:08:05',NULL,0,0),(50,29,'编辑','','','','#',NULL,1,0,0,'monitor:task:edit','F',1,1,'2026-01-18 07:09:11','2026-01-18 07:09:11',NULL,0,0),(51,29,'删除任务','','','','#',NULL,1,0,0,'monitor:task:remove','F',1,1,'2026-01-18 07:12:31','2026-01-18 07:12:31',NULL,0,0),(59,32,'详情','','','','#',NULL,1,0,0,'monitor:operlog:query','F',1,1,'2026-01-18 08:39:08','2026-01-18 08:39:08',NULL,0,0),(60,32,'删除','','','','#',NULL,1,0,0,'monitor:operlog:remove','F',1,1,'2026-01-18 08:39:46','2026-01-18 08:39:46',NULL,0,0),(61,32,'导出','','','','#',NULL,1,0,0,'monitor:oper:export','F',1,1,'2026-01-18 08:40:29','2026-01-18 08:40:29',NULL,0,0);
/*!40000 ALTER TABLE `sys_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_menu_api`
--

DROP TABLE IF EXISTS `sys_menu_api`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menu_api` (
                                `menu_id` bigint unsigned NOT NULL COMMENT '菜单id',
                                `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '请求方法',
                                `path` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '按钮对应的接口路径'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menu_api`
--

LOCK TABLES `sys_menu_api` WRITE;
/*!40000 ALTER TABLE `sys_menu_api` DISABLE KEYS */;
INSERT INTO `sys_menu_api` VALUES (35,'POST','/api/dept'),(35,'POST','/api/dept'),(37,'DELETE','/api/dept:deptId'),(38,'POST','/api/post'),(38,'POST','/api/post'),(40,'DELETE','/api/post:ids'),(41,'POST','/api/dictItem'),(41,'POST','/api/dictItem'),(42,'GET','/api/dictItemoptions'),(42,'GET','/api/dictItemoptions'),(43,'PUT','/api/dict'),(43,'PUT','/api/dict'),(39,'PUT','/api/post'),(39,'PUT','/api/post'),(36,'PUT','/api/dept'),(36,'PUT','/api/dept'),(44,'DELETE','/api/dict:codes'),(44,'DELETE','/api/dict:codes'),(12,'GET','/api/role:id'),(13,'POST','/api/role'),(15,'DELETE','/api/role:roleIds'),(45,'PATCH','/api/roledataScope'),(46,'DELETE','/api/online:ids'),(46,'DELETE','/api/online:ids'),(47,'DELETE','/api/loginLog:ids'),(47,'DELETE','/api/loginLog:ids'),(48,'DELETE','/api/loginLogclean'),(48,'DELETE','/api/loginLogclean'),(49,'POST','/api/task'),(49,'POST','/api/task'),(50,'PUT','/api/task'),(50,'PUT','/api/task'),(51,'DELETE','/api/task:ids'),(51,'DELETE','/api/task:ids'),(60,'DELETE','/api/oper:ids'),(60,'DELETE','/api/oper:ids'),(4,'GET','/api/userinfo'),(5,'POST','/api/user'),(6,'PUT','/api/userprofile'),(7,'DELETE','/api/user:userIds'),(10,'PATCH','/api/userresetPassword'),(19,'POST','/api/menu'),(20,'PATCH','/api/menu'),(21,'DELETE','/api/menu:menuId'),(14,'PUT','/api/role');
/*!40000 ALTER TABLE `sys_menu_api` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_oper_record`
--

DROP TABLE IF EXISTS `sys_oper_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_oper_record` (
                                   `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                   `ip` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求ip',
                                   `location` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '操作地点',
                                   `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求方法',
                                   `path` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '请求uri',
                                   `trace_id` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '追踪ID',
                                   `status` bigint DEFAULT NULL COMMENT 'http状态码',
                                   `resp_code` bigint DEFAULT NULL COMMENT ' 逻辑响应码',
                                   `latency` bigint DEFAULT NULL COMMENT '延迟',
                                   `agent` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '代理',
                                   `msg` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '响应信息',
                                   `body` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '请求Body',
                                   `resp` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '响应Body',
                                   `oper_id` bigint unsigned DEFAULT NULL COMMENT '用户id',
                                   `oper_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名称',
                                   `created_at` datetime DEFAULT NULL,
                                   `updated_at` datetime DEFAULT NULL,
                                   `deleted_at` datetime DEFAULT NULL,
                                   PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_oper_record`
--

LOCK TABLES `sys_oper_record` WRITE;
/*!40000 ALTER TABLE `sys_oper_record` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_oper_record` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_post`
--

DROP TABLE IF EXISTS `sys_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_post` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `post_code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                            `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '岗位名称',
                            `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
                            `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
                            `status` tinyint DEFAULT '0' COMMENT '是否启用 (1:启用 2:禁用)',
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                            `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_post`
--

LOCK TABLES `sys_post` WRITE;
/*!40000 ALTER TABLE `sys_post` DISABLE KEYS */;
INSERT INTO `sys_post` VALUES (1,'PM_01','产品经理','P6',1,1,'2026-01-17 17:13:57','2026-01-17 17:13:57',NULL,1,0),(2,'SALES_02','销售代表','',1,1,'2026-01-17 17:14:39','2026-01-17 17:14:39',NULL,1,0),(3,'MKT_OP','新媒体运营','',1,1,'2026-01-17 17:14:57','2026-01-17 17:14:57',NULL,1,0),(4,'ADMIN_01','行政专员','',1,1,'2026-01-17 17:15:15','2026-01-17 17:15:15',NULL,1,0);
/*!40000 ALTER TABLE `sys_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role`
--

DROP TABLE IF EXISTS `sys_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色名称',
                            `role_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色权限字符',
                            `data_scope` tinyint(1) DEFAULT '1' COMMENT '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5:仅本人）',
                            `menu_check_strictly` tinyint DEFAULT '0',
                            `dept_check_strictly` tinyint DEFAULT '0',
                            `sort` bigint unsigned DEFAULT '0' COMMENT '排序',
                            `status` tinyint DEFAULT '0' COMMENT '是否启用 (1:启用 2:禁用)',
                            `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                            `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role`
--

LOCK TABLES `sys_role` WRITE;
/*!40000 ALTER TABLE `sys_role` DISABLE KEYS */;
INSERT INTO `sys_role` VALUES (1,'超级管理员','admin',1,0,0,1,1,'超级管理员','2024-12-03 10:41:18',NULL,NULL,NULL,NULL),(3,'普通角色','common',5,1,0,2,1,'普通角色','2024-12-03 10:41:39','2026-01-21 13:54:49',NULL,NULL,NULL);
/*!40000 ALTER TABLE `sys_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_dept`
--

DROP TABLE IF EXISTS `sys_role_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_dept` (
                                 `role_id` bigint unsigned NOT NULL,
                                 `dept_id` bigint unsigned NOT NULL,
                                 PRIMARY KEY (`role_id`,`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_dept`
--

LOCK TABLES `sys_role_dept` WRITE;
/*!40000 ALTER TABLE `sys_role_dept` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_role_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_menu`
--

DROP TABLE IF EXISTS `sys_role_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_menu` (
                                 `role_id` bigint unsigned NOT NULL,
                                 `menu_id` bigint unsigned NOT NULL,
                                 PRIMARY KEY (`role_id`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_menu`
--

LOCK TABLES `sys_role_menu` WRITE;
/*!40000 ALTER TABLE `sys_role_menu` DISABLE KEYS */;
INSERT INTO `sys_role_menu` VALUES (3,1),(3,2),(3,3),(3,11),(3,17),(3,23),(3,24),(3,25);
/*!40000 ALTER TABLE `sys_role_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_task`
--

DROP TABLE IF EXISTS `sys_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_task` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '任务名称',
                            `group` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'DEFAULT' COMMENT '任务组名',
                            `protocol` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '执行方式 FUNC:函数 HTTP:http',
                            `cron_expression` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'cron执行表达式',
                            `invoke_target` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '调用目标',
                            `args` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '目标参数',
                            `http_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'get' COMMENT 'http 请求方式 get post put patch delete等',
                            `concurrent` tinyint unsigned DEFAULT '1' COMMENT '是否并发执行（1允许 2禁止）',
                            `status` tinyint unsigned DEFAULT '0' COMMENT '是否启用（1正常 2停用）',
                            `entry_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'job启动时返回的id',
                            `notify_strategy` tinyint unsigned DEFAULT NULL COMMENT '执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知',
                            `notify_channel` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'web' COMMENT '通知方式：web,app,mail,sms等',
                            `notify_keyword` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '通知匹配关键字(多个用,分割)',
                            `remark` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注',
                            `last_exec_time` datetime DEFAULT NULL COMMENT '最近一次执行时间',
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                            `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_task`
--

LOCK TABLES `sys_task` WRITE;
/*!40000 ALTER TABLE `sys_task` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_task_log`
--

DROP TABLE IF EXISTS `sys_task_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_task_log` (
                                `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                `task_id` bigint unsigned NOT NULL COMMENT '任务id',
                                `task_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务名称',
                                `task_group` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务组名',
                                `invoke_target` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '调用目标字符串',
                                `status` tinyint(1) DEFAULT '1',
                                `result` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '执行输出结果',
                                `except_info` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '异常信息',
                                `retry_times` bigint DEFAULT NULL COMMENT '重试次数',
                                `start_time` datetime DEFAULT NULL COMMENT '开始时间',
                                `stop_time` datetime DEFAULT NULL COMMENT '停止时间',
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_task_log`
--

LOCK TABLES `sys_task_log` WRITE;
/*!40000 ALTER TABLE `sys_task_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_task_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user`
--

DROP TABLE IF EXISTS `sys_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `username` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
                            `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户登录密码',
                            `nick_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户昵称',
                            `user_type` varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户类型 默认00：系统用户',
                            `profile` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户简介',
                            `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮箱',
                            `dial_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '86' COMMENT '地区（国家）编码',
                            `phone` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '手机号',
                            `sex` tinyint DEFAULT '0' COMMENT '0未知 1男 2女',
                            `avatar` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户头像',
                            `dept_id` bigint unsigned DEFAULT NULL COMMENT '部门ID',
                            `home_path` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户首页',
                            `ip` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户最后登录ip',
                            `status` tinyint(1) DEFAULT '0' COMMENT '是否启用(1:启用 2:禁用)',
                            `remark` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '备注信息',
                            `last_login_time` datetime DEFAULT NULL,
                            `created_at` datetime DEFAULT NULL,
                            `updated_at` datetime DEFAULT NULL,
                            `deleted_at` datetime DEFAULT NULL,
                            `create_by` bigint unsigned DEFAULT '0' COMMENT '创建者',
                            `update_by` bigint unsigned DEFAULT '0' COMMENT '更新者',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `uni_sys_user_username` (`username`),
                            UNIQUE KEY `uni_sys_user_email` (`email`),
                            UNIQUE KEY `uni_sys_user_phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user`
--

LOCK TABLES `sys_user` WRITE;
/*!40000 ALTER TABLE `sys_user` DISABLE KEYS */;
INSERT INTO `sys_user` VALUES (1,'bailu','$2a$10$Z2V4ugIxQFzgKWfat6HmzuAAEfDEkc5BMfXVpXvNeazLADweG5eOO','freddie','','','ouamour@gmail.com','86','13545055335',1,'imgs/ee26908bf9629eeb4b37dac350f4754a_20251231131234',NULL,'','127.0.0.1',1,'','2026-01-25 14:49:05','2023-04-27 14:32:21','2024-12-29 09:42:29',NULL,0,1),(2,'snufkin','$2a$10$FCl2Lr7KcUToPt..S16njuUDuSHfNe/otgC3Y6DJwg0.XChjv2Riq','snufkin',NULL,NULL,'snufkin@gmail.com','','13500001234',1,'imgs/ee26908bf9629eeb4b37dac350f4754a_20260123151848',NULL,NULL,'127.0.0.1',1,'普通用户123456','2026-01-24 15:16:51',NULL,'2026-01-22 13:09:45',NULL,0,2);
/*!40000 ALTER TABLE `sys_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_post`
--

DROP TABLE IF EXISTS `sys_user_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_post` (
                                 `user_id` bigint unsigned NOT NULL,
                                 `post_id` bigint unsigned NOT NULL,
                                 PRIMARY KEY (`user_id`,`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_post`
--

LOCK TABLES `sys_user_post` WRITE;
/*!40000 ALTER TABLE `sys_user_post` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_user_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_role`
--

DROP TABLE IF EXISTS `sys_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_role` (
                                 `user_id` bigint unsigned NOT NULL,
                                 `role_id` bigint unsigned NOT NULL,
                                 PRIMARY KEY (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_role`
--

LOCK TABLES `sys_user_role` WRITE;
/*!40000 ALTER TABLE `sys_user_role` DISABLE KEYS */;
INSERT INTO `sys_user_role` VALUES (1,1),(2,3);
/*!40000 ALTER TABLE `sys_user_role` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-01-25 23:40:37

