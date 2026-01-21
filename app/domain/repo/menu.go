/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package repo

import (
	"context"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/app/domain/repo/util"
	"gorm.io/gorm"
)

// var MenuSet = wire.NewSet(wire.Struct(new(MenuRepo), "*"))
func NewMenuRepo(db *gorm.DB) *MenuRepo {
	r := base.Repository[entity.Menu]{db}
	return &MenuRepo{r}
}

type MenuRepo struct {
	base.Repository[entity.Menu]
}

// 根据角色集合查询对应的菜单
func (m *MenuRepo) FindByRolesAndTypeNot(ctx context.Context, roleIds []uint64, menuType string, pid *uint64) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	db := util.GetDBWithModel[entity.Menu](ctx, m.DB)
	db = db.Joins("JOIN sys_role_menu rm ON rm.menu_id = sys_menu.ID AND rm.role_id IN (?)", roleIds)
	if menuType != "" {
		db = db.Where("type <> ?", menuType)
	}
	if pid != nil {
		db = db.Where("pid = ?", pid)
	}
	db = db.Where("status=1")
	err := db.Order("sort ASC, created_at ASC").Find(&menus).Error
	return menus, err
}

//func (m *MenuRepo) FindMenus(ctx context.Context, roleIds []uint64, params dto.MenuParams) ([]*entity.Menu, error) {
//	var menus []*entity.Menu
//	db := util.GetDBWithModel[entity.Menu](ctx, m.DB)
//	//AND条件会对sys_role_menu起过滤作用，但是不会影响menu的查询结果，匹配不上的补NULL
//	//db = db.Joins("JOIN sys_role_menu rm ON rm.menu_id = sys_menu.ID AND rm.role_id IN (?)", roleIds)
//
//	//nil表示是超管
//	if roleIds != nil {
//		db = db.Joins("JOIN sys_role_menu rm ON rm.menu_id = sys_menu.id").Where("rm.role_id in ?", roleIds)
//	}
//
//	if params.Name != "" {
//		db = db.Where("sys_menu.name LIKE ?", "%"+params.Name+"%")
//	}
//	if params.ExcludeType != "" {
//		db = db.Where("sys_menu.type <> ?", params.ExcludeType)
//	}
//	if params.Status != nil {
//		db = db.Where("sys_menu.status = ?", params.Status)
//	}
//	if params.Pid != nil {
//		db = db.Where("sys_menu.pid = ?", params.Pid)
//	}
//	err := db.Order("sys_menu.sort desc").Order("sys_menu.id").Find(&menus).Error
//	return menus, err
//}

func (m *MenuRepo) FindMenus(ctx context.Context, params dto.MenuParams) ([]*entity.Menu, error) {
	builder := base.NewQueryBuilder()
	builder.WithWhereStruct(params).WithPreload("Apis").
		WithOrder("sys_menu.pid", "sys_menu.sort asc").
		WithOrder("created_at")

	menus, err := m.FindByBuilder(ctx, builder)
	return menus.([]*entity.Menu), err
}

func (m *MenuRepo) FindMenusByUserId(ctx context.Context, userId uint64, params dto.MenuParams) ([]*entity.Menu, error) {
	builder := base.NewQueryBuilder()
	builder.WithTable("sys_menu m").WithPreload("Apis").
		WithJoin("left join sys_role_menu rm on rm.menu_id=m.id").
		WithJoin("left join sys_user_role ur on rm.role_id=ur.role_id").
		WithJoin("left join sys_role r on r.id=ur.role_id").
		WithWhere("ur.user_id=?", userId).
		WithWhereStructAndAlias(params, "m").
		WithOrder("m.pid", "m.sort asc")
	menus, err := m.FindByBuilder(ctx, builder)
	return menus.([]*entity.Menu), err
}

// 通过roles和 MenuParams查询菜单
func (m *MenuRepo) FindByRolesAndParam(ctx context.Context, roleIds []uint64, params dto.MenuParams) ([]*entity.Menu, error) {
	builder := base.NewQueryBuilder()
	builder.WithTable("sys_menu m").WithPreload("Apis").
		WithJoin("JOIN sys_role_menu rm ON rm.menu_id = m.id").
		WithWhere("rm.role_id IN ?", roleIds).
		WithWhereStructAndAlias(params, "m"). // Apply additional filters from dto.MenuParams
		WithOrder("m.pid", "m.sort asc")      // Consistent ordering

	menus, err := m.FindByBuilder(ctx, builder)
	if err != nil {
		return nil, err
	}
	return menus.([]*entity.Menu), err
}

func (m *MenuRepo) FindByRoleId(ctx context.Context, roleId uint64) ([]uint64, error) {
	var role entity.Role
	if err := m.DB.WithContext(ctx).First(&role, roleId).Error; err != nil {
		return nil, err
	}

	var menus []uint64
	db := util.GetDBWithModel[entity.Menu](ctx, m.DB)
	db = db.Select("id").Joins("LEFT JOIN sys_role_menu rm ON rm.menu_id = sys_menu.ID")
	db = db.Where("rm.role_id=?", roleId)

	// if menuCheckStrictly is false, it means cascade selection, so exclude parent nodes
	if !role.MenuCheckStrictly {
		db = db.Where("sys_menu.id not in (select pid from sys_menu m inner join sys_role_menu rm ON rm.menu_id = sys_menu.ID and rm.role_id=?)", roleId)
	}

	err := db.Order("sys_menu.pid ASC, sys_menu.sort asc").Find(&menus).Error
	return menus, err
}

func (m *MenuRepo) FindByRoleIds(ctx context.Context, roleIds []uint64) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	db := util.GetDBWithModel[entity.Menu](ctx, m.DB)
	err := db.Joins("JOIN sys_role_menu rm ON rm.menu_id = sys_menu.id").
		Where("rm.role_id in ?", roleIds).
		Order("sys_menu.pid ASC, sys_menu.sort asc").Find(&menus).Error
	return menus, err
}
