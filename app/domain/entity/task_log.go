/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 任务日志
 */

package entity

import "time"

var _ IModel = (*TaskLog)(nil)

type TaskLog struct {
	ID           uint64         `gorm:"primarykey"`
	TaskId       uint64         `json:"taskId" gorm:"not null;comment:任务id"`
	TaskName     string         `json:"taskName" gorm:"size:64;not null;comment:任务名称"`
	TaskGroup    string         `json:"taskGroup" gorm:"size:64;not null;comment:任务组名"`
	InvokeTarget string         `json:"invokeTarget" gorm:"size:500;not null;comment:调用目标字符串"`
	Status       uint8          `json:"status" gorm:"type:tinyint(1);default:1;执行状态（ 1:执行中  2:执行完毕 3:执行失败 4:任务取消(上次任务未执行完成) 5:异步执行）"`
	Result       string         `json:"Result" gorm:"size:200;comment:执行输出结果"`
	ExceptInfo   string         `json:"except_info" gorm:"size:2000;comment:异常信息"`
	RetryTimes   int            `json:"retryTimes" gorm:"comment:重试次数"`
	StartTime    time.Time `json:"startTime" gorm:"comment:开始时间"`
	StopTime     time.Time `json:"stopTime" gorm:"autoUpdateTime;comment:停止时间"` //自动更新时间
	TotalTime    int            `json:"totalTime" gorm:"-"`                              // 执行总时长
}

var TaskLogTN = "sys_task_log"

func (t TaskLog) TableName() string {
	return TaskLogTN
}

func (t TaskLog) GetID() uint64 {
	return t.ID
}
