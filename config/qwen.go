package config

import (
	"errors"
	"os"
)

const (
	Qwen_Host_Url_Key  = "qwenUrl"
	Qwen_ApiKey_Key    = "qwenApiKey"
	Qwen_Model_Version = "qwenModelVersion"
)

type QwenConfig struct {
	HostUrl      string
	ApiKey       string
	ModelVersion string
}

func GetQwenConfig() (cfg *QwenConfig, err error) {
	cfg = &QwenConfig{
		HostUrl:      os.Getenv(Qwen_Host_Url_Key),
		ApiKey:       os.Getenv(Qwen_ApiKey_Key),
		ModelVersion: os.Getenv(Qwen_Model_Version),
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
