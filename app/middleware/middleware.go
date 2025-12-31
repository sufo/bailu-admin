/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package middleware

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/global/consts"
	jwtProvider "bailu/pkg/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp.MethodNotAllowed(c)
	}
}

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp.NotFound(c)
	}
}

type SkipperFunc func(*gin.Context) bool

func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := JoinRouter(c.Request.Method, c.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

func EmptyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// Get jwt token from header (Authorization: Bearer xxx)
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserId(c *gin.Context) (uint64, error) {
	if user, exists := c.Get(consts.REQUEST_USER); !exists {
		return jwtProvider.ParseUserID(GetToken(c))
	} else {
		return user.(*entity.OnlineUserDto).ID, nil
	}
}

func DisabledLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("_trace_", nil)
	}
}
