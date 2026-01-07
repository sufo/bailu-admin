/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 系统配置
 */

package system

import (
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/service/sys"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils/page"
	"strconv"
)

type SysConfigApi struct {
	ConfigSrv *sys.SysConfigService
}

func NewSysConfigApi(configSrv *sys.SysConfigService) *SysConfigApi {
	return &SysConfigApi{configSrv}
}

// @title 系统参数配置
// @Summary 系统参数配置表
// @Description 按条件查询系统参数配置列表
// @Tags SysConfig
// @Accept json
// @Produce json
// @Security Bearer
// @Param query query dto.ConfigParams false "查询条件"
// //@Success 200 {object} resp.Response{data=resp.PageResult[entity.SysConfig]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.SysConfig]]
// @Router /api/config [get]
func (s *SysConfigApi) Index(c *gin.Context) {
	var params dto.ConfigParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	}
	page.StartPage(c)
	result, err := s.ConfigSrv.List(c.Request.Context(), params)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 创建系统参数
// @Summary 创建系统参数
// @Description 创建系统参数
// @Tags SysConfig
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.SysConfig false "Sysconfig"
// @Success 200 {object} resp.Response[any]
// @Router /api/config [post]
func (s *SysConfigApi) Create(c *gin.Context) {
	var config entity.SysConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		ctx := c.Request.Context()
		//检查字段唯一性
		if !s.ConfigSrv.CheckUnique(ctx, "name=?", config.Name) {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", config.Name))
			return
		}
		if !s.ConfigSrv.CheckUnique(ctx, "key=?", config.Key) {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", config.Key))
			return
		}

		err := s.ConfigSrv.Create(c.Request.Context(), &config)
		if err != nil {
			resp.InternalServerError(c)
			c.Abort()
			return
		}
		resp.Ok(c)
	}
}

// @title 修改系统参数
// @Summary 修改系统参数
// @Description 修改系统参数
// @Tags SysConfig
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.SysConfig false "sysConfig"
// //@Success 200 {object} resp.Response[entity.SysConfig]
// @Router /api/config [put]
func (s *SysConfigApi) Edit(c *gin.Context) {
	var config entity.SysConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	} else {
		//检查字段唯一性
		if !s.ConfigSrv.CheckUnique(c.Request.Context(), "id!=? and name=?", config.ID, config.Name) {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", config.Name))
			return
		}
		if !s.ConfigSrv.CheckUnique(c.Request.Context(), "id!=? and key=?", config.ID, config.Key) {
			resp.FailWithMsg(c, i18n.DefTr("admin.existed", config.Key))
			return
		}
		task, err := s.ConfigSrv.Update(c.Request.Context(), &config)
		if err != nil {
			resp.InternalServerError(c)
			return
		}
		resp.OKWithData(c, task)
	}
}

// @title 批量删除系统参数
// @Summary 批量删除系统参数
// @Description 批量删除系统参数
// @Tags SysConfig
// @Accept json
// @Produce json
// @Security Bearer
// @Param ids path string true "ids"
// @Success 200 {object} resp.Response[any]
// @Router /api/config/{ids} [delete]
func (s *SysConfigApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	if err := s.ConfigSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 启用/禁用系统配置
// @Summary 启用/禁用系统配置
// @Description 启用/禁用系统配置
// @Tags SysConfig
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body object{id=integer,status=integer} true "id和status"
// @Success 200 {object} resp.Response[any]
// @Router /api/config/status [patch]
func (s *SysConfigApi) Status(c *gin.Context) {

	var params = make(map[string]string)
	if err := c.ShouldBindJSON(&params); err != nil {
		panic(respErr.BadRequestError)
	}
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	status, err := strconv.ParseInt(params["status"], 10, 32)
	if err != nil {
		panic(respErr.BadRequestError)
	}
	err = s.ConfigSrv.Status(c.Request.Context(), id, int(status))
	if err != nil {
		log.L.Error(err)
		resp.InternalServerError(c)
		return
	}
	resp.Ok(c)
}
