/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc casbin权限处理
 */

package middleware

import (
	"bailu/app/config"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/global/consts"
	respErr "bailu/pkg/exception"
	"bailu/pkg/log"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skipperFunc ...SkipperFunc) gin.HandlerFunc {
	cfg := config.Conf.Casbin
	if cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skipperFunc...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		val, _ := c.Get(consts.REQUEST_USER)
		user := val.(entity.OnlineUserDto)
		e, err := enforcer.Enforce(user.ID, p, m)
		if err != nil {
			if cfg.Debug {
				resp.FailWithError(c, err)
				return
			} else {
				log.L.Errorf("casebin err is ", err)
				panic(respErr.InternalServerError)
			}
		}
		if !e {
			panic(respErr.ForbiddenError)
		}
		c.Next()
	}
}
