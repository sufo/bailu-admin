/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package monitor

import (
	"bailu/app/api/admin"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	respErr "bailu/pkg/exception"
	"bailu/pkg/log"
	"bailu/utils/page"
	"github.com/gin-gonic/gin"
)

type OnlineUserApi struct {
	OnlineUserSrv *sys.OnlineService
}

func NewOnlineUserApi(srv *sys.OnlineService) *OnlineUserApi {
	return &OnlineUserApi{srv}
}

// @title 在线用户列表
// @Summary 在线用户列表接口
// @Description 按条件查询在线用户列表
// @Tags Online
// @Accept json
// @Produce json
// @Param username query string false "username"
// @Param addr query string false "所在地区"
// //////////@Success 200 {object} resp.Response{data=resp.PageResult[entity.OnlineUserDto]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.OnlineUserDto]]
// @Router /api/online [get]
// @Security Bearer
func (p *OnlineUserApi) Index(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	addr := c.DefaultQuery("addr", "")
	page.StartPage(c)
	result, err := p.OnlineUserSrv.List(c.Request.Context(), username, addr)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 批量用户下线
// @Summary 批量用户下线接口
// @Description 批量强制用户下线
// @Tags Online
// @Accept json
// @Produce json
// @Param ids path string false "多用户id"
// @Success 200 {object} resp.Response[any]
// @Router /api/online/{ids} [delete]
// @Security Bearer
func (p *OnlineUserApi) KickOut(c *gin.Context) {
	ids := admin.ParseParamArray[string](c, "ids")
	if err := p.OnlineUserSrv.BatchKickOut(c.Request, ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}
