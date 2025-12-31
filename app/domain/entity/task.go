/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 定时任务调度表
 */

package entity

import "bailu/utils/types"

const (
	NO_NOTIFY       = 1
	FAILURE_NOTIFY  = 2
	FINISHED_NOTIFY = 3
	MATCH_NOTIFY    = 4
)

//	type Task struct {
//		ID             uint64 `json:"id" gorm:"primarykey"`
//		Name           string `json:"name" gorm:"size:64;default:'';comment:任务名称"`
//		Group          string `json:"group" gorm:"size:64;default:'DEFAULT';comment:任务组名"`
//		InvokeTarget   string `json:"invokeTarget" gorm:"size:500;not null;size:64;调用目标字符串"`
//		CronExpression string `json:"cronExpression" gorm:"size:255;default:'';comment:cron执行表达式"`
//		MisfirePolicy  uint8  `json:"misfirePolicy" gorm:"type:tinyint(1);comment:cron计划策略;default:1"`
//		Concurrent     uint8  `json:"concurrent" gorm:"type:tinyint(1);default:1;comment:是否并发执行（1允许 2禁止）"`
//		Status         uint8  `json:"status" gorm:"type:tinyint(1);default:1;comment:任务状态（1正常 2暂停）"`
//		BaseEntity
//	}
//
// tinyint默认tinyint(4)
type Task struct {
	ID             uint64 `json:"id" gorm:"primarykey"`
	Name           string `json:"name" gorm:"size:64;default:'';comment:任务名称" binding:"required"`
	Group          string `json:"group" gorm:"size:64;default:'DEFAULT';comment:任务组名" binding:"required"`
	Protocol       string `json:"protocol" gorm:"size:20;comment:执行方式 FUNC:函数 HTTP:http; SHELL:shell" binding:"required"`
	CronExpression string `json:"cronExpression" gorm:"size:64;default:'';comment:cron执行表达式" binding:"required"`
	InvokeTarget   string `json:"invokeTarget" gorm:"size:255;not null;comment:调用目标" binding:"required"`
	Args           string `json:"args" gorm:"size:255;comment:目标参数"`
	HttpMethod     string `json:"httpMethod" gorm:"size:10;default:'get';comment:http 请求方式 get post put patch delete等;"`
	Concurrent     uint8  `json:"concurrent" gorm:"type:tinyint(4) unsigned;default:1;comment:是否并发执行（1允许 2禁止）"`
	//MisfirePolicy  uint8  `json:"misfirePolicy" gorm:"type:tinyint(4) unsigned;default:1;comment:执行策略（1默认 2执行一次 3立即执行）"`
	//Timeout             int    `json:"timeout" gorm:"type:int unsigned;comment:超时时间(单位:秒);"`
	//RetryTimes          uint8  `json:"retryTimes" gorm:"type:tinyint;default:0;comment:重试次数;"`
	//RetryInterval       int    `json:"retryInterval" gorm:"type:int;comment:重试间隔(单位:秒);"`
	Status         uint8  `json:"status" gorm:"type:tinyint(4) unsigned;default:0;comment:是否启用（1正常 2停用）"`
	EntryId        string `json:"entryId" gorm:"size:36;comment:job启动时返回的id"`
	NotifyStrategy uint8  `json:"notifyStrategy"  binding:"required" gorm:"type:tinyint unsigned;comment:执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知;"`
	NotifyChannel  string `json:"notifyChannel" gorm:"default:'web';size:20;default:'web';comment:通知方式：web,app,mail,sms等"`
	//NotifyReceiverEmail string         `json:"notifyReceiverEmail" gorm:"comment:接收者邮箱地址(多个用,分割)"`
	NotifyKeyword string         `json:"notifyKeyword" gorm:"comment:通知匹配关键字(多个用,分割)"`
	Remark        string         `json:"remark" gorm:"comment:备注"`
	LastExecTime  types.JSONTime `json:"lastExecTime" gorm:"comment:最近一次执行时间;default:null"`
	NextTime      string         `json:"nextTime" gorm:"-"` //下次执行时间，通过分析cron表达式得出
	BaseEntity
}

const (
	TaskHTTPMethodGet  uint8 = 1
	TaskHttpMethodPost uint8 = 2
)

var TaskTN = "sys_task"

func (t Task) TableName() string {
	return TaskTN
}

func (t Task) GetID() uint64 {
	return t.ID
}
