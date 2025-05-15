package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	AppID         string `mapstructure:"APP_ID" yaml:"APP_ID"`
	AppSecret     string `mapstructure:"APP_SECRET" yaml:"APP_SECRET"`
	WebHookURL    string `mapstructure:"WEBHOOK_URL" yaml:"WEBHOOK_URL"`
	WebHookSecret string `mapstructure:"WEBHOOK_SECRET" yaml:"WEBHOOK_SECRET"`
	AiKey         string `mapstructure:"AI_KEY" yaml:"AI_KEY"`
	AiBaseurl     string `mapstructure:"AI_BASEURL" yaml:"AI_BASEURL"`
	AiModel       string `mapstructure:"AI_MODEL" yaml:"AI_MODEL"`
}

// GetConf 获取配置实例（默认不生成yaml）
func GetConf() *Config {
	return GetConfWithGenerate(false)
}

// GetConfWithGenerate 可选是否生成 config.yaml 文件
func GetConfWithGenerate(generateYaml bool) *Config {
	once.Do(func() {
		initConf(generateYaml)
	})
	return conf
}

// 初始化配置
func initConf(generateYaml bool) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 如果找不到配置文件，则回退到环境变量
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("未找到 config.yaml，尝试读取环境变量")
		viper.AutomaticEnv()
		viper.SetEnvPrefix("FeiShuBot")
	}

	// 反序列化到结构体
	conf = &Config{}
	if err := viper.Unmarshal(conf); err != nil {
		panic(err)
	}

	// 可选写入 config.yaml
	if generateYaml {
		v := viper.New()
		v.Set("APP_ID", conf.AppID)
		v.Set("APP_SECRET", conf.AppSecret)
		v.Set("WEBHOOK_URL", conf.WebHookURL)
		v.Set("WEBHOOK_SECRET", conf.WebHookSecret)
		v.Set("AI_KEY", conf.AiKey)
		v.Set("AI_BASEURL", conf.AiBaseurl)
		v.Set("AI_MODEL", conf.AiModel)

		err := v.WriteConfigAs("config.yaml")
		if err != nil {
			fmt.Println("写入 config.yaml 失败:", err)
		} else {
			fmt.Println("已生成 config.yaml")
		}
	}
}
