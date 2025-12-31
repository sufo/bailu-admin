/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 操作日志表
 */

package entity

import "bailu/utils/types"

var _ IModel = (*DictItem)(nil)

type Log struct {
	ID           uint64         `json:"id,string" gorm:"primarykey"`
	Name         string         `json:"name" gorm:"comment:操作模块"`
	Action       int            `json:"action" gorm:"type:tinyint(1);comment:操作类型（业务类型（0其它 1新增 2修改 3删除））"`
	Method       string         `json:"method" gorm:"comment:请求方法"`
	OperatorType uint           `json:"operateType" gorm:"type:tinyint(1);comment:操作类别（0其它 1后台用户 2手机端用户）"`
	OperName     string         `json:"operName" gorm:"comment:操作人员姓名"`
	DeptName     string         `json:"deptName" gorm:"comment:操作人员部门名称"`
	OperUrl      string         `json:"operUrl" gorm:"comment:操作url"`
	OperIp       string         `json:"operIp" gorm:"comment:操作地址"`
	OperLoc      string         `json:"operLoc" gorm:"comment:操作地点"`
	OperParam    string         `json:"operParam" gorm:"comment:请求参数"`
	Result       string         `json:"result" gorm:"comment:返回结果"`
	Status       uint           `json:"status" gorm:"type:tinyint(1);comment:操作状态（0正常 1异常）"`
	ErrorMsg     string         `json:"errorMsg" gorm:"comment:错误消息"`
	OperTime     types.JSONTime `json:"operTime" gorm:"comment:操作时间"`
}

var LogTN = "sys_log"

func (i Log) GetID() uint64 {
	return i.ID
}

func (Log) TableName() string {
	return LogTN
}
