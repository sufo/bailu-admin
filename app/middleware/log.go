/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

//func LoggerMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if SkipHandler(c, skippers...) {
//			c.Next()
//			return
//		}
//		start := time.Now()
//		path := c.Request.URL.Path
//		c.Next()
//		cost := time.Since(start)
//		log.L.Info(path,
//			zap.Int("Status", c.Writer.Status()),
//			zap.String("Method", c.Request.Method),
//			zap.String("IP", c.ClientIP()),
//			zap.String("Path", path),
//			//zap.String("TraceId", uuidStr),
//			//zap.Int("UserId", userId),
//			zap.String("query", c.Request.URL.RawQuery),
//			zap.String("UserAgent", c.Request.UserAgent()),
//			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
//			zap.Duration("Cost", cost),
//		)
//	}
//}

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// 日志中间件。
func LoggerMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		// 记录请求时间
		start := time.Now()

		var requestBodyBytes []byte
		if c.Request.Body != nil {
			//requestBodyBytes, _ = io.ReadAll(c.Request.Body)
			requestBodyBytes, _ = c.GetRawData()
			// 打印请求信息
			fmt.Printf("[INFO] Request: %s %s %s\n", c.Request.Method, c.Request.RequestURI, requestBodyBytes)
		}
		// 将原body塞回去
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		// 使用自定义 ResponseWriter
		crw := CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw

		// 执行请求处理程序和其他中间件函数
		c.Next()

		// 记录回包内容和处理时间
		end := time.Now()
		latency := end.Sub(start)
		respBody := string(crw.body.Bytes())
		fmt.Printf("[INFO] Response: %s %s %s (%v)\n", c.Request.Method, c.Request.RequestURI, respBody, latency)
	}
}
