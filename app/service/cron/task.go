/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 定时任务
 */

package cron

import (
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/domain/vo"
	base2 "github.com/sufo/bailu-admin/app/service/base"
	"github.com/sufo/bailu-admin/utils"
	"github.com/sufo/bailu-admin/utils/time"
	"context"
	"github.com/adhocore/gronx"
	"github.com/google/wire"
)

var TaskSet = wire.NewSet(wire.Struct(new(TaskOption), "*"), NewTaskService)

//type TaskService struct {
//	TaskRepo    *repo.TaskRepo
//	TaskLogRepo *repo.TaskLogRepo
//	TransRepo   *repo.Trans
//	Cron        *CronTask
//}

type TaskOption struct {
	TaskRepo    *repo.TaskRepo
	TaskLogRepo *repo.TaskLogRepo
	TransRepo   *repo.Trans
	Cron        *CronTask
}
type TaskService struct {
	base2.BaseService[entity.Task]
	TaskOption
}

func NewTaskService(opt TaskOption) *TaskService {
	return &TaskService{base2.BaseService[entity.Task]{opt.TaskRepo.Repository}, opt}
}

func (t *TaskService) Create(ctx context.Context, task *entity.Task) error {
	err := t.TaskRepo.Create(ctx, task)
	if err != nil {
		return err
	}
	t.Cron.Add(ctx, task)
	return nil
}

// 查询所有符合条件的，无dataScope限制
func (p *TaskService) List(ctx context.Context, Task dto.TaskParams) (*resp.PageResult[entity.Task], error) {
	builder := base.NewQueryBuilder()
	builder.WithWhereStruct(Task).
		WithPagination(ctx)
	if result, err := p.TaskRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		//处理task.nextTime
		tasks := result.(*resp.PageResult[entity.Task])
		if len(tasks.List) > 0 {
			for _, item := range tasks.List {
				nextTime, err := gronx.NextTick(item.CronExpression, true)
				if err == nil {
					item.NextTime = nextTime.Format(time.CSTLayout)
				}
			}
		}

		return result.(*resp.PageResult[entity.Task]), err
	}
}

func (t *TaskService) Update(ctx context.Context, Task *entity.Task) (*entity.Task, error) {
	err := t.TaskRepo.Update(ctx, Task)
	if err != nil {
		return nil, err
	}

	t.Cron.RemoveAndAdd(ctx, Task)
	return Task, nil
}

func (t *TaskService) Delete(ctx context.Context, ids []uint64) error {
	err := t.TaskRepo.Delete(ctx, ids)
	if err != nil {
		return err
	}

	//删除job
	tags := make([]string, 0)
	for _, id := range ids {
		tags = append(tags, utils.Strval(id))
	}
	t.Cron.RemoveByTag(tags...)
	return nil
}

func (t *TaskService) DeleteLogs(ctx context.Context, ids []uint64) error {
	return t.TaskLogRepo.Delete(ctx, ids)
}

func (t *TaskService) FindById(ctx context.Context, id uint64) (*entity.Task, error) {
	return t.TaskRepo.FindById(ctx, id)
}

func (t *TaskService) ChangeStatus(ctx context.Context, id uint64, status uint8) error {
	err := t.TaskRepo.Where(ctx, "id=?", id).UpdateColumn("status", status).Error
	if err != nil {
		return err
	}

	task, err := t.TaskRepo.FindById(ctx, id)
	if status == entity.Enabled {
		if err != nil {
			return err
		}
		t.Cron.Add(ctx, task)
	} else {
		t.Cron.Remove(task.EntryId)
	}
	return nil
}

func (t *TaskService) Execute(ctx context.Context, id uint64) error {
	info, err := t.TaskRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	t.Cron.Run(ctx, info)
	return nil
}

// func job
func (t *TaskService) FindFuncJobsOptions(ctx context.Context) []*vo.KV {
	options := GetFuncJobs()
	return options
}

func (t *TaskService) TaskLog(ctx context.Context, taskId uint64, params dto.TaskLogParams) (*resp.PageResult[entity.TaskLog], error) {
	//return t.TaskLogRepo.FirstBy(ctx, "task_id=?", id)
	builder := base.NewQueryBuilder()
	builder.WithWhereStruct(params).
		WithWhere("task_id", taskId).
		WithPagination(ctx)
	if result, err := t.TaskLogRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		//处理task.nextTime
		taskLogs := result.(*resp.PageResult[entity.TaskLog])
		return taskLogs, err
	}
}
