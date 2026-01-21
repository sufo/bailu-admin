/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package system

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/service/sys"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/log"
	"strconv"
)

var DeptSet = wire.NewSet(wire.Struct(new(DeptApi), "*"))

type DeptApi struct {
	DeptSrv sys.DeptService
}

// @title 部门列表接口
// @Summary 部门列表接口
// @Description 可按部门名称和状态查询部门列表接口
// @Tags Department
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param name query string false "部门名称"
// @Param status query string false "状态（1:启动，2:禁用）"
// //@Success 200 {object} resp.Response[entity.Dept]
// @response default {object} resp.Response[[]entity.Dept]
// @Router /api/dept [get]
// @Security Bearer
func (d *DeptApi) Index(c *gin.Context) {
	name := c.Query("name")
	status := c.Query("status")
	depts, err := d.DeptSrv.Tree(c, name, status)
	if err != nil {
		resp.InternalServerError(c)
	}
	resp.OKWithData(c, depts)
}

// @title 创建部门
// @Summary 创建部门接口
// @Description 创建部门接口
// @Tags Department
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.Dept true "部门信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/dept [post]
func (d *DeptApi) Create(c *gin.Context) {
	var dept entity.Dept
	if err := c.ShouldBindJSON(&dept); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		err := d.DeptSrv.Create(c.Request.Context(), &dept)
		if err != nil {
			resp.InternalServerError(c)
			return
		}
		resp.Ok(c)
	}
}

// @title 编辑部门
// @Summary 编辑部门接口
// @Description 编辑部门接口
// @Tags Department
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param body body entity.Dept true "部门信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/dept [put]
func (d *DeptApi) Edit(c *gin.Context) {
	var item entity.Dept
	if err := c.ShouldBindJSON(&item); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		if item.ID == item.Pid {
			resp.FailWithMsg(c, "上级不能为自己")
			return
		}
		err := d.DeptSrv.Update(c.Request.Context(), &item)
		if err != nil {
			resp.InternalServerError(c)
			return
		}
		resp.Ok(c)
	}
}

// @title 删除部门
// @Summary 删除部门接口
// @Description 删除部门接口
// @Tags Department
// @Accept json
// @Produce json
// @Security Bearer
// @Param deptId path string true "部门id"
// @Success 200 {object} resp.Response[any]
// @Router /api/dept/{deptId} [delete]
func (d *DeptApi) Destroy(c *gin.Context) {
	var ctx = c.Request.Context()
	var val = c.Param("deptId")
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(respErr.BadRequestError)
	}

	exist, err := d.DeptSrv.HasChildByID(ctx, id)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	//存在子菜单则不允许删除
	if exist {
		resp.FailWithMsg(c, i18n.DefTr("admin.hasChildMenus"))
		return
	}
	err = d.DeptSrv.Delete(ctx, id)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// @title 获取部门树
// @Summary 获取部门树接口
// @Description 获取部门树接口
// @Tags Department
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// //@Success 200 {object} resp.Response{data=array{key=string,label=string,children=array{key=string,label=string}}}\esponse{data=array{key=string,label=string,children=array{key=string,label=string}}}
// @Success 200 {object} resp.Response[any]{data=array{key=string,label=string,children=array{key=string,label=string}}}
// @Router /api/dept/tree [get]
func (d *DeptApi) Tree(c *gin.Context) {
	tree, err := d.DeptSrv.Tree(c.Request.Context(), "", "1")
	if err != nil {
		resp.InternalServerError(c)
	}
	data, err := admin.TransformTree[*entity.Dept](tree, []string{"ID", "Name", "Children"}, []string{"key", "label", "children"})
	if err != nil {
		log.L.Error(err)
		resp.InternalServerError(c)
		return
	}
	resp.OKWithData(c, data)
}

// 查询部门树 和 根据角色id查询关联的部门
// TreeSelect 获取部门树和所在部门集合
// @Summary 根据角色id查询关联的部门集合
// @Description 根据角色id查询关联的部门集合
// @Tags Department
// @Accept json
// @Produce json
// @Security Bearer
// @Param roleId path string true "角色id"
// //@Success 200 {object} resp.Response{data=object{checkedIds=[]integer,list=array{key=string,label=string,children=array{key=string,label=string}}}}
// @Success 200 {object} resp.Response[any]{data=object{checkedIds=[]integer,list=array{key=string,label=string,children=array{key=string,label=string}}}}
// @Router /api/dept/tree/{roleId} [get]
func (d *DeptApi) TreeSelect(c *gin.Context) {
	ctx := c.Request.Context()
	roleId := admin.ParseParamId(c, "roleId")
	checkIds, err := d.DeptSrv.DeptIdsByRoleId(ctx, roleId)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
		return
	}
	depts, err := d.DeptSrv.Tree(ctx, "", "")
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
		return
	}
	var result = map[string]any{
		"checkedIds": checkIds,
		"list":       transformDeptTree(depts),
	}
	resp.OKWithData(c, result)
}
