/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 系统用户
 */

package entity

import (
	"github.com/sufo/bailu-admin/global/consts"
	"github.com/sufo/bailu-admin/utils/types"
	"strings"
	"time"
)

var _ IModel = (*User)(nil)

// `json:"xxx,omitempy"`  // 如果为类型零值或空值，序列化时忽略该字段
type User struct {
	//ID       uint64 `json:"id,string" gorm:"primarykey"`
	ID       uint64 `json:"id" gorm:"primarykey"`
	Username string `json:"username" gorm:"not null;size:30;comment:用户名;not null;unique" binding:"required,lte=30"`
	//pwd         string   `json:"password" gorm:"-" binding:"required"` //接受用户提交的密码
	Password      string                 `json:"-" gorm:"not null;size:100;comment:用户登录密码"`
	NickName      string                 `json:"nickName,omitempty" gorm:"size:30;comment:用户昵称" binding:"required"`
	UserType      string                 `json:"userType,omitempty" gorm:"default：'00';size:2;comment:用户类型 默认00：系统用户"`
	Profile       string                 `json:"profile,omitempty" gorm:"size:200;comment:用户简介"`
	Email         string                 `json:"email,omitempty" gorm:"size:100;comment:邮箱;unique" binding:"omitempty,email"`
	DialCode      string                 `json:"dialCode,omitempty" gorm:"size:10;default:'86';comment:地区（国家）编码"`
	Phone         string                 `json:"phone,omitempty" gorm:"size:11;comment:手机号;unique"`
	Sex           *uint8                 `json:"sex,omitempty" gorm:"type:tinyint(4);default:0;comment:0未知 1男 2女"`
	Avatar        string                 `json:"avatar,omitempty" gorm:"comment:用户头像"`
	DeptId        *uint64                `json:"deptId" gorm:"default:null;comment:部门ID"`
	DeptName      string                 `json:"deptName" gorm:"-"`
	HomePath      string                 `json:"homePath" gorm:"size:200;comment:用户首页"`
	Permissions   []string               `json:"permissions" gorm:"-"` //用户权限集合
	Roles         types.EmptySlice[Role] `json:"roles" gorm:"many2many:sys_user_role;"`
	Posts         types.EmptySlice[Post] `json:"posts" gorm:"many2many:sys_user_post"`
	Ip            string                 `json:"ip,omitempty" gorm:"size:20;comment:用户最后登录ip"`
	Status        uint8                  `json:"status" gorm:"type:tinyint(1);comment:是否启用(1:启用 2:禁用);default:0"`
	Remark        string                 `json:"remark,omitempty" gorm:"size:200;comment:备注信息"`
	LastLoginTime time.Time         `json:"lastLoginTime" gorm:"default null,comment:登录时间"`
	//创建时使用
	//RoleIds types.Uint64EmptySlice `json:"roleIds" gorm:"-"`
	//PostIds types.Uint64EmptySlice `json:"postIds" gorm:"-"`
	RoleIds []uint64 `json:"roleIds,omitempty" gorm:"-"`
	PostIds []uint64 `json:"postIds,omitempty" gorm:"-"`

	BaseEntity
}

var UserTN = "sys_user"

func (User) TableName() string {
	return UserTN
}

func (u *User) Alias() string {
	return strings.Split(u.TableName(), "_")[1]
}

func (u User) GetID() uint64 {
	return u.ID
}

// ///////////中间表//////////////
type UserRole struct {
	UserId uint64 `gorm:"column:user_id"`
	RoleId uint64 `gorm:"column:role_id"`
}

func (UserRole) TableName() string {
	return "sys_user_role"
}

func UserRoleTableName() string {
	return "sys_user_role" //和上面User结构体定义的要一致
}

//func (UserRole) TableName() string {
//	return "sys_user_role"
//}

type UserPost struct {
	UserId uint64
	PostId uint64
}

func (UserPost) TableName() string {
	return "sys_user_post"
}

func (o *User) IsSuper() bool {
	for _, e := range o.Roles {
		if consts.SUPER_ROLE_ID == e.ID {
			return true
		}
	}
	return false
}

func IsSuper(userId uint64) bool {
	return consts.SUPER_USER_ID == userId
}
