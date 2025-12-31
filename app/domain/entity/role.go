/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 角色
 */

package entity

var _ IModel = (*Role)(nil)

type Role struct {
	//ID        uint64 `json:"id,string" gorm:"primarykey"`
	ID        uint64 `json:"id" gorm:"primarykey"`
	Name      string `json:"name" gorm:"size:30;not null;comment:角色名称"`
	RoleKey   string `json:"roleKey" gorm:"size:100;not null;comment:角色权限字符"`
	DataScope string `json:"dataScope" gorm:"type:tinyint(1);default:1;comment:数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5:仅本人）"`
	Menus     []Menu `json:"-" gorm:"many2many:sys_role_menu;comment:权限菜单"`
	Depts     []Dept `json:"-" gorm:"many2many:sys_role_dept;"` //数据权限
	SortAndStatus
	Remark string `json:"remark" gorm:"size:500;comment:描述"`
	BaseEntity
}

var RoleTN = "sys_role"
var RoleMenuTN = "sys_role_menu"
var RoleDeptTN = "sys_role_dept"

func (Role) TableName() string {
	return RoleTN
}

func (r Role) GetID() uint64 {
	return r.ID
}

// 联合主键
type RoleMenu struct {
	RoleId uint64 `json:"-" gorm:"primaryKey"`
	MenuId uint64 `json:"-" gorm:"primaryKey"`
}

func (RoleMenu) TableName() string {
	return RoleMenuTN
}

// 联合主键
type RoleDept struct {
	RoleId uint64 `json:"-" gorm:"primaryKey"`
	DeptId uint64 `json:"-" gorm:"primaryKey"`
}

func (RoleDept) TableName() string {
	return RoleDeptTN
}
