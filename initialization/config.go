package initialization

import (
	"sync"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Initialized bool

	//app凭证
	FeishuAppId     string
	FeishuAppSecret string

	//事件订阅
	FeishuVerifiedToken string
	FeishuEncryptKey    string

	//推理请求地址
	LlamaUrl string
}

var (
	cfg    = pflag.StringP("config", "c", "./static/config.yaml", "apiserver config file path.")
	config *Config
	once   sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		config = LoadConfig(*cfg)
		config.Initialized = true
	})

	return config
}

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	return &Config{
		FeishuAppId:         viper.GetString("APP_ID"),
		FeishuAppSecret:     viper.GetString("APP_SECRET"),
		FeishuEncryptKey:    viper.GetString("APP_ENCRYPT_KEY"),
		FeishuVerifiedToken: viper.GetString("APP_VERIFICATION_TOKEN"),
		LlamaUrl:            viper.GetString("LLAMA_URL"),
	}
}
