/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                           // 级别
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                        // 输出
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                        // 日志前缀
	Director      string `mapstructure:"director" json:"director" yaml:"director"`                  // 日志文件夹
	LinkName      string `mapstructure:"link-name" json:"linkName" yaml:"link-name"`                // 软链接名称
	ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"show-line-name"`           // 显示代码行
	EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`       //编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 栈名
	LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`  // 是否同时输出到控制台
	CleanDaysAgo  int    `mapstructure:"clean-days-ago" json:"CleanDaysAgo" yaml:"clean-days-ago"`  //清理多少天之前的日志
}
