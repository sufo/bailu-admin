/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 中间件
 *		recover掉项目可能出现的panic，并使用zap记录相关日志
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

func RecoveryMiddleware(stack bool, trans ut.Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			//自定义错误
			if err := recover(); err != nil {
				if e, ok := err.(respErr.ResponseError); ok {
					resp.FailWithErrorAndRecordLog(c, e)
					c.Abort()
					return
				}
				//处理MySQLError类型错误
				if e, ok := err.(*mysql.MySQLError); ok {
					switch e.Number {
					case 1062: //Duplicate entry '%s' for key %s
						params := utils.REGEXP_1062.FindStringSubmatch(e.Message)
						if params != nil && len(params) > 1 {
							resp.FailWithMsg(c, i18n.DefTr("admin.existed", params[1]))
						} else {
							resp.FailWithErrorAndRecordLog(c, respErr.InternalServerError)
						}
					default:
						resp.FailWithErrorAndRecordLog(c, respErr.InternalServerError)
					}
					c.Abort()
					return
				}
				// 获取validator.ValidationErrors类型的errors
				if e2, ok := err.(validator.ValidationErrors); ok {
					var eInfo string
					if trans != nil {
						// validator.ValidationErrors类型错误则进行翻译
						// 并使用removeTopStruct函数去除字段名中的结构体名称标识
						errInfos := utils.RemoveTopStruct(e2.Translate(trans))
						//所有里面取出一个错误
						for _, v := range errInfos {
							eInfo = v
							break
						}
					} else {
						eInfo = e2.Error()
					}
					resp.FailWithMsg(c, eInfo)
					c.Abort()
					return
				}

				//Check for a broken connection, as it is not really a
				//condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.L.Error(c.Request.URL.Path,
						zap.Any("exception", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}
				if stack {
					log.L.Error("[Recovery from panic]",
						zap.Any("exception", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.L.Error("[Recovery from panic]",
						zap.Any("exception", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
