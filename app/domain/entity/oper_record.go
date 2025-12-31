/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package entity

import (
	"time"
)

var _ IModel = (*OperationRecord)(nil)

// 如果含有time.Time 请自行import time包
type OperationRecord struct {
	ID       uint64        `json:"id,string" gorm:"primarykey"`
	Ip       string        `json:"ip" form:"ip" gorm:"column:ip;comment:请求ip"` // 请求ip
	Location string        `json:"location" form:"location" gorm:"column:location;comment:操作地点"`
	Method   string        `json:"method" form:"method" gorm:"size:10;column:method;comment:请求方法"`     // 请求方法
	Path     string        `json:"path" form:"path" gorm:"column:path;comment:请求uri"`                    // 请求uri
	Status   int           `json:"status,omitempty" form:"status" gorm:"column:status;comment:http状态码"` // http状态码
	RespCode *int          `json:"respCode,omitempty" form:"respCode" gorm:"column:resp_code;comment: 逻辑响应码"`
	Latency  time.Duration `json:"latency" form:"latency" gorm:"column:latency;comment:延迟" swaggertype:"string"` // 延迟
	Agent    string        `json:"agent" form:"agent" gorm:"column:agent;comment:代理"`                            // 代理
	OS       string        `json:"os" gorm:"-"`
	Browser  string        `json:"browser" gorm:"-"`
	Msg      string        `json:"msg" form:"msg" gorm:"size:2000;column:msg;comment:响应信息"`
	Body     string        `json:"body" form:"body" gorm:"type:text;column:body;comment:请求Body"`            // 请求Body
	Resp     string        `json:"resp" form:"resp" gorm:"type:text;column:resp;comment:响应Body"`            // 响应Body
	OperId   uint64        `json:"operId" form:"operId" gorm:"column:oper_id;comment:用户id"`                 // 用户id
	OperName string        `json:"operName" form:"operName" gorm:"size:30;column:oper_name;comment:用户名称"` // 用户id
	Model
}

var OperRecordTN = "sys_oper_record"

func (OperationRecord) TableName() string {
	return OperRecordTN
}

func (d OperationRecord) GetID() uint64 {
	return d.ID
}
