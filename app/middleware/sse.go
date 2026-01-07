/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/app/api/admin/monitor"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/global/consts"
	"github.com/sufo/bailu-admin/pkg/log"
)

func SSEMiddleware(stream *monitor.Event) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment global variable ID
		req_user, _ := c.Get(consts.REQUEST_USER)
		u := req_user.(*entity.OnlineUserDto)
		// Initialize client
		client := monitor.Client{
			ID:      u.ID,
			Channel: make(chan monitor.Data),
		}

		// Send new connection to event to store
		stream.NewClients <- client

		defer func() {
			// Send closed connection to event server
			log.L.Infof("Closing connection : %d", client.ID)
			stream.ClosedClients <- client
		}()

		c.Set("client", client)
		c.Next()
	}
}
