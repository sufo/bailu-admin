/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package middleware

import (
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
	"strconv"
	"time"
)

// 根据ip限流或用户限流
// request rate limiter (per minute)
func RateLimiterMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.Conf.RateLimiter
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	rc := config.Conf.Store.Redis
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": rc.Addr,
		},
		Password: rc.Password,
		DB:       cfg.RedisDB,
	})
	limiter := redis_rate.NewLimiter(ring)
	limiter.Fallback = rate.NewLimiter(rate.Inf, 0)

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		limit := cfg.Count
		ip, err := utils.GetIP(c.Request)
		if err != nil {
			c.Next()
		}
		//rate, delay, allowed := limiter.AllowMinute(fmt.Sprintf("%d", c.Request.RemoteAddr), limit)
		rate, delay, allowed := limiter.AllowMinute(fmt.Sprintf("%d", ip), limit)
		if !allowed {
			h := c.Writer.Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(rate, 10))
			delaySec := int64(delay / time.Second)
			h.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
			resp.TooManyRequests(c)
			c.Abort()
			return
		}

		//根据用户来限流，但是前提时要先登录
		//if userID, isExist := c.Get(consts.REQUEST_USER); isExist {
		//	if userID != 0 {
		//		limit := cfg.Count
		//		rate, delay, allowed := limiter.AllowMinute(fmt.Sprintf("%d", userID), limit)
		//		if !allowed {
		//			h := c.Writer.Header()
		//			h.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		//			h.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
		//			delaySec := int64(delay / time.Second)
		//			h.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
		//			resp.TooManyRequests(c)
		//			return
		//		}
		//	}
		//}

		c.Next()

	}

}
