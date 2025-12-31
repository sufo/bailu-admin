/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 服务监控
 */

package config

type Monitor struct {
	Enable    bool   `mapstructure:"enable" json:"enable" yaml:"enable"`            // 是否启用
	Addr      string `mapstructure:"addr" json:"addr" yaml:"addr"`                  // 服务器地址:端口
	ConfigDir string `mapstructure:"config-dir" json:"configDir" yaml:"config-dir"` // 密码
}
