/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package vo

import "github.com/sufo/bailu-admin/app/domain/entity"

type User struct {
	//ID          uint64   `json:"id,string"`
	ID          uint64   `json:"id"`
	Username    string   `json:"username"`
	Password    string   `json:"-"`
	NickName    string   `json:"nickName"`
	UserType    string   `json:"userType"`
	Profile     string   `json:"profile"`
	Email       string   `json:"email"`
	DialCode    string   `json:"dialCode"`
	Phone       string   `json:"phone"`
	Sex         int      `json:"sex"`
	Avatar      string   `json:"avatar"`
	DeptId      *uint64  `json:"deptId"`
	DeptName    string   `json:"deptName"`
	HomePath    string   `json:"homePath"`
	Roles       []KV     `json:"roles"`
	Posts       []KV     `json:"posts"`
	Permissions []string `json:"permissions"`
	Ip          string   `json:"ip"`
	entity.BaseEntity
}
