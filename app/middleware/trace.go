/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 设置trace id
 */

package middleware

import (
	"bailu/global/consts"
	"bailu/utils"
	"github.com/gin-gonic/gin"
)

func TraceMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		traceID := c.GetHeader(consts.REQUEST_ID_KEY)
		if traceID == "" {
			traceID = utils.NewTraceId()
		}
		//c.Writer.Header().Set("X-Trace-Id", traceID)
		c.Writer.Header().Set(consts.REQUEST_ID_KEY, traceID)
		c.Set(consts.REQUEST_ID_KEY, traceID)

		c.Next()
	}
}
