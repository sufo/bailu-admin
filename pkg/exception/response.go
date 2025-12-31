/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package respErr

import (
	"bailu/app/config"
	"bailu/app/domain/resp/status"
	"bailu/global/consts"
	"errors"
	"fmt"
	"net/http"
)

type ResponseError struct {
	Code   int //自定义错误码
	Msg    string
	Status int   // 响应状态码
	ERR    error // 响应错误
}

func (r ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Msg
}

func WrapResponse(err error, code, status int, msg string, args ...interface{}) ResponseError {
	return ResponseError{
		Code:   code,
		Msg:    fmt.Sprintf(msg, args...),
		ERR:    err,
		Status: status,
	}
}

func New(code, status int, msg string, args ...interface{}) ResponseError {
	return ResponseError{
		Code:   code,
		Msg:    fmt.Sprintf(msg, args...),
		Status: status,
	}
}

func WrapLogicResp(msg string, args ...interface{}) ResponseError {
	return ResponseError{
		Code:   consts.ERROR,
		Msg:    fmt.Sprintf(msg, args...),
		Status: http.StatusOK,
	}
}

var (
	BadRequestError,
	UserDisableError,
	NotFoundError,
	ForbiddenError,
	TooManyRequestsError,
	UnauthorizedError,
	InternalServerError ResponseError
)

// 放进函数里面是为了延后赋值，为了等待config初始化完成
func Initial() {
	BadRequestError = New(0, status.StatusBadRequest, status.StatusText(status.StatusBadRequest))

	UserDisableError = New(0, status.StatusBadRequest, "user Forbidden")

	NotFoundError = New(0, status.StatusNotFound, status.StatusText(status.StatusNotFound))

	ForbiddenError = New(0, status.StatusForbidden, status.StatusText(status.StatusForbidden))

	TooManyRequestsError = New(0, status.StatusTooManyRequests, status.StatusText(status.StatusTooManyRequests))

	UnauthorizedError = New(0, status.StatusUnauthorized, status.StatusText(status.StatusUnauthorized))

	InternalServerError = New(0, status.StatusInternalServerError, status.StatusText(status.StatusInternalServerError))
}

func InternalServerErrorWithMsg(msg string) ResponseError {
	return New(0, status.StatusInternalServerError, msg)
}
func InternalServerErrorWithError(err error) ResponseError {
	return WrapResponse(err, 0, status.StatusInternalServerError, "")
}
func BadRequestErrorWithMsg(msg string) ResponseError {
	return New(0, status.StatusBadRequest, msg)
}
func BadRequestErrorWithError(err error) ResponseError {
	//不是debug，则不显示明确报错信息
	if config.Conf.Server.Mode != "debug" {
		err = errors.New(status.StatusText(status.StatusBadRequest))
	}
	return WrapResponse(err, 0, status.StatusBadRequest, "")
}
