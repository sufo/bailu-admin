/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 采用Option模式
 *		Option和Opt的设计视为了优雅的传递初始化参数
 *	     ...Opt是为了链式传递参数
 *       也可以设计为直接传值，则不需要Option、Opt这些东西
 */

package resp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/resp/status"
	"github.com/sufo/bailu-admin/global/consts"
	"github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"net/http"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type Options struct {
	HttpStatus int
	Response[any]
}

type Opt func(option *Options)

func Result(c *gin.Context, opt ...Opt) {
	option := Options{
		http.StatusOK,
		Response[any]{consts.SUCCESS, i18n.DefTr("api.successMessage"), struct{}{}},
	}

	for _, fn := range opt {
		fn(&option)
	}
	c.JSON(option.HttpStatus, option.Response)
}

//func ResultWithStatus(status, code int, data interface{}, msg string, c *gin.Context) {
//	// 开始时间
//	c.JSON(int(status), Response{
//		code,
//		data,
//		msg,
//	})
//}

// 设置Http状态码
func Status(status int) Opt {
	return func(option *Options) { option.HttpStatus = status }
}

func Code(code int) Opt {
	return func(option *Options) { option.Code = code }
}

func Data(data interface{}) Opt {
	return func(option *Options) { option.Data = data }
}

func Msg(msg string) Opt {
	return func(option *Options) { option.Msg = msg }
}

// 成功
func Ok(c *gin.Context) { Result(c) }

// 成功并自定义msg
func OkWithMsg(c *gin.Context, msg string) {
	Result(c, Msg(msg))
}

func OKWithData(c *gin.Context, data interface{}) {
	Result(c, Data(data))
}

// 存在严重问题，nil传递给interface{}，导致nil变的不为nil了
//func OKWithListData(c *gin.Context, data interface{}) {
//	println(reflect.TypeOf(data))
//	if data != nil {
//		Result(c, Data(data))
//	} else {
//		Result(c, Data([]struct{}{}))
//	}
//}

// 失败
func Fail(c *gin.Context) {
	Result(c, Status(http.StatusOK), Code(consts.ERROR), Msg(i18n.DefTr("admin.operationFailed"))) //"操作失败"
}

// 失败
func FailWithStatus(c *gin.Context, status int) {
	Result(c, Status(status), Code(consts.ERROR), Msg(i18n.DefTr("admin.operationFailed"))) //"操作失败"
}

// 失败
func FailWithMsg(c *gin.Context, msg string) {
	Result(c, Status(http.StatusOK), Code(consts.ERROR), Msg(msg))
}

// 失败
func FailWithStatusAndMsg(c *gin.Context, status int, msg string) {
	Result(c, Status(status), Code(consts.ERROR), Msg(msg))
}

func Unauthorized(c *gin.Context) {
	Result(c, Status(status.StatusUnauthorized), Code(status.StatusUnauthorized), Msg(status.StatusText(status.StatusUnauthorized)))
}

func MethodNotAllowed(c *gin.Context) {
	Result(c, Status(status.StatusMethodNotAllowed), Code(status.StatusMethodNotAllowed), Msg(status.StatusText(status.StatusMethodNotAllowed)))
}

func BadRequest(c *gin.Context) {
	Result(c, Status(status.StatusBadRequest), Code(status.StatusBadRequest), Msg(status.StatusText(status.StatusBadRequest)))
}

func Forbidden(c *gin.Context) {
	Result(c, Status(status.StatusForbidden), Code(status.StatusForbidden), Msg(status.StatusText(status.StatusForbidden)))
}

func NotFound(c *gin.Context) {
	Result(c, Status(status.StatusNotFound), Code(status.StatusNotFound), Msg(status.StatusText(status.StatusNotFound)))
}

func TooManyRequests(c *gin.Context) {
	Result(c, Status(status.StatusTooManyRequests), Code(status.StatusTooManyRequests), Msg(status.StatusText(status.StatusTooManyRequests)))
}

func InternalServerError(c *gin.Context) {
	Result(c, Status(status.StatusInternalServerError), Code(status.StatusInternalServerError), Msg(status.StatusText(status.StatusInternalServerError)))
}

func FailWithError(c *gin.Context, err error) {
	if e, ok := err.(respErr.ResponseError); ok {
		Result(c, Status(e.Status), Code(e.Code), Msg(e.Error()))
	} else {
		FailWithMsg(c, err.Error())
	}
}

func FailWithErrorAndRecordLog(c *gin.Context, err error) {
	var msg = ""
	if e, ok := err.(respErr.ResponseError); ok {
		msg = e.Error()
		Result(c, Status(e.Status), Code(e.Code), Msg(msg))
	} else {
		msg = err.Error()
		FailWithMsg(c, msg)
	}

	if config.Conf.Server.Mode != consts.MODE_RELEASE {
		fmt.Printf("[INFO] Response: %s %s %s\n", c.Request.Method, c.Request.RequestURI, msg)
	}
}
