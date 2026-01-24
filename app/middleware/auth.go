/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc jwt验证
 */

package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/global/consts"
	"github.com/sufo/bailu-admin/pkg/jwt"
)

func AuthMiddleware(provider *jwt.JwtProvider, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.Conf.JWT.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		token := GetToken(c)
		if token == "" {
			resp.Unauthorized(c)
			c.Abort()
			return
		}
		//解析key
		userKey, err := jwt.ParseUserKey(token)
		if err != nil {
			resp.Unauthorized(c)
			c.Abort()
			return
		}
		res, err := provider.Store.Get(config.Conf.JWT.OnlineKey + userKey)
		if err != nil {
			resp.Unauthorized(c)
			c.Abort()
			return
		}
		var onlineUser = entity.OnlineUserDto{}
		err2 := json.Unmarshal([]byte(res), &onlineUser)
		if err2 != nil {
			resp.Unauthorized(c)
			c.Abort()
			return
		}

		//续期
		_ = provider.CheckRenewal(token)
		//缓存到gin Context中
		c.Set(consts.REQUEST_USER, &onlineUser)
		c.Set(consts.REQ_TOKEN, token)
		vvv, exist := c.Get(consts.REQUEST_USER)
		if exist {
			fmt.Printf("user: %v", vvv)
		} else {
			fmt.Printf("user not exist")
		}
		fmt.Println()

		//保存到全局context
		ctx := appctx.SetAuth(c.Request.Context(), &onlineUser)
		//为了在gorm的hook函数中使用
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}

}
