/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/api/admin/monitor"
	"github.com/sufo/bailu-admin/app/middleware"
)

func (r *Router) RegisterStream(app *gin.Engine) {
	s := app.Group("/stream")
	s.Use(middleware.LocaleMiddleware(), monitor.StreamHeadersMiddleware(), r.Event.SSEMiddleware(r.TokenProvider))
	{
		s.GET("", r.Event.Stream)
	}

}
