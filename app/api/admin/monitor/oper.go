/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 操作日志
 */

package monitor

import (
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/service/sys"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils/page"
	"github.com/gin-gonic/gin"
)

type OperationApi struct {
	OperSrv *sys.OperationService
}

func NewOperApi(OperSrv *sys.OperationService) *OperationApi {
	return &OperationApi{OperSrv}
}

// @title 操作日志列表
// @Summary 操作日志列表接口
// @Description 按条件查询操作日志列表
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param query query dto.OperParams false "查询条件"
// //////@Success 200 {object} resp.Response{data=resp.PageResult[entity.OperationRecord]}
// @Success 200 {object} resp.Response[resp.PageResult[entity.OperationRecord]]
// @Router /api/oper [get]
// @Security Bearer
func (p *OperationApi) Index(c *gin.Context) {
	var queryParams dto.OperParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		panic(respErr.BadRequestError)
	}
	page.StartPage(c)
	result, err := p.OperSrv.List(c.Request.Context(), queryParams)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 批量删除操作日志
// @Summary 批量删除操作日志接口
// @Description 批量删除操作日志
// @Tags OperationLog
// @Accept json
// @Produce json
// @Param ids path string true "ids"
// @Success 200 {object} resp.Response[any]
// @Router /api/oper/{ids} [delete]
// @Security Bearer
func (p *OperationApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	if err := p.OperSrv.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}
