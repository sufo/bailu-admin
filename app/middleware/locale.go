/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"bailu/app/config"
	"bailu/pkg/i18n"
)

func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//lng := c.Query("lang")
		lang := c.Query("locale")
		if lang == "" {
			lang = c.Request.Header.Get("Accept-Language")
		}
		c.Set("i18n", lang)
		//设置全局配置
		config.Conf.Server.Locale = lang
		i18n.Default().DefaultLang = lang
		c.Next()
	}
}
