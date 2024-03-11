package config

import (
	"errors"
	"os"
)

const (
	Qwen_Host_Url_Key      = "qwenUrl"
	Qwen_ApiKey_Key        = "qwenApiKey"
	Qwen_Model_Version_Key = "qwenModelVersion"
	Qwen_Welcome_Reply_Key = "qwenWelcomeReply"
)

type QwenConfig struct {
	HostUrl      string
	ApiKey       string
	ModelVersion string
}

func GetQwenConfig() (cfg *QwenConfig, err error) {
	cfg = &QwenConfig{
		HostUrl:      GetQwenHostUrl(),
		ApiKey:       GetQwenApiKey(),
		ModelVersion: GetQwenModelVersion(),
	}

	if cfg.HostUrl == "" {
		err = errors.New("请配置qwenUrl")
		return
	}
	if cfg.ApiKey == "" {
		err = errors.New("请配置qwenApiKey")
		return
	}
	if cfg.ModelVersion == "" {
		err = errors.New("请配置qwenModelVersion")
		return
	}

	return
}

func GetQwenHostUrl() string {
	return os.Getenv(Qwen_Host_Url_Key)
}

func GetQwenApiKey() string {
	return os.Getenv(Qwen_ApiKey_Key)
}

func GetQwenModelVersion() string {
	return os.Getenv(Qwen_Model_Version_Key)
}

func GetQwenWelcomeReply() (r string) {
	r = os.Getenv(Qwen_Welcome_Reply_Key)
	if r == "" {
		r = "我是通义千问机器人，开始聊天吧！"
	}
	return
}
