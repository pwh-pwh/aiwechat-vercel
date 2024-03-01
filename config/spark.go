package config

import (
	"errors"
	"os"
)

const (
	Spark_Host_Url_Key   = "sparkUrl"
	Spark_App_Id_Key     = "sparkAppId"
	Spark_App_Secret_Key = "sparkAppSecret"
	Spark_ApiKey_Key     = "sparkApiKey"
	Spark_Domain_Version = "generalv3.5"
)

type SparkConfig struct {
	HostUrl   string
	AppId     string
	ApiSecret string
	ApiKey    string
}

func GetSparkConfig() (cfg SparkConfig, err error) {
	cfg = SparkConfig{
		HostUrl:   os.Getenv(Spark_Host_Url_Key),
		AppId:     os.Getenv(Spark_App_Id_Key),
		ApiSecret: os.Getenv(Spark_App_Secret_Key),
		ApiKey:    os.Getenv(Spark_ApiKey_Key),
	}

	if cfg.HostUrl == "" {
		err = errors.New("请配置sparkUrl")
		return
	}
	if cfg.AppId == "" {
		err = errors.New("请配置sparkAppId")
		return
	}
	if cfg.ApiSecret == "" {
		err = errors.New("请配置sparkAppSecret")
		return
	}
	if cfg.ApiKey == "" {
		err = errors.New("请配置sparkApiKey")
		return
	}

	return
}
