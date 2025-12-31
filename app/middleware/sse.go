/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"bailu/app/api/admin/monitor"
	"bailu/app/domain/entity"
	"bailu/global/consts"
	"bailu/pkg/log"
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
