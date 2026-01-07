/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package middleware

import (
	"github.com/sufo/bailu-admin/app/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	cfg := config.Conf.CORS
	return cors.New(cors.Config{
		//AllowAllOrigins: true,
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Second * time.Duration(cfg.MaxAge),
	})
	//return func(c *gin.Context) {
	//	method := c.Request.Method
	//	origin := c.Request.Header.Get("Origin")
	//	c.Header("Access-Control-Allow-Origin", origin)
	//	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
	//	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
	//	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	//	c.Header("Access-Control-Allow-Credentials", "true")
	//
	//	// 放行所有OPTIONS方法
	//	if method == "OPTIONS" {
	//		c.AbortWithStatus(http.StatusNoContent)
	//	}
	//	// 处理请求
	//	c.Next()
	//}
}
