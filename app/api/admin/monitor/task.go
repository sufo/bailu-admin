/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 定时任务
 */

package monitor

import (
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/service/cron"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils"
	"github.com/sufo/bailu-admin/utils/page"
	"github.com/adhocore/gronx"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TaskApi struct {
	TaskSrv *cron.TaskService
}

func NewTaskApi(TaskSrv *cron.TaskService) *TaskApi {
	return &TaskApi{TaskSrv}
}

// @title 定时任务列表
// @Summary 定时任务列表
// @Description 按条件查询定时任务列表
// @Tags Task
// @Accept json
// @Produce json
// @Param query query dto.TaskParams false "查询条件"
// //////@Success 200 {object} resp.Response{data=resp.PageResult[entity.Task]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.Task]]
// @Router /api/task [get]
// @Security Bearer
func (t *TaskApi) Index(c *gin.Context) {
	var params dto.TaskParams
	err := c.ShouldBindQuery(&params)
	if err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	}
	page.StartPage(c)
	result, err := t.TaskSrv.List(c.Request.Context(), params)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 创建定时任务
// @Summary 创建定时任务
// @Description 创建定时任务
// @Tags Task
// @Accept json
// @Produce json
// @Param body body entity.Task false "task"
// @Success 200 {object} resp.Response[any]
// @Router /api/task [post]
// @Security Bearer
func (t *TaskApi) Create(c *gin.Context) {
	var task entity.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		ctx := c.Request.Context()
		//检查字段唯一性
		if !t.TaskSrv.CheckUnique(ctx, "name=?", task.Name) {
			resp.FailWithMsg(c, i18n.DefTr("api.existed", task.Name))
			return
		}

		//检查cron表达式
		if !gronx.New().IsValid(task.CronExpression) {
			resp.FailWithMsg(c, i18n.DefTr("api.cornInvalid"))
			return
		}

		err := t.TaskSrv.Create(c.Request.Context(), &task)
		if err != nil {
			resp.InternalServerError(c)
			c.Abort()
			return
		}
		resp.Ok(c)
	}
}

// @title 修改定时任务
// @Summary 修改定时任务
// @Description 修改定时任务
// @Tags Task
// @Accept json
// @Produce json
// @Param body body entity.Task false "task"
// ////////@Success 200 {object} resp.Response{data=entity.Task}
// @Success 200 {object} resp.Response[entity.Task]
// @Router /api/task [put]
// @Security Bearer
func (t *TaskApi) Edit(c *gin.Context) {
	var task entity.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		//检查字段唯一性
		if !t.TaskSrv.CheckUnique(c.Request.Context(), "id!=? and name=?", task.ID, task.Name) {
			resp.FailWithMsg(c, i18n.DefTr("api.existed", task.Name))
			return
		}
		//检查cron表达式
		if !gronx.New().IsValid(task.CronExpression) {
			resp.FailWithMsg(c, i18n.DefTr("api.cornInvalid"))
			return
		}

		task, err := t.TaskSrv.Update(c.Request.Context(), &task)
		if err != nil {
			resp.InternalServerError(c)
			return
		}
		resp.OKWithData(c, task)
	}
}

// @title 批量删除定时任务
// @Summary 批量删除定时任务
// @Description 批量删除定时任务
// @Tags Task
// @Accept json
// @Produce json
// @Param ids path string true "ids"
// @Success 200 {object} resp.Response[any]
// @Router /api/task/{ids} [delete]
// @Security Bearer
func (p *TaskApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	if err := p.TaskSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 获取任务信息
// @Summary 根据id获取单个任务信息
// @Description 根据id获取单个任务信息
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string false "id"
// //@Success 200 {object} resp.Response{data=entity.Task}
// @Success 200 {object} resp.Response[entity.Task]
// @Router /api/task/{id} [get]
func (p *TaskApi) Detail(c *gin.Context) {
	id := admin.ParseParamId(c, "id")
	data, err := p.TaskSrv.TaskRepo.FindById(c.Request.Context(), id)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	}
	resp.OKWithData(c, data)
}

// @title 批量删除定时任务
// @Summary 批量删除定时任务
// @Description 批量删除定时任务
// @Tags Task
// @Accept json
// @Produce json
// @Param ids path string true "ids"
// @Success 200 {object} resp.Response[any]
// @Router /api/task/log/{ids} [delete]
// @Security Bearer
func (p *TaskApi) DestroyLogs(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	if err := p.TaskSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 任务执行日志
// @Summary 任务执行日志
// @Description 任务执行日志列表
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string false "taskId"
// @Param id path string false "id"
// @Param query query dto.TaskLogParams false "查询参数"
// //@Success 200 {object} resp.Response{data=entity.TaskLog]
// @Success 200 {object} resp.Response[entity.TaskLog]
// @Router /api/task/{id}/logs [get]
// @Security Bearer
func (p *TaskApi) Logs(c *gin.Context) {
	id := admin.ParseParamId(c, "id")
	var params dto.TaskLogParams
	err := c.ShouldBindQuery(&params)
	if err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	}
	page.StartPage(c)
	data, err := p.TaskSrv.TaskLog(c.Request.Context(), id, params)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	}
	resp.OKWithData(c, data)
}

// @title 任务立即执行
// @Summary 任务立即执行
// @Description 任务立即执行
// @Tags Task
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string false "id"
// @Success 200 {object} resp.Response[any]
// @Router /api/task/invoke/{id} [post]
func (t *TaskApi) Exec(c *gin.Context) {
	id := admin.ParseParamId(c, "id")
	err := t.TaskSrv.Execute(c.Request.Context(), id)
	if err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	}
	resp.Ok(c)
}

// @title 函数任务列表
// @Summary 函数任务列表
// @Description 服务端可运行的函数任务
// @Tags Task
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} resp.Response[vo.KV]
// @Router /api/task/jobs [get]
func (t *TaskApi) FuncJobs(c *gin.Context) {
	options := t.TaskSrv.FindFuncJobsOptions(c.Request.Context())
	resp.OKWithData(c, options)
}

// @title 启用/禁用定时任务
// @Summary 启用/禁用定时任务
// @Description 启用/禁用定时任务
// @Tags Task
// @Accept json
// @Produce json
// @Param id path integer true "id"
// @Param status path integer true "status"
// @Success 200 {object} resp.Response[any]
// @Router /api/task/status [patch]
// @Security Bearer
func (t *TaskApi) Status(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	statusStr := c.Param("status")
	status, err := utils.ToUint[uint8](statusStr)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	err = t.TaskSrv.ChangeStatus(c.Request.Context(), id, status)
	if err != nil {
		resp.FailWithError(c, err)
	}
	resp.Ok(c)
}
