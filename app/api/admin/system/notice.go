/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Update 2024/6/21
 * @Desc 通知消息(通知公告)
 */

package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"bailu/app/api/admin"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	respErr "bailu/pkg/exception"
	"bailu/pkg/log"
	"bailu/utils/page"
)

var NoticeSet = wire.NewSet(wire.Struct(new(NoticeApi), "*"))

type NoticeApi struct {
	NoticeSrv *sys.NoticeService
}

// @title 通知公告列表
// @Summary 通知公告列表接口
// @Description 通知公告列表接口
// @Tags Notice
// @Accept json
// @Produce json
// @Security Bearer
// @Param query query dto.NoticeParams false "查询参数"
// //@Success 200 {object} resp.Response{data=resp.PageResult[entity.Notice]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.Notice]]
// @Router /api/notice [get]
func (n *NoticeApi) Index(c *gin.Context) {
	var queryParams dto.NoticeParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		panic(respErr.BadRequestError)
	}
	page.StartPage(c)
	result, err := n.NoticeSrv.List(c.Request.Context(), queryParams)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// Create 创建通知公告
// @Summary 创建通知公告接口
// @Description 创建通知公告接口
// @Tags Notice
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.Notice true "通知信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/notice [post]
func (a *NoticeApi) Create(c *gin.Context) {
	var anc entity.Notice
	if err := c.ShouldBindJSON(&anc); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//创建
	if err := a.NoticeSrv.Create(ctx, &anc); err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// @title 修改通知公告
// @Summary 修改通知公告接口
// @Description 修改通知公告接口
// @Tags Notice
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body entity.Dept true "通知信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/notice [put]
func (a *NoticeApi) Edit(c *gin.Context) {
	var anc entity.Notice
	if err := c.ShouldBindJSON(&anc); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	if result, err := a.NoticeSrv.Update(ctx, &anc); err != nil {
		resp.FailWithError(c, respErr.InternalServerError)
		return
	} else {
		resp.OKWithData(c, result)
	}
}

// @title 删除通知
// @Summary 删除通知接口
// @Description 删除通知接口
// @Tags Notice
// @Accept json
// @Produce json
// @Security Bearer
// @Param ids path string true "通知id数组"
// @Success 200 {object} resp.Response[any]
// @Router /api/notice/{ids} [delete]
func (a *NoticeApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	if err := a.NoticeSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 发布通知
// @Summary 发布通知接口
// @Description 发布通知接口
// @Tags Notice
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "通知id"
// @Success 200 {object} resp.Response[any]
// @Router /api/notice/release/{id} [patch]
func (a *NoticeApi) Release(c *gin.Context) {
	id := admin.ParseParamId(c, "id")
	ctx := c.Request.Context()
	//创建
	if err := a.NoticeSrv.ReleaseNotice(ctx, id); err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// @title 撤销通知
// @Summary 撤销通知接口
// @Description 撤销通知
// @Tags Notice
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "通知id"
// @Success 200 {object} resp.Response[any]
// @Router /api/notice/revoke/{id} [patch]
func (a *NoticeApi) Revoke(c *gin.Context) {
	id := admin.ParseParamId(c, "id")
	ctx := c.Request.Context()
	//创建
	if err := a.NoticeSrv.CancelNotice(ctx, id); err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}
