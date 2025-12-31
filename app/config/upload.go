/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package config

type Qiniu struct {
	Region        string `mapstructure:"region" json:"region" yaml:"region"`                            // 存储区域
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`                            // 空间名称
	ImgPath       string `mapstructure:"img-path" json:"img-path" yaml:"img-path"`                      // CDN加速域名
	UseHTTPS      bool   `mapstructure:"use-https" json:"use-https" yaml:"use-https"`                   // 是否使用https
	AccessKey     string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`                // 秘钥AK
	SecretKey     string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`                // 秘钥SK
	UseCdnDomains bool   `mapstructure:"use-cdn-domains" json:"use-cdn-domains" yaml:"use-cdn-domains"` // 上传是否使用CDN上传加速
}

type Upload struct {
	Model              string `json:"model" yaml:"model"` //上传模式 local、cloud、other
	Type               string `json:"type" yaml:"type"`   //上传类型
	MaxMultipartMemory int    `mapstructure:"max-multipart-memory" json:"maxMultipartMemory" yaml:"max-multipart-memory"`
	MimeType           string `mapstructure:"mime-file" json:"mimeType" yaml:"mime-file"` //允许的上传文件类型
	Qiniu              Qiniu  `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`            //
}
