/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc CasbinAdapter
 */

package casbin

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/pkg/log"
	"strings"
)

var _ persist.Adapter = (*CasbinAdapter)(nil)

type CasbinAdapter struct {
	UserRepo *repo.UserRepo
	RoleRepo *repo.RoleRepo
}

func (a *CasbinAdapter) LoadPolicy(model model.Model) error {
	ctx := context.Background()
	err := a.LoadRolePolicy(ctx, model)
	if err != nil {
		log.L.Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}
	err = a.LoadUserPolicy(ctx, model)
	if err != nil {
		log.L.Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}
	return nil
}

func (a *CasbinAdapter) SavePolicy(model model.Model) error {
	return nil
}

func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

var CasbinAdapterSet = wire.NewSet(wire.Struct(new(CasbinAdapter), "*"), wire.Bind(new(persist.Adapter), new(*CasbinAdapter)))

// Load role policy (p,role_id,path,method)
func (a *CasbinAdapter) LoadRolePolicy(ctx context.Context, m model.Model) error {

	builder := base.NewQueryBuilder()
	builder.WithPreload("Menus", "type != ?", "M").
		WithJoin("left join sys_user_role as ur on ur.role_id=sys_role.id").
		WithJoin("left join sys_user as u on u.id=ur.user_id").
		WithJoin("left join sys_dept as d on d.id=u.dept_id").
		WithWhere("sys_role.status=1").
		WithDataScope(ctx, "d", "u").
		WithPagination(ctx)

	roles, err := a.RoleRepo.FindAllByBuilder(ctx, builder)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return nil
	}
	for _, r := range roles {
		for _, menu := range r.Menus {
			// 1. Process permission identifier (e.g., "sys:user:list") to generate policy by convention
			if menu.Permission != nil && *menu.Permission != "" {
				key := *menu.Permission
				parts := strings.Split(key, ":")
				if len(parts) == 3 {
					resource := parts[1]
					action := parts[2]

					var method, path string
					// Convention: API path is /api/{resource}
					path = "/api/" + resource

					switch action {
					case "list", "query", "get":
						method = "GET"
					case "add", "create":
						method = "POST"
					case "edit", "update":
						method = "PUT"
					case "delete", "remove":
						method = "DELETE"
					}

					if method != "" {
						line := fmt.Sprintf("p,%d,%s,%s", r.ID, path, method)
						persist.LoadPolicyLine(line, m)
					}
				}
			}

			//if menu.Path != "" {
			//	line := fmt.Sprintf("p,%d,%s,%s", r.ID, menu.Path, "GET")                                                                                                  │
			//	persist.LoadPolicyLine(line, m)
			//}

			// 2. Process explicitly bound APIs (for buttons)
			for _, api := range menu.Apis {
				line := fmt.Sprintf("p,%d,%s,%s", r.ID, api.Path, api.Method)
				persist.LoadPolicyLine(line, m)
			}
		}
	}
	return nil
}

// Load user policy (g,user_id,role_id)
func (a *CasbinAdapter) LoadUserPolicy(ctx context.Context, m model.Model) error {
	var users []*entity.User
	err := a.UserRepo.DB.WithContext(ctx).Preload("Roles").Where("status = ?", 1).Find(&users).Error
	if err != nil {
		return err
	}

	for _, u := range users {
		for _, r := range u.Roles {
			line := fmt.Sprintf("g,%d,%d", u.ID, r.ID)
			persist.LoadPolicyLine(line, m)
		}
	}
	return nil
}
