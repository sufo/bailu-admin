/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package dto

import "bailu/global/consts"

type RoleParams struct {
	Name    string `json:"name" form:"name" query:"name,like"`
	RoleKey string `json:"roleKey" form:"roleKey" query:"role_key"`
	Status  *int   `json:"status" form:"status" query:"status"`
	//ltefield=EndDate（小与等与EndDate）
	//BeginDate string `form:"beginDate" binding:"omitempty,len=8,ltefield=EndDate" query:"createAt,between endDate"`
	//EndDate   string `form:"endDate" binding:"omitempty,len=8"`
	BeginDate string `json:"beginDate" form:"beginDate" binding:"omitempty,ltefield=EndDate" query:"created_at,between EndDate"`
	EndDate   string `json:"endDate" form:"endDate" binding:"omitempty"`
	//CreatedAt []string `form:"createdAt between" `
}

type Role struct {
	ID             uint64   `json:"id,string" form:"id,string"`
	Name           string   `json:"name" form:"name" binding:"required"`
	RoleKey        string   `json:"roleKey" form:"roleKey" binding:"required"`
	DataScope      string   `json:"dataScope" form:"dataScope"`
	Menus          []uint64 `json:"menus" form:"menus"`
	Depts          []uint64 `json:"depts" form:"depts"` //数据权限
	Sort           uint     `json:"sort" form:"sort"`
	Status         uint8    `json:"status" form:"status" binding:"required"`
	Remark         string   `json:"remark" form:"remark"`
	IsMenusChanged bool     `json:"isMenusChanged" form:"isMenusChanged"`
}

type RoleDataScope struct {
	ID        uint64   `json:"id,string" form:"id,string" binding:"required"`
	DataScope string   `json:"dataScope" form:"dataScope" binding:"required"`
	DeptIds   []uint64 `json:"deptIds" form:"deptIds"` //数据权限
}

func (r *Role) IsSuper() bool {
	return consts.SUPER_ROLE_ID == r.ID
}
