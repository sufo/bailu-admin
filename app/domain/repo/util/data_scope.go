/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package util

import (
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/global/consts"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// dataScope
func DataScope(ctx context.Context, db *gorm.DB, deptAlias string, userAlias string) *gorm.DB {
	if deptAlias == "" {
		deptAlias = entity.DeptTN
	}
	if userAlias == "" {
		userAlias = entity.UserTN
	}
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if user == nil {
		log.L.Warn("current user not found")
		db.Where(fmt.Sprintf("%s.dept_id = 0", deptAlias)) //什么也查不到
		return db
	}
	// 如果是超级管理员，则不过滤数据
	if user.IsSuper() {
		return db
	}
	var condition []string
	var sqlBuilder strings.Builder
	for _, role := range user.Roles {
		dataScope := role.DataScope
		if consts.DATA_SCOPE_CUSTOM != dataScope && utils.ContainsInSlice(condition, dataScope) {
			continue
		}
		if consts.DATA_SCOPE_ALL == dataScope {
			sqlBuilder = strings.Builder{}
			break
		} else if consts.DATA_SCOPE_CUSTOM == dataScope {
			sql := fmt.Sprintf(" OR sys_dept.dept_id IN ( SELECT dept_id FROM sys_role_dept WHERE role_id = %d ) ", role.GetID())
			sqlBuilder.WriteString(sql)
		} else if consts.DATA_SCOPE_DEPT == dataScope {
			sql := fmt.Sprintf(" OR %s.dept_id = %d ", deptAlias, user.DeptId)
			sqlBuilder.WriteString(sql)
		} else if consts.DATA_SCOPE_DEPT_AND_CHILD == dataScope {
			sql := fmt.Sprintf(" OR %s.dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = %d or find_in_set( %d , ancestors ) )", deptAlias, user.DeptId, user.DeptId)
			sqlBuilder.WriteString(sql)
		} else if consts.DATA_SCOPE_SELF == dataScope {
			sql := fmt.Sprintf(" OR %s.user_id = %d ", userAlias, user.ID)
			sqlBuilder.WriteString(sql)
		}
		condition = append(condition, dataScope)
	}
	// 多角色情况下，所有角色都不包含传递过来的权限字符, 则不查询任何数据
	if condition == nil || len(condition) == 0 {
		sqlBuilder.WriteString(fmt.Sprintf("%s.dept_id = 0", deptAlias))
	}

	sqlString := sqlBuilder.String()
	if sqlString != "" {
		//sqlString = " AND (" + sqlString[4:] + ")"
		//return db.Exec(sqlString)
		db.Where(db.Where(sqlString[4:]))
	}
	return db
}
