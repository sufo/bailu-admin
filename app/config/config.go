/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

import (
	"bailu/global"
	"bailu/global/consts"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

var Conf = new(Config)

type Config struct {
	Version string `json:"version" yaml:"version"`
	Server  Server `json:"server" yaml:"server"`
	Local   Local  `json:"Local" yaml:"Local"` //本地路径
	//DBType        string      `json:"DBType" yaml:"db-type"`               // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	//OssType       string      `json:"ossType" yaml:"oss-type"`             // Oss类型
	//UseMultipoint string      `json:"useMultipoint" yaml:"use-multipoint"` // 多点登录拦截
	Zap Zap `mapstructure:"zap" json:"zap" yaml:"zap"`
	//Redis       Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
	DataSource  DataSource  `mapstructure:"datasource" json:"datasource" yaml:"datasource"`
	Monitor     Monitor     `mapstructure:"Monitor" json:"monitor" yaml:"monitor"`
	CORS        CORS        `mapstructure:"cors" json:"cors" yaml:"cors"`
	GZIP        GZIP        `json:"GZIP" yaml:"gzip"`
	Swagger     bool        `json:"swagger" yaml:"swagger"`
	JWT         JWT         `json:"jwt" yaml:"jwt"`
	Captcha     Captcha     `json:"captcha" yaml:"captcha"`
	Store       Store       `json:"store" yaml:"store"`
	Signature   Signature   `json:"signature" yaml:"signature"` //接口签名验证
	RateLimiter RateLimiter `mapstructure:"rate-limiter" json:"rateLimiter" yaml:"rate-limiter"`
	RSA         RSA         `json:"-" yaml:"rsa"`
	SMS         SMS         `json:"-" yaml:"sms"`
	Menu        Menu        `json:"-" yaml:"menu"`
	Upload      Upload      `json:"-" yaml:"upload"`
	OperLog     OperLog     `json:"-" yaml:"operLog" mapstructure:"oper-log"`
	Casbin      Casbin      `json:"-" yaml:"casbin"`
}

type Local struct {
	Path string `json:"path" yaml:"path"`
	Dir  string `json:"dir" yaml:"dir"`
}
type RSA struct {
	PrivateKey string `mapstructure:"private-key" yaml:"private-key"`
}

type GZIP struct {
	Enable bool `yaml:"enable"`
	//排除的文件扩展名
	ExcludedExt []string `yaml:"excluded-ext" mapstructure:"excluded-ext"`
	//排除的请求路径
	ExcludedPaths []string `yaml:"excluded-paths" mapstructure:"excluded-paths"`
}

type Captcha struct {
	Length      int    `mapstructure:"length" yaml:"length"`
	Width       int    `mapstructure:"width" yaml:"width"`
	Height      int    `mapstructure:"height" yaml:"height"`
	CaptchaType string `mapstructure:"captcha-type" yaml:"captcha-type" mapstructure:"captcha-type"`
	Store       string `mapstructure:"store" yaml:"store"`
	RedisDB     int    `mapstructure:"redis-db" yaml:"redis-db" mapstructure:"redis-db"`
	Expire      int    `mapstructure:"expire" yaml:"expire"`
	Prefix      string `mapstructure:"prefix" yaml:"prefix"`
}

type SMS struct {
	AppKey       string `mapstructure:"app-key" yaml:"app-key"`
	AppSecret    string `mapstructure:"app-secret" yaml:"app-secret"`
	SignName     string `mapstructure:"sign-name" yaml:"sign-name"`
	TemplateCode string `mapstructure:"template-code" yaml:"template-code"`
	RegionId     string `mapstructure:"region-id" yaml:"region-id"`
	Expired      int    `yaml:"expired"` //过期时长 s
}

// 初始化导入菜单
type Menu struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
}

type OperLog struct {
	Enable   bool `yaml:"enable"`
	Interval int  `yaml:"interval"`
}

type Casbin struct {
	Enable           bool   `yaml:"enable"`
	Debug            bool   `yaml:"debug"`
	Model            string `yaml:"model"`
	AutoLoad         bool   `yaml:"auto-load"`                                            //自动加载
	AutoLoadInterval int    `yaml:"auto-load-interval" mapstructure:"auto-load-interval"` //加载间隔
}

func (c *Config) IsDebug() bool {
	return c.Server.Mode == consts.MODE_DEBUG
}

// 初始化默认值，将config/config.yml作为默认值
func (c *Config) Default() error {
	file, err := os.ReadFile(path.Join(global.Root, consts.ConfigDefault))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return yaml.Unmarshal(file, c)
}

// 处理默认值
//func (c *Config) UnmarshalJSON(data []byte) error {
//	type configAlias Config
//}
