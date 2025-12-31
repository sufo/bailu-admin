/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 系统访问日志表
 */

package entity

import "bailu/utils/types"

type LoginInfo struct {
	ID uint64 `json:"id" gorm:"primarykey"`
	//UserId    uint64         `json:"userId" gorm:"not null;comment:用户名id;not null;" binding:"required"`
	Username  string         `json:"username" gorm:"not null;size:30;comment:用户名;not null;" binding:"required,lte=30"`
	Ip        string         `json:"ip" gorm:"size:255;comment:登录IP"`
	Addr      string         `json:"addr" gorm:"size:255;comment:登录IP地址"`
	Browser   string         `json:"browser" gorm:"size:50;comment:浏览器"`
	Os        string         `json:"os" gorm:"size:50;comment:操作系统"`
	Status    int            `json:"status" gorm:"type:tinyint(4);comment:登录状态（0成功 1失败）"`
	Msg       string         `json:"msg" gorm:"comment:提示消息"`
	LoginTime types.JSONTime `json:"loginTime" gorm:"comment:访问时间"`
}

var LoginInfoTN = "sys_login_info"

func (LoginInfo) TableName() string {
	return LoginInfoTN
}

func (l LoginInfo) GetID() uint64 {
	return l.ID
}
