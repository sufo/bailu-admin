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
	"strconv"
)

//var DeptSet = wire.NewSet(wire.Struct(new(DeptRepo), "*"))

func NewDeptRepo(db *gorm.DB) *DeptRepo {
	r := base.Repository[entity.Dept]{db}
	return &DeptRepo{r}
}

type DeptRepo struct {
	base.Repository[entity.Dept]
}

func (d *DeptRepo) FindDepts(ctx context.Context, name string, status string) (any, error) {
	build := base.NewQueryBuilder()
	//build.WithTable("sys_dept d").
	//	WithJoin("join sys_role_dept rd ON d.id = rd.dept_id")
	if name != "" {
		build.WithWhere("name like ?", "%"+name+"%")
	}
	if status != "" {
		_enable, _ := strconv.Atoi(status)
		build.WithWhere("status = ?", _enable)
	}
	build.WithPagination(ctx).WithDataScope(ctx, "d", "u").WithOrder("pid", "sort")
	if build.Paginate == nil {
		return d.FindByBuilder(ctx, build)
	} else {
		return d.ListByBuilder(ctx, build)
	}
}

// CheckStrictly选择联动
func (d *DeptRepo) DeptIdsByRoleId(ctx context.Context, roleId uint64, checkStrictly bool) ([]uint64, error) {
	var deptIds = make([]uint64, 0)
	build := base.NewQueryBuilder()
	build.WithTable("sys_dept d").WithSelect("", "id").
		WithJoin("join sys_role_dept rd ON d.id = rd.dept_id").
		WithWhere("rd.role_id=?", roleId).
		WithOrder("d.pid", "d.sort")
	if checkStrictly {
		build.WithWhere("d.id not in (select pid from sys_dept d left join sys_role_dept rd ON d.id=rd.dept_id where rd.role_id=?)", roleId)
	}
	err := d.WithQueryBuilder(ctx, build).Find(&deptIds).Error
	return deptIds, err
}
