// Package config 配置包
package config

import (
	"github.com/spf13/viper"
	"os"
)

var globalConf *GlobalConf

// GlobalConf 全局配置文件
type GlobalConf struct {
	// SystemConf 系统配置
	SystemConf SystemConf `json:"systemConf" yaml:"systemConf"`
	WeConfig   WeConfig   `json:"weConfig" yaml:"weConfig"`
}

// SystemConf 系统配置
type SystemConf struct {
	Port        string     `json:"port" yaml:"port"`
	MsgMode     string     `json:"msgMode" yaml:"msgMode"`
	CallBackUrl string     `json:"callBackUrl" yaml:"callBackUrl"`
	LogConf     LoggerConf `json:"logConf" yaml:"logConf"`
}

// DBConf 配置
type WeConfig struct {
	Corpid           string `json:"corpid" yaml:"corpid"`
	CorpSecret       string `json:"corpSecret" yaml:"corpSecret"`
	AgentId          int    `json:"agentId" yaml:"agentId"`
	WeApiRCallToken  string `json:"weApiRCallToken" yaml:"weApiRCallToken"`
	WeApiEncodingKey string `json:"weApiEncodingKey" yaml:"weApiEncodingKey"`
	WeChatApiAddr    string `json:"weChatApiAddr" yaml:"weChatApiAddr"`
}

type LoggerConf struct {
	LogLevel      string `json:"logLevel,omitempty" yaml:"logLevel"`
	LogOutPutMode string `json:"logOutPutMode" yaml:"logOutPutMode"`
	LogOutPutPath string `json:"logOutPutPath,omitempty" yaml:"logOutPutPath"`
	LogFormatter  string `json:"logFormatter,omitempty" yaml:"logFormatter"`
	LogTimeFormat string `json:"logTimeFormat" yaml:"logTimeFormat"`
}

func GetSystemConf() *SystemConf {
	if globalConf == nil {
		return nil
	}
	return &globalConf.SystemConf
}

func GetWechatConf() *WeConfig {
	if globalConf == nil {
		return nil
	}
	return &globalConf.WeConfig
}

// LoadConf  加载配置文件
func init() {
	v := viper.New()
	v.SetConfigType("yaml")
	WorkPath, _ := os.Getwd()
	v.AddConfigPath(WorkPath + "/config/")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&globalConf); err != nil {
		panic(err)
	}
}
