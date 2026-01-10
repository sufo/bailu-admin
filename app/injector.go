/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc injector
 */

package app

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/router"
	"github.com/sufo/bailu-admin/app/service/cron"
	"github.com/sufo/bailu-admin/app/service/sys"
	"go.uber.org/zap"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	Engine         *gin.Engine
	Logger         *zap.SugaredLogger
	Router         router.IRouter //如果前面也使用的IRouter，这里就不能使用 *router.Router
	MenuSrv        *sys.MenuService
	Job            *cron.CronTask
	Injector2Job   *cron.Inject2Jobs
	CasbinEnforcer *casbin.SyncedEnforcer
}
