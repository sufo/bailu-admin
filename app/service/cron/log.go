/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package cron

import (
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"time"
)

type JobLog struct {
	taskLogRepo *repo.TaskLogRepo
}

// 创建任务日志
func (j *JobLog) CreateTaskLog(taskModel *entity.Task, status uint8) (uint64, error) {
	taskLogModel := new(entity.TaskLog)
	taskLogModel.TaskId = taskModel.ID
	taskLogModel.TaskName = taskModel.Name
	taskLogModel.TaskGroup = taskModel.Group
	taskLogModel.InvokeTarget = taskModel.InvokeTarget
	taskLogModel.StartTime = time.Now()
	taskLogModel.Status = status
	return j.taskLogRepo.CreateTaskLog(taskLogModel)
}

// 更新任务日志
func (j *JobLog) updateTaskLog(taskLogId uint64, taskResult dto.TaskResult) (uint64, error) {
	result := taskResult.Result
	var status uint8
	if taskResult.Err != nil {
		status = entity.Failure
	} else {
		status = entity.Finish
	}
	err := j.taskLogRepo.UpdateTaskLog(taskLogId, map[string]any{
		"retry_times": taskResult.RetryTimes,
		"status":      status,
		"result":      result,
		"except_info": taskResult.Err.Error(),
	})
	return taskLogId, err
}
