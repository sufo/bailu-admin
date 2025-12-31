/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 初始化用户数据库
 */

package system

import (
	"bailu/app/config"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	"bailu/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var DBApiSet = wire.NewSet(wire.Struct(new(InitDBApi), "*"))

type InitDBApi struct {
	DB        *gorm.DB
	dbService sys.DBService
}

// TODO
func (d *InitDBApi) InitDB(c *gin.Context) {
	if d.DB != nil {
		log.L.Error("已存在数据库配置!")
		resp.FailWithMsg(c, "已存在数据库配置")
		return
	}
	if config.Conf.DataSource.DbType == "mysql" {
		var dbInfo = config.Conf.DataSource.Mysql
		if err := c.ShouldBindJSON(&dbInfo); err != nil {
			log.L.Error("参数校验不通过!", zap.Error(err))
			resp.FailWithMsg(c, "参数校验不通过")
			return
		}

		if err := d.dbService.InitMysqlDB(dbInfo); err != nil {
			log.L.Error("自动创建数据库失败!", zap.Error(err))
			resp.FailWithMsg(c, "自动创建数据库失败，请查看后台日志，检查后在进行初始化")
			return
		}
	}
	resp.OkWithMsg(c, "自动创建数据库成功")
}
