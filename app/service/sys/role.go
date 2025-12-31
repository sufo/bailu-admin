/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package sys

import (
	"context"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"bailu/app/core/appctx"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	"bailu/app/domain/repo/base"
	"bailu/app/domain/resp"
	"bailu/app/domain/vo"
	base2 "bailu/app/service/base"
	"bailu/global/consts"
	respErr "bailu/pkg/exception"
	"bailu/pkg/log"
)

// var RoleSet = wire.NewSet(wire.Struct(new(RoleService), "*"))
var RoleSet = wire.NewSet(wire.Struct(new(RoleOption), "*"), NewRoleService)

type RoleOption struct {
	RoleRepo  *repo.RoleRepo
	TransRepo *repo.Trans
}

type RoleService struct {
	base2.BaseService[entity.Role]
	RoleOption
}

func NewRoleService(opt RoleOption) *RoleService {
	return &RoleService{base2.BaseService[entity.Role]{opt.RoleRepo.Repository}, opt}
}

func (r *RoleService) UntiedMenu(menuId uint64) error {
	return r.RoleRepo.UntiedMenu(menuId)
}

func (r *RoleService) List(ctx context.Context, params dto.RoleParams) (*resp.PageResult[entity.Role], error) {

	builder := base.NewQueryBuilder()
	builder.WithJoin("left join sys_user_role as ur on ur.role_id=sys_role.id").
		WithJoin("left join sys_user as u on u.id=ur.user_id").
		WithJoin("left join sys_dept as d on d.id=u.dept_id").
		WithWhereStructAndAlias(params, "sys_role").
		WithDataScope(ctx, "d", "u").
		WithPagination(ctx).WithOrder("sort asc")
	if result, err := r.RoleRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		return result.(*resp.PageResult[entity.Role]), err
	}
}

func (r *RoleService) find(ctx context.Context, params dto.RoleParams) (any, error) {
	builder := base.NewQueryBuilder()
	builder.WithJoin("left join sys_user_role as ur on ur.role_id=sys_role.id").
		WithJoin("left join sys_user as u on u.id=ur.user_id").
		WithJoin("left join sys_dept as d on d.id=u.dept_id").
		WithWhereStructAndAlias(params, "sys_role").
		WithDataScope(ctx, "d", "u").
		WithPagination(ctx)
	return r.RoleRepo.FindByBuilder(ctx, builder)
}

func (r *RoleService) Get(ctx context.Context, id uint64) (*entity.Role, error) {
	return r.RoleRepo.FindById(ctx, id)
}

func (r *RoleService) Create(ctx context.Context, role *dto.Role) error {
	return r.TransRepo.Exec(ctx, func(ctx context.Context) error {
		var _role = &entity.Role{}
		err := copier.Copy(_role, role)
		if err != nil {
			return err
		}
		err = r.RoleRepo.Create(ctx, _role)
		if err != nil {
			return err
		}
		//创建后，_role会生成ID
		role.ID = _role.ID
		return r.addRoleMenu(ctx, role)
	})
}

func (r *RoleService) Update(ctx context.Context, role *dto.Role) error {
	r.CheckRoleAllowed(role)
	ids := []uint64{role.ID}
	r.CheckDataScope(ctx, ids)
	return r.TransRepo.Exec(ctx, func(ctx context.Context) error {
		//修改role
		var _role = &entity.Role{}
		if err := copier.Copy(_role, role); err != nil {
			return err
		}
		err := r.RoleRepo.UpdateColumns(ctx, _role.ID, _role).Error
		if err != nil {
			return err
		}
		//判断关联的菜单是否更改
		if role.IsMenusChanged {
			// 删除角色与菜单关联
			err = r.RoleRepo.GetDB(ctx).Unscoped().Where("role_id=?", role.ID).
				Delete(&entity.RoleMenu{}).Error
			if err != nil {
				return err
			}
			//插入角色
			return r.addRoleMenu(ctx, role)
		} else {
			return nil
		}
	})
}

func (r *RoleService) addRoleMenu(ctx context.Context, role *dto.Role) error {
	var roleMenus = make([]map[string]interface{}, 0)
	for _, mID := range role.Menus {
		rm := map[string]interface{}{
			"role_id": role.ID,
			"menu_id": mID,
		}
		roleMenus = append(roleMenus, rm)
	}
	if len(roleMenus) > 0 {
		return r.RoleRepo.Model(ctx, &entity.RoleMenu{}).
			Create(roleMenus).Error
	} else {
		return nil
	}
}

//func (r *RoleService) Delete(ctx context.Context, ids []uint64) error {
//	//check
//	for _, id := range ids {
//		r.CheckRoleAllowed(&dto.Role{ID: id})
//	}
//	r.CheckDataScope(ctx, ids)
//
//	return r.TransRepo.Exec(ctx, func(ctx context.Context) error {
//		//解除角色和用户关系
//		var userRole entity.UserRole
//		err := repo.GetDB(ctx, nil).Where("role_id IN ?", ids).Unscoped().Delete(&userRole).Error
//		if err != nil {
//			return err
//		}
//		//解除角色和权限菜单
//		var RoleMenu entity.RoleMenu
//		err = repo.GetDB(ctx, nil).Where("role_id in ?", ids).Unscoped().Delete(&RoleMenu).Error
//		if err != nil {
//			return err
//		}
//		return r.RoleRepo.Delete(ctx, ids)
//	})
//}

func (r *RoleService) Delete(ctx context.Context, ids []uint64) error {
	//check
	for _, id := range ids {
		r.CheckRoleAllowed(&dto.Role{ID: id})
	}
	r.CheckDataScope(ctx, ids)
	return r.TransRepo.Exec(ctx, func(ctx context.Context) error {
		//如果角色还绑定了用户，解除角色和用户关系
		var userRole entity.UserRole
		err := base.GetDB(ctx, nil).Where("role_id IN ?", ids).Unscoped().Delete(&userRole).Error
		if err != nil {
			return err
		}
		//解除角色和权限菜单
		var RoleMenu entity.RoleMenu
		err = base.GetDB(ctx, nil).Where("role_id in ?", ids).Unscoped().Delete(&RoleMenu).Error
		if err != nil {
			return err
		}
		//解除角色和部门
		var RoleDept entity.RoleDept
		err = base.GetDB(ctx, nil).Where("role_id in ?", ids).Unscoped().Delete(RoleDept).Error
		if err != nil {
			return err
		}
		return r.RoleRepo.Delete(ctx, ids)
	})
}

func (r *RoleService) AuthDataScope(ctx context.Context, roleDS *dto.RoleDataScope) error {
	if roleDS.ID == consts.SUPER_ROLE_ID {
		panic(respErr.ForbiddenError)
	}
	ids := []uint64{roleDS.ID}
	r.CheckDataScope(ctx, ids)

	return r.TransRepo.Exec(ctx, func(ctx context.Context) error {
		//修改角色
		if err := r.RoleRepo.Where(ctx, "id=?", roleDS.ID).Update("data_scope", roleDS.DataScope).Error; err != nil {
			return err
		}
		//删除角色与部门关联, 这里不要直接用r.RoleRepo.Where，因为这里是中间表
		if err := r.RoleRepo.GetDB(ctx).Where("role_id=?", roleDS.ID).Delete(&entity.RoleDept{}).Error; err != nil {
			return err
		}
		//新增角色和部门信息
		return r.insertRoleDept(ctx, roleDS)
	})
}

// 插入 sys_role_dept
func (r *RoleService) insertRoleDept(ctx context.Context, role *dto.RoleDataScope) error {
	if role.DeptIds == nil || len(role.DeptIds) == 0 {
		return nil
	}
	var roleDepts = make([]entity.RoleDept, 0)
	for _, deptId := range role.DeptIds {
		var roleDept entity.RoleDept
		roleDept.RoleId = role.ID
		roleDept.DeptId = deptId
		roleDepts = append(roleDepts, roleDept)
	}
	return r.RoleRepo.Model(ctx, &entity.RoleDept{}).Create(&roleDepts).Error
}

// 取消用户角色授权
func (r *RoleService) DeleteUserRole(ctx context.Context, params *dto.RoleUsers) error {
	builder := &base.QueryBuilder{}
	builder.WithTable("sys_user_role").
		WithWhereStruct(params)
	return r.RoleRepo.WithQueryBuilder(ctx, builder).Unscoped().Delete(&entity.UserRole{}).Error
}

// 批量选择用户授权
func (r *RoleService) InsertAuthUser(ctx context.Context, params dto.RoleUsers) error {
	r.CheckDataScope(ctx, []uint64{params.RoleId})
	var roleUserList = make(map[string]uint64, len(params.UserIds))
	for _, userId := range params.UserIds {
		roleUserList["role_id"] = params.RoleId
		roleUserList["user_id"] = userId
	}
	return r.RoleRepo.Create(ctx, params)
}

func (r *RoleService) UpdateStatus(ctx context.Context, id uint64, status uint8) error {
	return r.RoleRepo.Where(ctx, "id=?", id).UpdateColumn("status", status).Error
}

func (r *RoleService) CheckRoleAllowed(role *dto.Role) {
	if role.IsSuper() {
		panic(respErr.ForbiddenError)
	}
}

// 检查数据权限
func (r *RoleService) CheckDataScope(ctx context.Context, roleIds []uint64) {
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if !user.IsSuper() {
		var roles []*entity.Role
		builder := base.NewQueryBuilder()
		builder.WithJoin("left join sys_user_role as ur on ur.role_id=sys_role.id").
			WithJoin("left join sys_user as u on u.id=ur.user_id").
			WithJoin("left join sys_dept as d on d.id=u.dept_id").
			WithWhere("status=? and roleId in ?", 1, roleIds).
			WithDataScope(ctx, "d", "u")

		err := r.RoleRepo.WithQueryBuilder(ctx, builder).Find(&roles).Error
		if err != nil {
			panic(respErr.InternalServerErrorWithError(err))
		}
		if len(roles) == 0 {
			//panic(respErr.InternalServerErrorWithMsg("没有权限访问角色数据！"))
			log.L.Error("没有权限访问角色数据！")
			panic(respErr.ForbiddenError)
		}
	}
}

func (r *RoleService) FindOptions(ctx context.Context) ([]*vo.Option[uint64], error) {

	builder := base.NewQueryBuilder()
	builder.WithTable("sys_role r").
		WithSelect("r.id as value").WithSelect("r.name as label").
		WithJoin("left join sys_user_role as ur on ur.role_id=r.id").
		WithJoin("left join sys_user as u on u.id=ur.user_id").
		WithJoin("left join sys_dept as d on d.id=u.dept_id").
		WithWhere("r.status=?", 1). //启用
		WithDataScope(ctx, "d", "u").
		WithOrder("r.sort asc")
	options := make([]*vo.Option[uint64], 0)
	err := r.RoleRepo.FindModelByBuilder(ctx, builder, &options)
	return options, err
}
