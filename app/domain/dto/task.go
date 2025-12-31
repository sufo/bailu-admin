/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package dto

type TaskParams struct {
	Name   string `json:"name" query:"name,like"`
	Group  string `json:"group" query:"group, -"`
	Status string `json:"status" query:"status"`
}

type TaskLogParams struct {
	//TaskName  string `json:"taskName" query:"task_name,like"`
	//TaskGroup string `json:"TaskGroup" query:"task_group"`
	Status    string `json:"status" query:"status"`
	BeginDate string `form:"beginDate" query:"start_time,between endDate" binding:"omitempty,len=8"` //omitempty表示可选，存在则继续向后校验 YYYYMMDD
	EndDate   string `form:"endDate" binding:"omitempty,len=8"`
}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}
