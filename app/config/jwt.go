/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

type JWT struct {
	Enable bool `mapstructure:"enable" json:"enable" yaml:"enable"` // jwt签名`json:"enable" yaml:"enable"`
	//签名方式(支持：HS512/HS384/HS512)
	SigningMethod string `mapstructure:"signingMethod" json:"signingMethod" yaml:"signing-method"`
	//签名key
	SigningKey string `mapstructure:"signingKey" json:"signingKey" yaml:"signing-key"`
	//过期时间（单位秒）
	Expired int64 `mapstructure:"expired" json:"expired" yaml:"expired"`
	//存储(支持：file/redis)
	Store string `mapstructure:"store" json:"store" yaml:"store"`
	//文件路径
	//FilePath string `mapstructure:"filePath" json:"filePath" yaml:"file-path"` //"data/jwt_auth.db"
	//redis 数据库(如果存储方式是redis，则指定存储的数据库)
	//RedisDB uint `mapstructure:"redisDB" json:"redisDB" yaml:"redis-db"`
	//存储到 redis 数据库中的键名前缀
	RedisPrefix string `mapstructure:"redisPrefix" json:"redisPrefix" yaml:"redis-prefix"`
	//online-key 在线用户key
	OnlineKey string `mapstructure:"online-key" json:"onlineKey" yaml:"online-key"`
	//token 续期检查时间范围（默认30分钟，单位毫秒），在token即将过期的一段时间内用户操作了，则给用户的token续期
	Detect int64 `mapstructure:"detect" json:"-" ymal:"detect"`
	//续期时间范围，默认1小时，单位毫秒
	Renew int64 `mapstructure:"renew" json:"-" ymal:"renew"`
	//Authorization
	Header string `mapstructure:"header" json:"-" ymal:"header"`
	//令牌前缀Bearer
	TokenStartWith string `mapstructure:"tokenStartWith" json:"-" ymal:"token-start-with"`
}
