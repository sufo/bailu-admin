package status

import (
	"github.com/sufo/bailu-admin/app/config"
	"strings"
)

// HTTP status codes, defined in RFC 2616.
const (
	StatusContinue           = 100
	StatusSwitchingProtocols = 101

	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent            = 204
	StatusResetContent         = 205
	StatusPartialContent       = 206

	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	StatusTemporaryRedirect = 307

	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
	StatusPreconditionRequired         = 428
	StatusTooManyRequests              = 429
	StatusRequestHeaderFieldsTooLarge  = 431
	StatusUnavailableForLegalReasons   = 451

	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusNetworkAuthenticationRequired = 511
)

var statusText = map[int]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",

	StatusOK:                   "OK",
	StatusCreated:              "Created",
	StatusAccepted:             "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	StatusNoContent:            "No Content",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",

	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

var statusTextZH = map[int]string{
	StatusContinue:           "继续",
	StatusSwitchingProtocols: "切换协议",

	StatusOK:                   "请求成功",
	StatusCreated:              "已创建",
	StatusAccepted:             "已接受请求",
	StatusNonAuthoritativeInfo: "非授权信息",
	StatusNoContent:            "无内容",
	StatusResetContent:         "重置内容",
	StatusPartialContent:       "服务器成功处理了部分GET请求",

	StatusMultipleChoices:   "多种选择",
	StatusMovedPermanently:  "请求的资源已被永久的移动到新URI",
	StatusFound:             "临时移动资源",
	StatusSeeOther:          "查看其它地址",
	StatusNotModified:       "所请求的资源未修改",
	StatusUseProxy:          "请求的资源必须通过代理访问",
	StatusTemporaryRedirect: "临时重定向",

	StatusBadRequest:                   "客户端请求错误",
	StatusUnauthorized:                 "未认证",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "拒绝执行此请求",
	StatusNotFound:                     "未找到资源",
	StatusMethodNotAllowed:             "请求方法不允许",
	StatusNotAcceptable:                "客户端无法完成解析",
	StatusProxyAuthRequired:            "代理身份需要认证",
	StatusRequestTimeout:               "请求超时",
	StatusConflict:                     "请求冲突",
	StatusGone:                         "资源已经不存在",
	StatusLengthRequired:               "请求请携带Content-Length参数",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "请求实体过大",
	StatusRequestURITooLong:            "URI过长",
	StatusUnsupportedMediaType:         "不支持的媒体格式",
	StatusRequestedRangeNotSatisfiable: "请求的范围无效",
	StatusExpectationFailed:            "无法满足Expect的请求头信息",
	StatusTeapot:                       "I'm a teapot",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "请求太频繁",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	StatusInternalServerError:           "服务器内部错误",
	StatusNotImplemented:                "服务器不支持请求",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "服务器暂不可用",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "不支持请求的HTTP协议的版本",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

// 返回httpcode对应的 状态码描述信息
// 返回空字符串表示状态码 unknown
func StatusText(code int) string {
	locale := strings.ToLower(config.Conf.Server.Locale)
	if strings.HasPrefix(locale, "zh") {
		return statusTextZH[code]
	}
	return statusText[code]
}
