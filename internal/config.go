package internal

import (
	"github.com/spf13/viper"
)

type ConfigStruct struct {
	// wire 配置
	Wire ConfigWireStruct
}

type ConfigWireStruct struct {
	// 排除的目录
	Exclude []string
}

var Config = &ConfigStruct{}

func init() {
	var configViperConfig = viper.New()
	configViperConfig.AddConfigPath("./")
	configViperConfig.SetConfigName(".tools")
	configViperConfig.SetConfigType("yaml")
	//读取配置文件内容
	if err := configViperConfig.ReadInConfig(); err != nil {
		//fmt.Println("未指定.tools.yaml，使用默认配置", err)
		return
	}
	if err := configViperConfig.Unmarshal(&Config); err != nil {
		panic(err)
	}
}
