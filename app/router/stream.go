/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package router

import (
	"github.com/gin-gonic/gin"
	"bailu/app/api/admin/monitor"
	"bailu/app/middleware"
)

func (r *Router) RegisterStream(app *gin.Engine) {
	s := app.Group("/stream")
	s.Use(middleware.LocaleMiddleware(), monitor.StreamHeadersMiddleware(), r.Event.SSEMiddleware(r.TokenProvider))
	{
		s.GET("", r.Event.Stream)
	}

}
