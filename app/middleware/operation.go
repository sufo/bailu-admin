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
	"mime"
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

// CustomResponseWriter 封装 gin ResponseWriter 用于获取回包内容。
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
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

		contentType := c.GetHeader("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			// For file uploads, only log form fields, not the file content.
			if err := c.Request.ParseMultipartForm(8 << 20); err == nil { // 8 MB max memory
				body, _ = json.Marshal(c.Request.PostForm)
			}
		} else if c.Request.Method != http.MethodGet {
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
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
		}

		// The original truncation logic is no longer needed here as we don't read the file content.

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
		
		            // Assign all record fields
		            record.OperId = operId
		            record.OperName = operName
		            record.TraceID = ctx.GetString(consts.REQUEST_ID_KEY)
		            record.Location = utils.GetAddr(record.Ip)
		            record.Latency = time.Since(_start)
		            record.Status = crw.Status()
		
		            if ctx.Errors != nil {
		                record.Msg = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		            }
		
		            // Conditional response logging
		            if ctx.Request.Method != http.MethodGet {
		                // For non-GET requests, log the response body
		                disposition := crw.Header().Get("Content-Disposition")
		                isAttachment := strings.Contains(disposition, "attachment")
		                isBinary := strings.Contains(crw.Header().Get("Content-Type"), "application/octet-stream") ||
		                    strings.Contains(crw.Header().Get("Content-Type"), "application/force-download") ||
		                    strings.Contains(crw.Header().Get("Content-Type"), "application/download")
		
		                if isAttachment || isBinary {
		                    filename := ""
		                    // Try to parse filename from Content-Disposition
		                    _, params, err := mime.ParseMediaType(disposition)
		                    if err == nil {
		                        filename = params["filename"]
		                    }
		                    fileSize := crw.Header().Get("Content-Length")
		                    fileInfo := map[string]string{
		                        "message":  "File Download",
		                        "filename": filename,
		                        "size":     fileSize,
		                    }
		                    respBytes, _ := json.Marshal(fileInfo)
		                    record.Resp = string(respBytes)
		                } else {
		                    respBody := crw.body.String()
		                    if len(respBody) > 1024 {
		                        record.Resp = respBody[:1024] + "..." // Truncate long responses
		                    } else {
		                        record.Resp = respBody
		                    }
		                }
		            }
		
		            // For all requests, try to unmarshal for business code and message
		            if record.Status == http.StatusOK && len(crw.body.Bytes()) > 0 {
		                var result Response
		                err = json.Unmarshal(crw.body.Bytes(), &result)
		                if err == nil {
		                    record.RespCode = &result.Code
		                    // Only override msg if it's not a file download
		                    if record.Msg == "" || result.Msg != "" {
		                        record.Msg = result.Msg
		                    }
		                }
		            }
		
		            if err := operSrv.Create(ctx.Request.Context(), &record); err != nil {
		                log.L.Error("create operation record error:", zap.Error(err))
		            }
		        }(start, crw)	}
}
