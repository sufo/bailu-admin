/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

type Lumberjack struct {
	Filename   string `mapstructure:"filename" json:"filename" yaml:"filename"`         // 日志路径
	Maxsize    int    `mapstructure:"max-size" json:"maxSize" yaml:"max-size"`          // 单个文件最大尺寸，默认单位M
	MaxBackups int    `mapstructure:"max-backups" json:"maxBackups" yaml:"max-backups"` // 保留旧文件的最大个数
	MaxAge     int    `mapstructure:"max-age" json:"maxAge" yaml:"max-age"`             // 最大时间
	LocalTime  bool   `mapstructure:"local-time" json:"localTime" yaml:"local-time"`    // 使用本地时间
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`         // 是否压缩
}
