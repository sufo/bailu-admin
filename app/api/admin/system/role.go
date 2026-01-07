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
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/service/sys"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils/page"
	"strconv"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleApi), "*"))

type RoleApi struct {
	RoleSrv   *sys.RoleService
	DeptSrv   *sys.DeptService
	UserSrv   *sys.UserService
	OnlineSrv *sys.OnlineService
}

// @title 角色列表
// @Summary 角色列表接口
// @Description 按条件查询角色列表接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param query query dto.RoleParams false "查询条件"
// ////////@Success 200 {object} resp.Response{data=resp.PageResult[entity.Role]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.Role]]
// @Router /api/role [get]
func (r *RoleApi) Index(c *gin.Context) {
	//开启分页
	page.StartPage(c)

	var queryParams dto.RoleParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		panic(respErr.BadRequestError)
	}

	roles, err := r.RoleSrv.List(c.Request.Context(), queryParams)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, roles)
}

// @title 创建角色
// @Summary 创建角色接口
// @Description 创建角色接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.Role true "角色信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/role [post]
func (r *RoleApi) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var role dto.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		panic(respErr.BadRequestError)
	}

	//检查角色名
	if !r.RoleSrv.CheckUnique(ctx, "name=?", role.Name) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", role.Name))
		return
	}
	//检查roleKey
	if !r.RoleSrv.CheckUnique(ctx, "role_key=?", role.RoleKey) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", role.RoleKey))
		return
	}

	err := r.RoleSrv.Create(ctx, &role)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.Ok(c)
}

// @title 修改角色
// @Summary 修改角色接口
// @Description 修改角色接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.Role true "角色信息"
// @Success 200 {object} resp.Response[any]
// @Router /api/role [put]
func (r *RoleApi) Edit(c *gin.Context) {
	ctx := c.Request.Context()

	var role dto.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		panic(respErr.BadRequestError)
	}

	//检查角色名
	if !r.RoleSrv.CheckUnique(ctx, "name=? and id !=?", role.Name, role.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", role.Name))
		return
	}
	//检查roleKey
	if !r.RoleSrv.CheckUnique(ctx, "role_key=? and id !=?", role.RoleKey, role.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", role.RoleKey))
		return
	}
	if err := r.RoleSrv.Update(ctx, &role); err != nil {
		resp.FailWithError(c, err)
		return
	}
	//更新缓存信息
	if !role.IsSuper() {
		reqUser, _ := c.Get(consts.REQUEST_USER)
		var onlineUser = reqUser.(*entity.OnlineUserDto)
		user, err := r.UserSrv.FindById(ctx, onlineUser.ID)
		if err != nil {
			log.L.Errorf("%d 更新缓存失败！", user.ID)
		}
		r.OnlineSrv.Save(user, c.Request, onlineUser.Token)
	}
	resp.Ok(c)
}

// @title 修改角色状态
// @Summary 修改角色状态接口
// @Description 修改角色状态接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body object{id=integer,status=integer} true "id和status"
// @Success 200 {object} resp.Response[any]
// @Router /api/role/status [patch]
func (r *RoleApi) Status(c *gin.Context) {
	ctx := c.Request.Context()
	//参数处理
	var params = make(map[string]string)
	if err := c.ShouldBindJSON(&params); err != nil {
		panic(respErr.BadRequestError)
	}
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	status, err := strconv.ParseUint(params["status"], 10, 8)
	if err != nil {
		panic(respErr.BadRequestError)
	}

	r.RoleSrv.CheckRoleAllowed(&dto.Role{ID: id})
	var ids = []uint64{id}
	r.RoleSrv.CheckDataScope(ctx, ids)

	if err := r.RoleSrv.UpdateStatus(ctx, id, uint8(status)); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 批量删除角色
// @Summary 批量删除角色接口
// @Description 批量删除角色接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param deptId path string true "角色id集合"
// @Success 200 {object} resp.Response[any]
// @Router /api/role/{roleIds} [delete]
func (r *RoleApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "roleIds")
	if err := r.RoleSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// 根据角色id查询角色信息
// @title id查询角色
// @Summary 根据id查询角色信息接口
// @Description 根据id查询角色信息接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param roleId path integer false "roleId"
// //////@Success 200 {object} resp.Response{data=entity.Role}
// @Success 200 {object} resp.Response[entity.Role]
// @Router /api/role/{id} [get]
func (r *RoleApi) Detail(c *gin.Context) {
	id := admin.ParseParamId(c, "roleId")
	if result, err := r.RoleSrv.Get(c.Request.Context(), id); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.OKWithData(c, result)
	}
}

// 修改保存用户数据权限
// @title 修改用户数据权限
// @Summary 修改用户数据权限接口
// @Description 修改用户数据权限接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body dto.RoleDataScope true "dataScope"
// @Success 200 {object} resp.Response[any]
// @Router /api/dataScope [patch]
func (r *RoleApi) DataScope(c *gin.Context) {
	var roleDataScope dto.RoleDataScope
	if err := c.ShouldBindJSON(&roleDataScope); err != nil {
		panic(respErr.BadRequestError)
	}
	if err := r.RoleSrv.AuthDataScope(c.Request.Context(), &roleDataScope); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// 批量取消用户授权
// @title 批量取消用户授权
// @Summary 批量取消用户授权接口
// @Description 批量取消用户授权接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body dto.RoleUsers true "请求参数"
// @Success 200 {object} resp.Response[any]
// @Router /api/role/authUser [delete]
func (r *RoleApi) CancelAuthUsers(c *gin.Context) {
	var params dto.RoleUsers
	if err := c.ShouldBindJSON(&params); err != nil {
		panic(respErr.BadRequestError)
	}
	err := r.RoleSrv.DeleteUserRole(c.Request.Context(), &params)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
		return
	}
}

// 批量用户授权
// @title 批量用户授权
// @Summary  批量用户授权接口
// @Description 批量用户授权接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body dto.RoleUsers true "请求参数"
// @Success 200 {object} resp.Response[any]
// @Router /api/role/authUser [post]
func (r *RoleApi) AuthUsers(c *gin.Context) {
	var params dto.RoleUsers
	if err := c.ShouldBindJSON(&params); err != nil {
		panic(respErr.BadRequestError)
	}

	err := r.RoleSrv.InsertAuthUser(c.Request.Context(), params)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
		return
	}
}

// @title 角色下拉选项列表
// @Summary 角色下拉选项列表接口
// @Description 获取角色下拉选项列表接口
// @Tags Role
// @Accept json
// @Produce json
// @Security Bearer
// //////@Success 200 {object} resp.Response{data=vo.Option[uint64]}
// //@Success 200 {object} resp.Response[vo.Option[uint64]]
// @Success 200 {object} resp.Response[any]{data=array{value=integer,label=string,isDefault=bool}}
// @Router /api/role/options [get]
func (r *RoleApi) Options(c *gin.Context) {
	result, err := r.RoleSrv.FindOptions(c.Request.Context())
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// 转换结果
func transformDeptTree(tree []*entity.Dept) []map[string]any {
	var depts = make([]map[string]any, 0)
	for _, item := range tree {
		var temp = make(map[string]any, 0)
		temp["key"] = item.ID
		temp["label"] = item.Name
		var children = item.Children
		if children != nil && len(children) > 0 {
			temp["children"] = transformDeptTree(children)
		}
		depts = append(depts, temp)
	}
	return depts
}
