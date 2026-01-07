/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 定时任务
 */

package cron

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/service/sys"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils"
	"time"
)

import (
	"github.com/sufo/bailu-admin/app/domain/repo"
	"sync"
)

const JobPrefix = "JOB:"

// 暂时没有放开配置
const RetryInterval = 5 //失败重试间隔（s）
const RetryTimes = 3    //失败重试次数

// 任务计数
type taskCount struct {
	wg   sync.WaitGroup
	exit chan struct{}
}

func (tc *taskCount) i() {}

func (tc *taskCount) Add() {
	tc.wg.Add(1)
}

func (tc *taskCount) Done() {
	tc.wg.Done()
}

func (tc *taskCount) Exit() {
	tc.wg.Done()
	<-tc.exit
}

func (tc *taskCount) Wait() {
	tc.Add()
	tc.wg.Wait()
	close(tc.exit)
}

// 记录正在执行的任务
// 任务ID作为Key
type Instance struct {
	m sync.Map
}

// 是否有任务处于运行中
func (i *Instance) has(key uint64) bool {
	_, ok := i.m.Load(key)

	return ok
}

func (i *Instance) add(key uint64) {
	i.m.Store(key, struct{}{})
}

func (i *Instance) done(key uint64) {
	i.m.Delete(key)
}

type CronTask struct {
	noticeSrv   *sys.NoticeService
	taskRepo    *repo.TaskRepo
	taskLog     *JobLog
	taskCount   *taskCount
	cron        gocron.Scheduler
	runInstance *Instance //记录运行job
}

func NewCronTask(TaskRepo *repo.TaskRepo, TaskLogRepo *repo.TaskLogRepo, NoticeSrv *sys.NoticeService) *CronTask {
	cron, err := gocron.NewScheduler()
	if err != nil {
		return nil
	}
	return &CronTask{
		noticeSrv: NoticeSrv,
		taskRepo:  TaskRepo,
		taskLog:   &JobLog{TaskLogRepo},
		//robfig/cron 秒级操作、函数没执行完就跳过本次函数
		//cron: goCron.New(goCron.WithSeconds(), goCron.WithChain(goCron.SkipIfStillRunning(goCron.DefaultLogger))),
		cron: cron,
		taskCount: &taskCount{
			wg:   sync.WaitGroup{},
			exit: make(chan struct{}),
		},
		runInstance: &Instance{},
	}
}

func (s *CronTask) createTask(ctx context.Context, task *entity.Task) gocron.Task {
	return gocron.NewTask(s.createJob(ctx), task)
}

//func jobName(taskId uint64) string {
//	return JobPrefix + utils.Strval(taskId)
//}
//func getIdByJobName(name string) uint64 {
//	t, err := utils.ToUint[uint64](strings.Split(name, JobPrefix)[1])
//	if err != nil {
//		log.L.Error("getIdByJobName failed", name, err)
//		return 0
//	}
//	return t
//}

// 仅funJob需要ctx
func (s *CronTask) Add(ctx context.Context, task *entity.Task) {
	cronName := utils.Strval(task.ID)

	jobOpts := []gocron.JobOption{
		gocron.WithTags(cronName),        //设置tag
		gocron.WithName(string(task.ID)), //设置job name
	}
	if task.Concurrent == 2 { //禁止该任务并发
		jobOpts = append(jobOpts, gocron.WithSingletonMode(gocron.LimitModeWait))
	}
	//创建Job
	job, err := s.cron.NewJob(gocron.CronJob(task.CronExpression, true),
		s.createTask(ctx, task),
		jobOpts...,
		//gocron.WithEventListeners(
		//	gocron.BeforeJobRuns(
		//		func(jobID uuid.UUID, jobName string) {
		//			s.taskCount.Add()
		//			defer s.taskCount.Done()
		//		}),
		//	gocron.AfterJobRuns(
		//		func(jobID uuid.UUID, jobName string) {
		//			// do something after the job completes
		//			s.afterExecJob(task)
		//		},
		//	),
		//),
	)

	//更新job id
	task.EntryId = job.ID().String()
	s.taskRepo.UpdateEntryId(task.ID, "entry_id", task.EntryId)
	if err != nil {
		log.L.Error(fmt.Errorf("执行任务：", task.ID, err))
	}
}

// 批量添加任务
func (s *CronTask) BatchAdd(ctx context.Context, tasks []*entity.Task) {
	for _, item := range tasks {
		s.RemoveAndAdd(ctx, item)
	}
}

func (s *CronTask) RemoveByTag(taskTag ...string) {
	s.cron.RemoveByTags(taskTag...)
}

func (s *CronTask) Remove(entryId string) {
	//转UUID
	_entryId := []byte(entryId)
	var result uuid.UUID
	copy(result[:], _entryId)

	s.cron.RemoveJob(result)
}

func (s *CronTask) RemoveAndAdd(ctx context.Context, task *entity.Task) {
	cronName := utils.Strval(task.ID)
	s.Remove(cronName)
	s.Add(ctx, task)
}

// 直接运行任务
func (s *CronTask) Run(ctx context.Context, task *entity.Task) {
	go s.createJob(ctx)(task)
}

func (s *CronTask) Stop() {
	defer func() { _ = s.cron.Shutdown() }()
	_ = s.cron.StopJobs()
	s.taskCount.Exit()
}

func (s *CronTask) Start(ctx context.Context) {
	s.cron.Start()
	go s.taskCount.Wait()

	var tempTask = make([]*entity.Task, 0)
	err := s.taskRepo.Where(ctx, "status=1").FindInBatches(&tempTask, 50,
		func(tx *gorm.DB, batch int) error {
			s.BatchAdd(ctx, tempTask)
			return nil
		}).Error
	if err != nil {
		log.L.Fatal("cron initialize tasks count err", zap.Error(err))
	}
}

func (s *CronTask) createJob(ctx context.Context) func(task *entity.Task) {
	return func(task *entity.Task) {
		var j JobProxy
		if task.Protocol == "FUNC" {
			j = (&ExecJob{ctx, *task})
		} else if task.Protocol == "HTTP" {
			j = &HttpJob{*task}
		} else {
			j = &ShellJob{*task}
		}

		s.taskCount.Add()
		defer s.taskCount.Done()

		taskLogId := s.beforeExecJob(task)
		if taskLogId <= 0 { //说明取消了本次执行，因为之前的任务还没有结束
			return
		}
		//记录正在运行job
		s.runInstance.add(task.ID)
		defer s.runInstance.done(task.ID)

		log.L.Infof("开始执行任务#%s#命令-%s", task.Name, task.InvokeTarget)
		taskResult := execJob(j, taskLogId)

		//更新执行时间
		s.taskRepo.UpdateColumn(task.ID, "last_exec_time", time.Now())

		log.L.Infof("任务完成#%s#命令-%s", task.Name, task.InvokeTarget)
		s.afterExecJob(ctx, task, taskResult, taskLogId)
	}
}

// 任务执行前操作
func (s *CronTask) beforeExecJob(taskModel *entity.Task) uint64 {
	//if taskModel.Multi == 0 && runInstance.has(taskModel.Id) {
	//	createTaskLog(taskModel, models.Cancel)
	//	return
	//}
	// 单实例运行
	if s.runInstance.has(taskModel.ID) {
		_, err := s.taskLog.CreateTaskLog(taskModel, entity.Cancel)
		if err != nil {
			log.L.Error("任务开始执行#写入任务日志失败-", err)
		}
		return 0
	}

	taskLogId, err := s.taskLog.CreateTaskLog(taskModel, entity.Running)
	if err != nil {
		log.L.Error("任务开始执行#写入任务日志失败-", err)
		return 0
	}
	log.L.Debugf("任务命令-%s", taskModel.InvokeTarget)
	return taskLogId
}

// 任务执行后置操作
func (s *CronTask) afterExecJob(ctx context.Context, taskModel *entity.Task, taskResult dto.TaskResult, taskLogId uint64) {
	_, err := s.taskLog.updateTaskLog(taskLogId, taskResult)
	if err != nil {
		log.L.Error("任务结束#更新任务日志失败-", err)
	}
	//notice
	go s.notify(ctx, taskModel, taskResult)
}

// 执行具体任务
func execJob(job JobProxy, taskLogId uint64) dto.TaskResult {
	defer func() {
		if err := recover(); err != nil {
			log.L.Error("panic#service/task.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var count int8 = 0
	var output string
	var err error
LOOP:
	if count < RetryTimes {
		output, err = job.Run()
		if err != nil {
			count = count + 1
			if count <= RetryTimes {
				//重试间隔
				time.Sleep(time.Duration(count) * RetryInterval * time.Second)
				goto LOOP
			}
		}
	}
	return dto.TaskResult{Result: output, Err: err, RetryTimes: count}
}

// 通知
// TODO
func (s *CronTask) notify(ctx context.Context, taskModel *entity.Task, taskResult dto.TaskResult) {
	strategy := taskModel.NotifyStrategy
	switch strategy {
	case entity.FAILURE_NOTIFY:
		if taskResult.Err != nil {
			s.noticeSrv.TaskNotice(ctx, taskModel, taskResult)
		}
	case entity.FINISHED_NOTIFY:
		s.noticeSrv.TaskNotice(ctx, taskModel, taskResult)
	case entity.MATCH_NOTIFY:
		//TODO
	}
}
