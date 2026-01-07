/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/pkg/i18n"
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
