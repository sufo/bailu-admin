/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package repo

import (
	"context"
	"gorm.io/gorm"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
)

//var RoleSet = wire.NewSet(wire.Struct(new(RoleRepo), "*"))

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	r := base.Repository[entity.Role]{db}
	return &RoleRepo{r}
}

type RoleRepo struct {
	base.Repository[entity.Role]
}

func (r *RoleRepo) findRoles(ctx context.Context, name string, status string) {
	var params map[string]interface{}
	if name != "" {
		params["name"] = name
	}
	if status != "" {
		params["status"] = status
	}
	//r.RoleRepo.WithWhere(ctx, params).find
}

//	关于Gorm执行原生SQL
//
// **********语句字段要小写************
// ***********查询用db.Raw,其他用db.Exec
// *********** 字段大小写要对应上 **************
// *************** 注意要 defer rows.Close()
func (r *RoleRepo) UntiedMenu(menuId uint64) error {
	return r.DB.Exec("delete from sys_role_menu where menu_id = ?", menuId).Error
}
