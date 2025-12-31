/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 跨域配置
 */

package config

type CORS struct {
	Enable           bool     `json:"enable" yaml:"enable"`
	AllowOrigins     []string `mapstructure:"allow-origins" json:"allowOrigins" yaml:"allow-origins"`
	AllowMethods     []string `mapstructure:"allow-methods" json:"allowMethods" yaml:"allow-methods"`
	AllowHeaders     []string `mapstructure:"allow-headers" json:"allowHeaders" yaml:"allow-headers"`
	AllowCredentials bool     `mapstructure:"allow-credentials" json:"allowCredentials" yaml:"allow-credentials"`
	MaxAge           int      `mapstructure:"max-age" json:"maxAge" yaml:"max-age"`
}
