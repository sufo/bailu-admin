/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 操作中间件
 */

package middleware

import (
	"bailu/app/domain/entity"
	"bailu/app/service/sys"
	"bailu/global/consts"
	"bailu/pkg/jwt"
	"bailu/pkg/log"
	"bailu/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var respPool sync.Pool

func init() {
	respPool.New = func() interface{} {
		return make([]byte, 1024)
	}
}

func OperationMiddleware(operSrv *sys.OperationService, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if operSrv == nil {
			return
		}
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}
		// 记录请求时间
		start := time.Now()
		var body []byte
		var err error
		if c.Request.Method != http.MethodGet {
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				log.L.Error("read body from request error:", zap.Error(err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		}

		record := entity.OperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			//Path:   c.Request.RequestURI,
			Path:  c.Request.URL.Path,
			Agent: c.Request.UserAgent(),
			Body:  string(body),
			//OperId:   operId,
			//OperName: operName,
		}

		// 上传文件时候 中间件日志进行裁断操作
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			if len(record.Body) > 1024 {
				// 截断
				newBody := respPool.Get().([]byte)
				copy(newBody, record.Body)
				record.Body = string(newBody)
				defer respPool.Put(newBody[:0])
			}
		}

		// 使用自定义 ResponseWriter
		crw := CustomResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = crw

		// 执行请求处理程序和其他中间件函数
		c.Next()

		ctx := c.Copy()
		defer func(_start time.Time, crw CustomResponseWriter) {
			var operId uint64
			var operName string
			//操作人
			user, exist := ctx.Get(consts.REQUEST_USER)
			if exist {
				onlineUser := user.(*entity.OnlineUserDto)
				operId = onlineUser.ID
				operName = onlineUser.Username
			} else {
				token, exist := ctx.Get(consts.REQ_TOKEN)
				if exist {
					operId, err = jwt.ParseUserID(token.(string))
					if err != nil {
						operId = 0
					}
				}
			}
			record.OperId = operId
			record.OperName = operName

			//操作模块,
			//如果需要记录操作对应的描述信息，那么这个描述需要自己添加，在每个需要记录的controller层增加
			// 通过c *gin.Context c.Set("title","操作描述信息")
			//Record.Title = c.Get("title")

			//地址
			record.Location = utils.GetAddr(record.Ip)

			// 记录回包内容和处理时间
			latency := time.Since(_start)
			if ctx.Errors != nil {
				record.Msg = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
			}
			record.Status = ctx.Writer.Status()
			record.Latency = latency
			//w := ctx.Writer.(CustomResponseWriter)
			record.Resp = crw.body.String()

			//逻辑响应码和message
			if record.Status == http.StatusOK {
				var result Response
				err = json.Unmarshal(crw.body.Bytes(), &result)
				if err == nil {
					record.RespCode = &result.Code
					record.Msg = result.Msg
				}
			}

			if strings.Contains(crw.Header().Get("Pragma"), "public") ||
				strings.Contains(crw.Header().Get("Expires"), "0") ||
				strings.Contains(crw.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
				strings.Contains(crw.Header().Get("Content-Type"), "application/force-download") ||
				strings.Contains(crw.Header().Get("Content-Type"), "application/octet-mq") ||
				strings.Contains(crw.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
				strings.Contains(crw.Header().Get("Content-Type"), "application/download") ||
				strings.Contains(crw.Header().Get("Content-Disposition"), "attachment") ||
				strings.Contains(crw.Header().Get("Content-Transfer-Encoding"), "binary") {
				if len(record.Resp) > 1024 {
					// 截断
					newBody := respPool.Get().([]byte)
					copy(newBody, record.Resp)
					record.Resp = string(newBody)
					defer respPool.Put(newBody[:0])
				}
			}
			if err := operSrv.Create(ctx.Request.Context(), &record); err != nil {
				log.L.Error("create operation record error:", zap.Error(err))
			}
		}(start, crw)
	}
}
