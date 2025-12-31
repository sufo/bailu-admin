/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package system

import (
	"bailu/app/api/admin"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"strconv"
)

var MenuSet = wire.NewSet(wire.Struct(new(MenuApi), "*"))

type MenuApi struct {
	MenuSrv *sys.MenuService
}

// @title 菜单列表
// @Summary 菜单列表接口
// @Description 查询菜单列表
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param query query dto.MenuParams false "查询参数"
// ////@response default {object} resp.Response{data=resp.PageResult[entity.Menu]}
// @response default {object} resp.Response[resp.PageResult[entity.Menu]]
// @Success 200 {object} resp.Response[any]
// @Router /api/menu [get]
// @Security Bearer
func (m *MenuApi) Index(c *gin.Context) {
	//name := c.Query("name")
	//enableStr := c.Query("enable")
	var menuParams dto.MenuParams
	if err := c.ShouldBindQuery(&menuParams); err != nil {
		panic(respErr.BadRequestError)
	}

	if menus, err := m.MenuSrv.FindMenuList(c, menuParams); err != nil {
		resp.InternalServerError(c)
	} else {
		//BuildTree存在问题，当menus如果不存在pid为0的顶级菜单，则会返回空
		//menuTree := m.MenuSrv.BuildTree(menus, 0)

		menuTree := m.MenuSrv.BuildAnyTree(menus)
		resp.OKWithData(c, menuTree)
	}
}

// get dynamic routes
// @title 动态路由
// @Summary 获取动态路由接口
// @Description 获取动态路由接口
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @response default {object} resp.Response[resp.PageResult[vo.Route]]
// @Success 200 {object} resp.Response[any]
// @Router /api/menu/routes [get]
func (m *MenuApi) Routes(c *gin.Context) {
	routes, err := m.MenuSrv.FindRoutes(c)
	if err != nil {
		panic(respErr.InternalServerError)
	}
	resp.OKWithData(c, routes)
}

// get dynamic routes
// @title 菜单数
// @Summary 获取子菜单树接口
// @Description 获取子菜单树接口
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param menuId path string true "menuId"
// //@Success 200 {object} resp.Response{data=resp.PageResult[entity.Menu]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.Menu]]
// @Router /api/menu/{menuId} [get]
func (m *MenuApi) SubMenus(c *gin.Context) {
	val := c.Param("menuId")
	pid, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	menus, err := m.MenuSrv.FindByPid(c, pid)
	if err != nil {
		log.L.Error(err)
		panic(respErr.InternalServerError)
	}
	resp.OKWithData(c, menus)
}

// 获取目录和菜单 排除按钮
// @title 菜单树(不包含按钮)
// @Summary 获取菜单树(不包含按钮)接口
// @Description 获取菜单树(不包含按钮)接口
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// //@Success 200 {object} resp.Response{data=resp.PageResult[entity.Menu]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.Menu]]
// @Router /api/menu/menus [get]
func (m *MenuApi) Menus(c *gin.Context) {
	var status = 1
	var params = dto.MenuParams{Status: &status, ExcludeType: "F"}
	if menus, err := m.MenuSrv.FindMenuList(c, params); err != nil {
		resp.InternalServerError(c)
	} else {
		menuTree := m.MenuSrv.BuildTree(menus, 0)
		resp.OKWithData(c, menuTree)
	}
}

// Create 创建菜单
// @title 创建菜单
// @Summary 创建菜单接口
// @Description 创建菜单接口
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param request body entity.Menu true "菜单信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/menu [post]
func (m *MenuApi) Create(c *gin.Context) {
	var item entity.Menu
	if err := c.ShouldBindJSON(&item); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		err := m.MenuSrv.Save(c, item)
		if err != nil {
			//resp.FailWithStatusAndMsg(c, status.StatusBadRequest, err.Error())
			resp.InternalServerError(c)
			return
		}
		resp.Ok(c)
	}
}

// @title 修改菜单
// @Summary 修改菜单接口
// @Description 修改菜单接口
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.Dept true "菜单信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/menu [put]
func (m *MenuApi) Edit(c *gin.Context) {
	var item entity.Menu
	if err := c.ShouldBindJSON(&item); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		if item.ID == *item.Pid {
			resp.FailWithMsg(c, "上级不能为自己")
			return
		}
		err := m.MenuSrv.Update(c.Request.Context(), item)
		if err != nil {
			resp.InternalServerError(c)
			return
		}
		resp.Ok(c)
	}
}

// @title 修改菜单
// @Summary 删除菜单接口
// @Description 删除菜单接口
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param menuId path string true "菜单id"
// @Success 200 {object} resp.Response[any]
// @Router /api/menu/{menuId} [delete]
func (m *MenuApi) Destroy(c *gin.Context) {
	var val = c.Param("menuId")
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(respErr.BadRequestError)
	}

	exist, err := m.MenuSrv.HasChildByID(c, id)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	//存在子菜单则不允许删除
	if exist {
		resp.FailWithMsg(c, i18n.DefTr("admin.hasChildMenus"))
		return
	}
	err = m.MenuSrv.Delete(c.Request.Context(), id)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// 角色新增时，下拉菜单树
// @title 获取菜单树
// @Summary 获取下拉菜单树
// @Description 角色新增时，下拉菜单树
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// //@Success 200 {object} resp.Response{data=array{key=string,label=string,children=array{key=string,label=string}}}
// @Success 200 {object} resp.Response[any]{data=array{key=string,label=string,children=array{key=string,label=string}}}
// @Router /api/menu/tree [get]
func (m *MenuApi) Tree(c *gin.Context) {
	var params = dto.MenuParams{}
	menus, err := m.MenuSrv.FindMenuList(c, params)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
		return
	}
	tree := m.MenuSrv.BuildTree(menus, 0)
	resp.OKWithData(c, transRoleMenuTree(tree))
}

// 获取角色对应的菜单树和已选择菜单
// @title 菜单树和已选择菜单
// @Summary 获取角色对应的菜单树和已选择菜单
// @Description 获取角色对应的菜单树和已选择菜单
// @Tags Menu
// @Accept json
// @Produce json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param roleId path string true "角色id"
// //@Success 200 {object} resp.Response{data=object{checkedIds=[]integer,menus=array{key=string,label=string,children=array{key=string,label=string}}}}
// @Success 200 {object} resp.Response[any]{data=object{checkedIds=[]integer,menus=array{key=string,label=string,children=array{key=string,label=string}}}}
// @Router /api/menu/tree/{roleId} [get]
func (m *MenuApi) TreeSelect(c *gin.Context) {
	var roleId = admin.ParseParamId(c, "roleId")
	mIds, err := m.MenuSrv.FindCheckedMenusByRoleId(c, roleId)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
		return
	}
	var params = dto.MenuParams{}
	menus, err := m.MenuSrv.FindMenuList(c, params)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
		return
	}

	//数据转换
	tree := m.MenuSrv.BuildTree(menus, 0)

	var data = map[string]any{
		"checkedIds": mIds,
		"menus":      transRoleMenuTree(tree),
	}
	resp.OKWithData(c, data)
}

// 转换结果
func transRoleMenuTree(menuTree []*entity.Menu) []map[string]any {
	var menus = make([]map[string]any, 0)
	for _, item := range menuTree {
		var temp = make(map[string]any, 2)
		temp["key"] = item.ID
		temp["label"] = item.Name
		var children = item.Children
		if children != nil && len(children) > 0 {
			temp["children"] = transRoleMenuTree(children)
		}
		menus = append(menus, temp)
	}
	return menus
}
