/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"bailu/app/api/admin"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/log"
	"bailu/utils/page"
)

var DictSet = wire.NewSet(wire.Struct(new(DictApi), "*"))

type DictApi struct {
	DictSrv *sys.DictService
}

// @title Index 字典列表
// @Summary 字典列表接口
// @Description 可按字典名称查询字典列表接口
// @Tags Dict
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param name query string false "字典名称"
// //////// @Security ApiKeyAuth
// //// @response default {object} resp.Response{data=resp.PageResult[entity.Dict]}
// @response default {object} resp.Response[resp.PageResult[entity.Dict]]
// @Success 200 {object} resp.Response[any]
// @Router /api/dict [get]
// @Security Bearer
func (d *DictApi) Index(c *gin.Context) {
	search := c.DefaultQuery("name", "")
	page.StartPage(c)
	result, err := d.DictSrv.List(c.Request.Context(), search)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title Create 创建字典接口
// @Summary 创建字典接口
// @Description 创建字典接口
// @Tags Dict
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body entity.Dict true "字典信息"
// //////// @Security ApiKeyAuth
// @Success 200 {object} resp.Response[any]
// @Router /api/dict [post]
// @Security Bearer
func (d *DictApi) Create(c *gin.Context) {
	var dict entity.Dict
	if err := c.ShouldBindJSON(&dict); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//检查字段唯一性
	if !d.DictSrv.CheckUnique(ctx, "code=?", dict.Code) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dict.Code))
		return
	}
	//检查字段唯一性
	if !d.DictSrv.CheckUnique(ctx, "name=?", dict.Name) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dict.Name))
		return
	}
	err := d.DictSrv.CreateDict(ctx, &dict)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.Ok(c)
}

// @title Edit 修改字典接口
// @Summary 修改字典接口
// @Description 修改字典接口
// @Tags Dict
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body entity.Dict true "字典信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/dict [put]
// @Security Bearer
func (d *DictApi) Edit(c *gin.Context) {
	var dict entity.Dict
	if err := c.ShouldBindJSON(&dict); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//检查字段唯一性
	if !d.DictSrv.CheckUnique(ctx, "code=? and id != ?", dict.Code, dict.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dict.Code))
		return
	}
	//检查字段唯一性
	if !d.DictSrv.CheckUnique(ctx, "name=? and id != ?", dict.Name, dict.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dict.Name))
		return
	}
	_dict, err := d.DictSrv.UpdateDict(ctx, &dict)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, _dict)
}

// @title Destroy 批量删除字典接口
// @Summary 批量删除字典接口
// @Description 批量删除字典接口
// @Tags Dict
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param codes path []string true "字典编码集合"
// @Success 200 {object} resp.Response[any]
// @Router /api/dict/{codes} [delete]
// @Security Bearer
func (d *DictApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamArray[string](c, "codes")
	if err := d.DictSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}
