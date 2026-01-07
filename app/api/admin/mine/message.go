/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 个人所有种类消息处理
 */

package mine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/service/message"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
)

var MessageSet = wire.NewSet(wire.Struct(new(Message), "*"))

type Message struct {
	MsgSrv *message.MessageService
}

// @title 我的消息列表
// @Summary 我的消息列表
// @Description 按条件查询消息
// @Tags Mine
// @Accept json
// @Produce json
// @Param query query dto.MessageParams false "查询条件"
// @Success 200 {object} resp.Response[resp.PageResult[any]]
// @Security Bearer
// @Router /api/mine/message/unread [get]
func (m *Message) Unread(c *gin.Context) {
	var params dto.MessageParams
	if err := c.ShouldBindQuery(&params); err != nil {
		panic(respErr.BadRequestError)
	}
	result, err := m.MsgSrv.UnreadList(c.Request.Context(), params)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	resp.OKWithData(c, result)
}

// @title 我的消息批量删除
// @Summary 我的消息批量删除
// @Description 我的消息批量删除
// @Tags Mine
// @Accept json
// @Produce json
// @Param msgType path string true "消息类型（notice,event,chat）"
// @Param ids path string true "id集合"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/mine/{msgType}/{ids} [delete]
func (m *Message) Destroy(c *gin.Context) {
	ids := admin.ParseParamIDs(c, "ids")
	msgType := c.Param("msgType")
	if msgType == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("path Param must be required")))
	}
	err := m.MsgSrv.Delete(c.Request.Context(), ids, msgType)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}
	resp.Ok(c)
}

// @title 我的消息清空
// @Summary 我的消息清空
// @Description 我的消息清空
// @Tags Mine
// @Accept json
// @Produce json
// @Param msgType path string true "消息类型（notice,event,chat）"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/mine/{msgType}/clear [delete]
func (m *Message) Clear(c *gin.Context) {
	msgType := c.Param("msgType")
	if msgType == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("path Param must be required")))
	}
	err := m.MsgSrv.Clear(c.Request.Context(), msgType)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}
	resp.Ok(c)
}

// @title 我的未读消息数量
// @Summary 未读消息数量
// @Description 按条件查询未读消息数量
// @Tags Mine
// @Accept json
// @Produce json
// @Param msgType path string true "查询类型：notice、event、chat、msg查所有"
// @Success 200 {object} resp.Response[any]{data=object{notice=int, event=int,chat=int,total=int}}
// @Security Bearer
// @Router /api/mine/{msgType}/unread_count [get]
func (m *Message) UnreadCount(c *gin.Context) {
	msgType := c.Param("msgType")
	if msgType == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("msgType Param must be required")))
	}
	result, err := m.MsgSrv.UnreadCount(c, msgType)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}
	resp.OKWithData(c, result)
}

// @title 所有消息设为已读
// @Summary 所有未读消息消息设为已读
// @Description 按msgType设置所有消息已读
// @Tags Mine
// @Accept json
// @Produce json
// @Param msgType path string true "设置已读：notice、event、chat、msg设置所有"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/mine/{msgType}/read_all [put]
func (m *Message) ReadAll(c *gin.Context) {
	msgType := c.Param("msgType")
	if msgType == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("msgType Param must be required")))
	}
	err := m.MsgSrv.ReadAll(c, msgType)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}
	resp.Ok(c)
}

// @title 消息设为已读
// @Summary 消息设为已读
// @Description 消息设为已读
// @Tags Mine
// @Accept json
// @Produce json
// @Param msgType path string true "msgType"
// @Param id path string true "id"
// @Success 200 {object} resp.Response[any]
// @Security Bearer
// @Router /api/mine/{msgType}/read/{id} [put]
func (m *Message) Read(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("id Param must be required")))
	}
	msgType := c.Param("msgType")
	if msgType == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("msgType Param must be required")))
	}
	err := m.MsgSrv.Read(c, msgType, id)
	if err != nil {
		panic(respErr.InternalServerErrorWithError(err))
	}
	resp.Ok(c)
}
