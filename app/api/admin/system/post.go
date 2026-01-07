/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 岗位
 */

package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/service/sys"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils/page"
)

var PostSet = wire.NewSet(wire.Struct(new(PostApi), "*"))

type PostApi struct {
	PostSrv *sys.PostService
}

// @title 岗位列表
// @Summary 岗位列表接口
// @Description 可按条件查询岗位列表接口
// @Tags Post
// @Accept json
// @Produce json
// @Security Bearer
// @Param query query dto.PostParams false "查询参数"
// //@Success 200 {object} resp.Response{data=resp.PageResult[entity.Post]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.Post]]
// @Router /api/post [get]
func (p *PostApi) Index(c *gin.Context) {
	var queryParams dto.PostParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		panic(respErr.BadRequestError)
	}
	page.StartPage(c)
	result, err := p.PostSrv.List(c.Request.Context(), queryParams)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 创建岗位
// @Summary 创建岗位接口
// @Description 创建岗位接口
// @Tags Post
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body entity.Post true "岗位信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/post [post]
func (p *PostApi) Create(c *gin.Context) {
	var post entity.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//检查字段唯一性
	if !p.PostSrv.CheckUnique(ctx, "name=?", post.Name) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", post.Name))
		return
	}
	if !p.PostSrv.CheckUnique(ctx, "post_code=?", post.PostCode) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", post.PostCode))
		return
	}

	//创建
	if err := p.PostSrv.Create(ctx, &post); err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// @title 修改岗位
// @Summary 修改岗位接口
// @Description 修改岗位接口
// @Tags Post
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body entity.Post true "岗位信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/post [put]
func (p *PostApi) Edit(c *gin.Context) {
	var post entity.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//检查字段唯一性
	if !p.PostSrv.CheckUnique(ctx, "id!=? and name=?", post.ID, post.Name) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", post.Name))
		return
	}
	if !p.PostSrv.CheckUnique(ctx, "id!= ? and post_code=?", post.ID, post.PostCode) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", post.PostCode))
		return
	}
	if result, err := p.PostSrv.Update(ctx, &post); err != nil {
		resp.FailWithError(c, respErr.InternalServerError)
		return
	} else {
		resp.OKWithData(c, result)
	}
}

// @title 删除岗位
// @Summary 批量删除岗位接口
// @Description 批量删除岗位接口
// @Tags Post
// @Accept json
// @Produce json
// @Security Bearer
// @Param ids path string true "岗位id集合"
// @Success 200 {object} resp.Response[any]
// @Router /api/post/{ids} [delete]
func (p *PostApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	if err := p.PostSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 岗位下拉选项列表
// @Summary 岗位下拉选项列表接口
// @Description 岗位下拉选项列表接口
// @Tags Post
// @Accept json
// @Produce json
// @Security Bearer
// //@Success 200 {object} resp.Response[vo.Option[integer]]
// @Success 200 {object} resp.Response[any]{data=array{value=integer,label=string,isDefault=bool}}
// @Router /api/post/options [get]
func (p *PostApi) Options(c *gin.Context) {
	result, err := p.PostSrv.FindOptions(c.Request.Context())
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}
