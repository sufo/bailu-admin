/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package config

type DataSource struct {
	DbType string `mapstructure:"db-type" json:"dbType" yaml:"db-type"` //数据库类型
	//Debug             bool   `mapstructure:"debug" json:"debug" yaml:"debug"`                          //数据库调试是否打开
	TablePrefix       string `mapstructure:"table-prefix" json:"tablePrefix" yaml:"table-prefix"`      //表名前缀
	MaxIdleConns      int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns      int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	MaxLifetime       int    `mapstructure:"max-lifetime" json:"maxLifetime" yaml:"max-lifetime"`      // 设置连接可以重用的最长时间(单位：秒)
	LogZap            bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`                     // 是否通过zap写入日志文件
	EnableAutoMigrate bool   `json:"enableAutoMigrate" yaml:"enable-auto-migrate"`
	Mysql             Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
}

type Mysql struct {
	Driver   string `mapstructure:"driver" json:"driver" yaml:"driver"`       //数据库类型
	Host     string `mapstructure:"host" json:"host" yaml:"host"`             //ip
	Port     string `mapstructure:"port" json:"port" yaml:"port"`             //端口
	Username string `mapstructure:"username" json:"username" yaml:"username"` //数据库用户名
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 数据库密码
	DbName   string `mapstructure:"db-name" json:"dbName" yaml:"db-name"`     // 数据库名
	Params   string `mapstructure:"params" json:"params" yaml:"params"`       // 高级配置
	//LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`                  // 是否开启Gorm全局日志
	//LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`                     // 是否通过zap写入日志文件
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.DbName + "?" + m.Params
}
