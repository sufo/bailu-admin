/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 登录访问日志
 */

package monitor

import (
	"bailu/app/api/admin"
	"bailu/app/domain/dto"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	respErr "bailu/pkg/exception"
	"bailu/pkg/log"
	"bailu/utils/page"
	"github.com/gin-gonic/gin"
)

type LoginLogApi struct {
	LoginLogSrv *sys.LoginLogService
}

func NewLoginLogApi(loginLogSrv *sys.LoginLogService) *LoginLogApi {
	return &LoginLogApi{loginLogSrv}
}

// @title 登录日志列表
// @Summary 登录日志列表接口
// @Description 按条件查询登录日志列表
// @Tags LoginLog
// @Accept json
// @Produce json
// @Param query query dto.LoginLogParams false "查询条件"
// ///////@Success 200 {object} resp.Response{data=resp.PageResult[entity.LoginInfo]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.LoginInfo]]
// @Router /api/loginLog [get]
// @Security Bearer
func (l *LoginLogApi) Index(c *gin.Context) {
	var queryParams dto.LoginLogParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		panic(respErr.BadRequestError)
	}
	page.StartPage(c)
	result, err := l.LoginLogSrv.List(c.Request.Context(), queryParams)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 批量删除登录日志
// @Summary 批量删除登录日志接口
// @Description 批量删除登录日志
// @Tags LoginLog
// @Accept json
// @Produce json
// @Param ids path string true "ids"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/loginLog/{ids} [delete]
func (l *LoginLogApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	if err := l.LoginLogSrv.Destroy(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 清空登录日志
// @Summary 清空登录日志接口
// @Description 清空登录日志
// @Tags LoginLog
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} resp.Response[any]
// @Router /api/loginLog/clean [delete]
func (l *LoginLogApi) Clean(c *gin.Context) {
	if err := l.LoginLogSrv.Clean(c.Request.Context()); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title ID查询登录日志
// @Summary 根据ID查询登录日志接口
// @Description 按条件查询登录日志
// @Tags LoginLog
// @Accept json
// @Produce json
// @Success 200 {object} resp.Response[entity.LoginInfo]
// @Router /api/loginLog/{userId} [get]
// @Security Bearer
func (l *LoginLogApi) FindByName(c *gin.Context) {
	var user = admin.GetLoginUser(c)
	result, err := l.LoginLogSrv.FindByName(c.Request.Context(), user.Username)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}
