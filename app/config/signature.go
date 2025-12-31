package config

import "time"

/**
* 接口签名验证，防参数篡改
 */
type Signature struct {
	Key    string        `json:"key" yaml:"key"`
	Secret string        `json:"secret" yaml:"secret"`
	TTL    time.Duration `json:"ttl" yaml:"ttl"`
	Enable bool          `json:"enable" yaml:"enable"`
}
