/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

type RateLimiter struct {
	Enable  bool  `mapstructure:"enable" json:"enable" yaml:"enable"`    // 是否开启请求频率限制
	Count   int64 `mapstructure:"count" json:"count" yaml:"count"`       // 每分钟每个用户允许的最大请求数量
	RedisDB int   `mapstructure:"redisDB" json:"redisDB" yaml:"redisDB"` //
}
