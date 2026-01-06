-- create user
-- CREATE USER 'username'@'host' IDENTIFIED BY 'password';
-- grant
-- GRANT privileges ON databasename.tablename TO 'username'@'host';
-- GRANT ALL ON `bailu-admin`.* TO 'test'@'%';


-- Create a database
CREATE DATABASE IF NOT EXISTS `bailu-admin` DEFAULT CHARACTER SET = `utf8mb4`;
-- CREATE SCHEMA `bailu-admin` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_ci ;

-- 删除所有表
-- 产生删除表的sql语句：
-- select concat("DROP TABLE IF EXISTS ", table_name, ";") from information_schema.tables where table_schema="bailu-admin"
use `bailu-admin`;
-- 然后执行上面生成的语句


-- 初始化超级管理员
insert into `sys_role` values(1, '超级管理员', 'admin', 1, 1, 1, '超级管理员',sysdate(),null, null, null,null);
insert into `sys_role` values(3, '普通角色', 'common',  5, 2, 1, '普通角色',sysdate(),null, null, null,null);
# insert into sys_role values('2', '普通角色',    'common', 2, 2, 1, 1, '0', '0', 'admin', sysdate(), '', null, '普通角色');

-- 初始化菜单
#
# insert into `sys_menu` values(2, 0, 'system', 'system', 'Layout', '系统管理','carbon:dashboard','',false,false,true,false,false,'M','',2,true,sysdate(),null, null, null,null);
-- ----------------------------
-- 初始化-用户信息表数据
-- ----------------------------
insert into `sys_user` values(1, 'bailu', '$2a$10$Z2V4ugIxQFzgKWfat6HmzuAAEfDEkc5BMfXVpXvNeazLADweG5eOO', '', '', '', '', '86', '13500000000', '0', '', NULL, '', '127.0.0.1', '1','', null, '2023-04-27 14:32:21', '2023-04-27 14:32:21', NULL, 0, 0);

-- 角色用户中间表
insert into `sys_user_role` values(1,1);
-- 角色菜单中间表
# insert into `sys_role_menu` values(1,2)
