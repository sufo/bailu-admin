/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 服务器信息
 */

package monitor

import (
	"bailu/app/config"
	"bailu/app/domain/resp"
	"bailu/app/domain/vo"
	"bailu/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServerInfo struct{}

func NewServerInfo() *ServerInfo {
	return &ServerInfo{}
}

func (s *ServerInfo) GetSeverConfig(c *gin.Context) config.Server {
	return config.Conf.Server
}

// @title 服务器信息
// @Summary 获取服务器信息接口
// @Description 获取服务器信息
// @Tags Server
// @Accept json
// @Produce json
// ////@Success 200 {object} resp.Response{data=vo.Server}
// @Success 200 {object} resp.Response[vo.Server]
// @Security Bearer
// @Router /api/server [get]
func (s *ServerInfo) GetServerInfo(c *gin.Context) {
	if err := vo.ServerInfo.CopyTo(c); err != nil {
		log.L.Error("func GetServerInfo Failed", zap.String("err", err.Error()))
		resp.InternalServerError(c)
		return
	}
	resp.OKWithData(c, vo.ServerInfo)
}
