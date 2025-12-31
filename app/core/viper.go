/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 配置管理
 */
package core

import (
	"bailu/app/config"
	"bailu/global/consts"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

var Viper = viper.New()

func InitViper(path ...string) {
	var confPath string
	if len(path) == 0 || path[0] == "" {
		if configEnv := os.Getenv(consts.ConfigEnv); configEnv == "" {
			confPath = consts.ConfigFile
			fmt.Printf("您正在使用config的默认值，config的路径为%v\n", consts.ConfigFile)
		} else {
			confPath = configEnv
			fmt.Printf("您正在使用BAILU_CONFIG环境变量，config的路径为%v\n", confPath)
		}
	} else {
		confPath = path[0]
		fmt.Printf("您正在使用命令行的-c参数传递的值，config的路径为%v\n", confPath)
		//fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", confPath)
	}
	//v := viper.New()
	Viper.SetConfigFile(confPath)
	Viper.SetConfigType(consts.ConfigType)

	if err := Viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal exception config file: #{err} \n", err))
	}

	Viper.WatchConfig()
	Viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file has changed: ", e.Name)
		if err := Viper.Unmarshal(&config.Conf); err != nil {
			fmt.Println(err)
		}
	})
	if err := Viper.Unmarshal(&config.Conf); err != nil {
		fmt.Println(err)
	}
}
