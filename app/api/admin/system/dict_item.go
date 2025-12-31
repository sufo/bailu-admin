/**
 * Create by sufo
 * @Email ouamour@gmail.com
 * @Desc
 */

package system

import (
	"bailu/app/api/admin"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	"bailu/global/consts"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/log"
	"bailu/utils/page"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var DictItemSet = wire.NewSet(wire.Struct(new(DictItemApi), "*"))

type DictItemApi struct {
	DictItemSrv *sys.DictItemService
}

// @title Index 字典项列表
// @Summary 字典项列表接口
// @Description 可按字典名称查询字典项列表接口
// @Tags DictItem
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param code query string false "字典编码"
// @Param label query string false "字典项标签"
// @Param status query string false "字典状态"
// @response default {object} resp.Response[resp.PageResult[entity.DictItem]]
// @Success 200 {object} resp.Response[any]
// @Router /api/dictItem/{code} [get]
// @Security Bearer
func (d *DictItemApi) Index(c *gin.Context) {
	dictCode := c.Param("code")
	label := c.DefaultQuery("label", "")
	status := c.DefaultQuery("status", "")
	if dictCode == "" {
		resp.BadRequest(c)
		return
	}
	page.StartPage(c)
	result, err := d.DictItemSrv.List(c.Request.Context(), dictCode, label, status)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title Create 创建字典项
// @Summary 创建字典项接口
// @Description 创建字典项接口
// @Tags DictItem
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body entity.DictItem true "字典信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/dictItem [post]
// @Security Bearer
func (d *DictItemApi) Create(c *gin.Context) {
	var dictItem entity.DictItem
	if err := c.ShouldBindJSON(&dictItem); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//检查字段唯一性
	if !d.DictItemSrv.CheckUnique(ctx, "code=? and label=?", dictItem.Code, dictItem.Label) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dictItem.Label))
		return
	}
	//检查字段唯一性
	if !d.DictItemSrv.CheckUnique(ctx, "code=? and value=?", dictItem.Code, dictItem.Value) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dictItem.Value))
		return
	}
	err := d.DictItemSrv.Create(ctx, &dictItem)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.Ok(c)
}

func (d *DictItemApi) canEdit(c *gin.Context, ids []uint64) bool {
	//检查是否可以更改
	user, _ := c.Get(consts.REQUEST_USER)
	userDto := user.(*entity.OnlineUserDto)
	if !userDto.IsSuper() {
		oldItem, err := d.DictItemSrv.DictItemRepo.FindByIds(c.Request.Context(), ids)
		if err != nil {
			panic(respErr.InternalServerErrorWithError(err))
		}
		for _, item := range oldItem {
			if *item.Fixed {
				return false
			}
		}
	}
	return true
}

// @title Edit 修改字典项
// @Summary 修改字典项接口
// @Description 修改字典项接口
// @Tags DictItem
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body entity.DictItem true "字典信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/dictItem [put]
// @Security Bearer
func (d *DictItemApi) Edit(c *gin.Context) {
	var dItem entity.DictItem
	if err := c.ShouldBindJSON(&dItem); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//检查字段唯一性
	if !d.DictItemSrv.CheckUnique(ctx, "code=? and label=? and id != ?", dItem.Code, dItem.Label, dItem.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dItem.Label))
		return
	}
	//检查字段唯一性
	if !d.DictItemSrv.CheckUnique(ctx, "code=? and value=? and id != ?", dItem.Code, dItem.Value, dItem.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", dItem.Value))
		return
	}
	//检查是否可以更改
	can := d.canEdit(c, []uint64{dItem.ID})
	if !can {
		panic(i18n.DefTr("tip.noSupportAction"))
	}
	//可以更改
	_dict, err := d.DictItemSrv.Update(ctx, &dItem)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, _dict)
}

// @title Destroy 批量删除字典项
// @Summary 批量删除字典项接口
// @Description 批量删除字典项接口
// @Tags DictItem
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param ids path []integer true "字典编码集合"
// @Success 200 {object} resp.Response[any]
// @Router /api/dictItem/{ids} [delete]
// @Security Bearer
func (d *DictItemApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	//检查是否可以更改
	can := d.canEdit(c, ids)
	if !can {
		panic(i18n.DefTr("tip.noSupportAction"))
	}

	if err := d.DictItemSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title ChangeStatus 修改字典项状态
// @Summary 修改字典项状态接口
// @Description 修改字典项状态接口
// @Tags DictItem
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body dto.StatusParam true "id和状态"
// @Success 200 {object} resp.Response[any]
// @Router /api/dictItem [patch]
// @Security Bearer
func (d *DictItemApi) Status(c *gin.Context) {
	var param dto.StatusParam
	if err := c.ShouldBindJSON(&param); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	err := d.DictItemSrv.ChangeStatus(ctx, param)
	if err != nil {
		resp.InternalServerError(c)
		return
	}
	resp.Ok(c)
}

// @title 字典项下拉选项列表
// @Summary 字典项下拉选项列表接口
// @Description 字典项下拉选项列表接口
// @Tags DictItem
// @Accept json
// @Produce json
// @Security Bearer
// @Param code query string true "字典编码"
// // @Success 200 {object} resp.Response[vo.Option[string]]
// @Success 200 {object} resp.Response[any]{data=array{value=string,label=string,isDefault=bool}}
// @Router /api/dictItem/options [get]
func (d *DictItemApi) Options(c *gin.Context) {
	dictCode, bool := c.GetQuery("code")
	if !bool {
		resp.BadRequest(c)
		return
	}
	result, err := d.DictItemSrv.FindOptions(c.Request.Context(), dictCode)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}
