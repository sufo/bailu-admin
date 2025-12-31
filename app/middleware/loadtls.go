package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

/*
*
// 用https把这个中间件在router里面use一下就好
//大部分时候，都会在外层来个nginx作反向代理，此时可以由nginx来做https的解析工作，
//最后将处理后的请求转发给你的 http web server，此时你就不需要使用 TSL 服务了。
// 如果使用endless下ListenAndServeTLS,则不需要使用此中间件
*/
func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			// 如果出现错误，请不要继续
			fmt.Println(err)
			return
		}
		// 继续往下处理
		c.Next()
	}
}
