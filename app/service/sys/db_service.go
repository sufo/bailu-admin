/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 数据库
 */

package sys

import "bailu/app/config"

type DBService struct {
}

// 创建数据库并初始化 总入口
func (dbService *DBService) InitMysqlDB(conf config.Mysql) error {
	panic("TODO")
}
