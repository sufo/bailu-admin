/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
}

type BuntDb struct {
	FilePath string `mapstructure:"file-path" json:"filePath" yaml:"file-path"` //":memory:":表示buntdb不会将数据保存到磁盘
}

type Store struct {
	//存储类型
	StoreType   string `mapstructure:"store-type" json:"storeType" yaml:"store-type"`
	Redis       Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	BuntDb      BuntDb `mapstructure:"bunt-db" json:"buntDb" yaml:"bunt-db"`
	EnableCache bool   `mapstructure:"enable-cache" json:"enableCache" yaml:"enable-cache"`
}
