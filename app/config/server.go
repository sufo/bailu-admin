/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

type Server struct {
	Mode          string `mapstructure:"mode" json:"mode" yaml:"mode"`                              //
	Port          int64  `mapstructure:"port" json:"port" yaml:"port"`                              //服务器端口
	Host          string `mapstructure:"host" json:"host" yaml:"host"`                              //服务器ip
	Name          string `mapstructure:"name" json:"name" yaml:"name"`                              //服务器名称
	ReadTimeout   int    `mapstructure:"read-timeout" json:"readTimeout" yaml:"read-timeout"`       //服务器名称
	WriterTimeout int    `mapstructure:"writer-timeout" json:"writerTimeout" yaml:"writer-timeout"` //服务器名称
	OssType       string `mapstructure:"oss-type" json:"ossType" yaml:"oss-type"`                   //存储类型
	//UseMultipoint bool   `json:"useMultipoint" yaml:"use-multipoint"`                               // 多点登录拦截
	UseMultiDevice bool   `json:"useMultiDevice" yaml:"use-multi-device"`
	UseRedis       bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"` // 使用redis
	Locale         string `json:"locale" yaml:"locale"`                                //服务端语言
	TimeZone       string `json:"timeZone" yaml:"time-zone"`
	CertFile       string `mapstructure:"cert-file" json:"certFile" yaml:"cert-file"`
	KeyFile        string `mapstructure:"key-file" json:"keyFile" yaml:"key-file"`
	//MaxMultipartMemory int    `mapstructure:"max-multipart-memory" json:"maxMultipartMemory" yaml:"max-multipart-memory"`
	//MimeType           string `mapstructure:"mime-file" json:"mimeType" yaml:"mime-file"` //允许的上传文件类型
}
