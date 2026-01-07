/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package entity

import (
	"github.com/sufo/bailu-admin/global/consts"
	"encoding/json"
	"time"
)

type OnlineUserDto struct {
	Username string `json:"username"`
	ID       uint64 `json:"id"`
	NickName string `json:"nickName"`
	Token    string `json:"-"` //用户唯一标识（uuid），不是jwt
	Browser  string `json:"browser"`
	Os       string `json:"os"`
	Ip       string `json:"ip"`
	Addr     string `json:"addr"`
	//RoleIds  []uint64 `json:"-"`
	Permissions []string `json:"permissions"`
	Roles       []Role   `json:"roles"`
	//Roles  []uint64 `json:"roles"`
	DeptId   uint64 `json:"-"`
	DeptName string `json:"deptName"`
	/**
	 * 登录时间
	 */
	LoginTime time.Time `json:"loginTime"`
}

func (o *OnlineUserDto) IsSuper() bool {
	for _, e := range o.Roles {
		if consts.SUPER_ROLE_ID == e.ID {
			return true
		}
	}
	return false
}

// redis set must implement
// notice: OnlineUserDto must be not point
func (onlineUser OnlineUserDto) MarshalBinary() ([]byte, error) {
	return json.Marshal(onlineUser)
}
