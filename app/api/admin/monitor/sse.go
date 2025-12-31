/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc SSE
 */

package monitor

import (
	"bailu/app/config"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/pkg/jwt"
	"bailu/pkg/log"
	"bailu/pkg/mq"
	"bailu/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

const (
	MSG_TYPE = "message"
)

// Data to be broadcasted to a client.
type Data struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	To      uint64 `json:"receiver"`
}

// Uniquely defines an incoming client.
type Client struct {
	// Unique Client ID <-> Data.To
	ID uint64
	// Client channel
	Channel chan Data
}

// Keeps track of every SSE events.
type Event struct {
	// Data are pushed to this channel
	Message chan Data

	// New client connections
	NewClients chan Client

	// Closed client connections
	ClosedClients chan Client

	// Total client connections
	Clients map[uint64]chan Data
}

// Initializes Event and starts the event listener
func NewEvent() (event *Event) {
	event = &Event{
		Message:       make(chan Data),
		NewClients:    make(chan Client),
		ClosedClients: make(chan Client),
		Clients:       make(map[uint64]chan Data),
	}
	go event.listen()
	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.Clients[client.ID] = client.Channel
			log.L.Infof("Added client. %d registered clients", len(stream.Clients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.Clients, client.ID)
			close(client.Channel)
			log.L.Infof("Removed client. %d registered clients", len(stream.Clients))

		// Broadcast message to a specific client with client ID fetched from eventMsg.To
		case eventMsg := <-stream.Message:
			stream.Clients[eventMsg.To] <- eventMsg
		}
	}
}

// Mandatory Headers which should be set in the Response header for SSE to work.
func StreamHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

// 中间件
// 中间件必须有前置的AuthMiddleware中间件才能工作
// 这样就必须依赖前端EventSource通过header传递token
//	func (stream *Event) SSEMiddleware() gin.HandlerFunc {
//		return func(c *gin.Context) {
//			// Increment global variable ID
//			req_user, _ := c.Get(consts.REQUEST_USER)
//			u := req_user.(*entity.OnlineUserDto)
//			// Initialize client
//			client := Client{
//				ID:      u.ID,
//				Channel: make(chan Data),
//			}
//
//			// Send new connection to event to store
//			stream.NewClients <- client
//
//			defer func() {
//				// Send closed connection to event server
//				log.L.Infof("Closing connection : %d", client.ID)
//				stream.ClosedClients <- client
//			}()
//
//			c.Set("client", client)
//			c.Next()
//		}
//	}

// 不需要前置AuthMiddleware中间件
// 因为H5 EventSource无法传递header
func (stream *Event) SSEMiddleware(provider *jwt.JwtProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Increment global variable ID
		var token = c.Query("token")
		if token == "" {
			resp.Unauthorized(c)
			c.Abort()
			return
		}

		//解析key
		userKey, err := jwt.ParseUserKey(token)
		if err != nil {
			resp.Unauthorized(c)
			c.Abort()
			return
		}
		res, err := provider.Store.Get(config.Conf.JWT.OnlineKey + userKey)
		if err != nil {
			resp.Unauthorized(c)
			c.Abort()
			return
		}
		var u = entity.OnlineUserDto{}
		err2 := json.Unmarshal([]byte(res), &u)
		if err2 != nil {
			resp.Unauthorized(c)
			c.Abort()
			return
		}

		// Initialize client
		client := Client{
			ID:      u.ID,
			Channel: make(chan Data),
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

// @title SSE
// @Summary SSE
// @Description SSE消息通知
// @Tags Server Send Event
// @Accept json
// @Produce octet-stream
// @Security Bearer
// @Success 200 {string} binary
// @Router /api/sse [get]
func (e *Event) Stream(c *gin.Context) {
	v, ok := c.Get("client")
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}
	client, ok := v.(Client)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	// This goroutine will send the above data to Message channel
	// Which will pass through listen(), where it will get sent to the specified client (To)
	go e.handleMsg(c.Request.Context(), client)

	c.Stream(func(w io.Writer) bool {
		// Stream data to client
		for {
			select {
			// Send msg to the client
			case msg, ok := <-client.Channel:
				if !ok {
					return false
				}
				c.SSEvent("message", msg)
				return true
			// Client exit
			case <-c.Request.Context().Done():
				return false
			}
		}
	})
}

func (e *Event) handleMsg(c context.Context, client Client) {
	msgChan, err := mq.Consumer.GetMsgChan(c)
	if err != nil {
		log.L.Infof("get massage chan err: +%v", err)
		return
	}
	//测试
	data := Data{
		Message: "New Client in town",
		To:      1, // To send this data to a specified client, you can change this to the specific client ID
	}
	e.Message <- data

	for {
		select {
		case msgs, ok := <-msgChan:
			if !ok {
				return
			}
			if e.Clients[client.ID] == nil {
				log.L.Infof("Receiver - %d doesn't exist or disconnected.", client.ID)
			} else {
				sended := make([]*mq.StreamMsg, 0)
				for _, msg := range msgs {
					receiverId, err := utils.ToUint[uint64](msg.Values["receiveId"])
					if err != nil {
						continue
					}
					//当前连接的客户端=消息接收方
					if client.ID == receiverId {
						data := Data{
							Message: utils.Map2String(msg.Values),
							Type:    MSG_TYPE,
							To:      client.ID, // To send this data to a specified client, you can change this to the specific client ID
						}
						e.Message <- data
						sended = append(sended, msg)
					}

				}

				//确认已处理消息
				if len(sended) > 0 {
					err := mq.Consumer.AckMsg(c, sended)
					if err != nil {
						log.L.Infof("consumer AckMsg err: +%v", err)
					}
				}
			}
		}

	}
}
