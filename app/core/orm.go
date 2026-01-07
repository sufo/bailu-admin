/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package core

import (
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/pkg/gormx"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

//@description: 初始化数据库并产生数据库全局变量
// gorm维护了一个连接池，初始化db之后所有的连接都由库来管理。所以不需要使用者手动关闭。
//@return: *gorm.DB

func InitGorm() (*gorm.DB, func(), error) {
	cfg := config.Conf.DataSource
	db, err := Gorm()
	if err != nil {
		return nil, nil, err
	}
	cleanFunc := func() {}

	if cfg.EnableAutoMigrate {
		err2 := AutoMigrate(db)
		if err2 != nil {
			return nil, cleanFunc, err2
		}
	}
	//registerCallback(db)
	return db, cleanFunc, nil
}

func Gorm() (*gorm.DB, error) {
	dbType := config.Conf.DataSource.DbType
	switch dbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}
}

// AutoMigrate
// @function: MysqlTables
// @description: 注册数据库表专用
// @param: db *gorm.DB
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 系统模块表
		&entity.User{},
		&entity.Dept{},
		&entity.Role{},
		&entity.Post{},
		&entity.Menu{},
		&entity.MenuApi{},
		//&entity.Parameter{},
		&entity.Dict{},
		&entity.DictItem{},
		&entity.Notice{},
		&entity.Log{},
		&entity.Task{},
		&entity.TaskLog{},
		&entity.OperationRecord{},

		&entity.LoginInfo{},
		&entity.SysConfig{},

		//通知公告、事件提醒、私信
		&entity.Notice{},
		&entity.NoticeSend{},
		&entity.RemindSubscription{},
		&entity.MsgUserConfig{},

		//文件
		&entity.FileInfo{},
		&entity.FileCategory{},
	)
	//if err != nil {
	//	log.L.Error("register table failed", zap.Any("err", err))
	//	os.Exit(0)
	//}
	//log.L.Info("register table success")
}

func GormMysql() (*gorm.DB, error) {
	dsc := config.Conf.DataSource
	if dsc.Mysql.DbName == "" {
		return nil, errors.New("未配置DbName")
	}
	dsn := dsc.Mysql.Dsn()
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  //禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig()); err != nil {
		return nil, err
	} else {
		//if dsc.Debug {
		//	db = db.Debug()
		//}
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(dsc.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dsc.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(dsc.MaxLifetime) * time.Second)
		return db, nil
	}
}

// @description: 根据配置决定是否开启日志
// @param: mod bool
// @return: *gorm.Config
func gormConfig() *gorm.Config {
	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Conf.DataSource.TablePrefix,
			SingularTable: true, //(遵守单数形式)
		},
		Logger: gormx.Default.LogMode(logger.Info), //设置日志级别，可以作为配置
		//GORM creates database foreign key constraints automatically when AutoMigrate or CreateTable,
		//disable this by setting it to true, refer Migration for details
		DisableForeignKeyConstraintWhenMigrating: true, //true禁止创建外键约束
	}
}

// Before After
// gorm:create gorm:delete gorm:query gorm:update
//func registerCallback(db *gorm.DB) {
//	// 自动添加updateBy CreateBy
//	err := db.Callback().Create().Before("gorm:create").Register("create_by", func(db *gorm.DB) {
//		db.Statement.SetColumn("create_by", )
//	})
//	err := db.Callback().Create().Before("gorm:update").Register("update_by", func(db *gorm.DB) {
//		db.Statement.SetColumn("update_by")
//	})
//	if err != nil {
//		log.L.Panicf("err: %+v", err)
//	}
//}
