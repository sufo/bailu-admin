/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc casbin权限处理
 */

package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/log"
	"gorm.io/gorm/utils"
	"strconv"
)

// Define permission equivalences. If the primary check for 'key' fails,
// Casbin will also check for 'value' as an alternative.
// This handles cases where different API endpoints fundamentally grant
// access to the same underlying resource data, e.g., a tree view vs. a list view.
var permissionEquivalence = map[string]string{
	"/api/dept/tree": "/api/dept", // Access to dept tree is also granted by dept list permission
	"api/menu/tree":  "/api/menu",
	// Add other equivalences here as needed:
	// "/api/user/options": "/api/user", // e.g., user dropdown might need general user list permission
}

func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skipperFunc ...SkipperFunc) gin.HandlerFunc {
	cfg := config.Conf.Casbin
	if !cfg.Enable {
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
		user, ok := val.(*entity.OnlineUserDto)
		if !ok {
			panic(respErr.ForbiddenError)
		}

		e, err := enforcer.Enforce(utils.ToString(user.ID), p, m)
		if err != nil {
			if cfg.Debug {
				resp.FailWithError(c, err)
				return
			} else {
				log.L.Errorf("casbin err is: %v", err)
				panic(respErr.InternalServerError)
			}
		}

		// If the initial Casbin check fails, try checking for an equivalent permission.
		if !e {
			if equivalentPath, ok := permissionEquivalence[p]; ok {
				log.L.Debugf("Casbin initial check failed for %s %s. Checking equivalent path: %s", m, p, equivalentPath)
				e, err = enforcer.Enforce(strconv.FormatUint(user.ID, 10), equivalentPath, m)
				if err != nil {
					if cfg.Debug {
						resp.FailWithError(c, err)
						return
					} else {
						log.L.Errorf("casbin equivalent path check err is: %v", err)
						panic(respErr.InternalServerError)
					}
				}
				if e {
					log.L.Debugf("Casbin check passed for equivalent path: %s %s", m, equivalentPath)
				}
			}
		}
		if !e {
			panic(respErr.ForbiddenError)
		}
		c.Next()
	}
}
