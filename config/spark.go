package config

import (
	"errors"
	"os"
	"regexp"
)

const (
	Spark_Host_Url_Key   = "sparkUrl"
	Spark_App_Id_Key     = "sparkAppId"
	Spark_App_Secret_Key = "sparkAppSecret"
	Spark_ApiKey_Key     = "sparkApiKey"
)

type SparkConfig struct {
	HostUrl            string
	AppId              string
	ApiSecret          string
	ApiKey             string
	SparkDomainVersion string
}

func extractVersion(url string) string {
	// 使用正则表达式匹配版本号
	regex := regexp.MustCompile(`v(\d+)\.(\d+)`)
	matches := regex.FindStringSubmatch(url)
	if len(matches) != 3 {
		return ""
	}

	// 返回版本号
	return matches[1] + "." + matches[2]
}

func GetSparkConfig() (cfg *SparkConfig, err error) {
	var sparkUrl string = os.Getenv(Spark_Host_Url_Key)
	if sparkUrl == "" {
		sparkUrl = "wss://spark-api.xf-yun.com/v3.5/chat"
	}
	version := extractVersion(sparkUrl)
	var sparkDomainVersion = ""

	switch version {
	case "3.5":
		sparkDomainVersion = "generalv3.5"
	case "3.1":
		sparkDomainVersion = "generalv3"
	case "2.1":
		sparkDomainVersion = "generalv2"
	case "1.5":
		sparkDomainVersion = "general"
	default:
		sparkDomainVersion = ""
	}

	cfg = &SparkConfig{
		HostUrl:            sparkUrl,
		AppId:              os.Getenv(Spark_App_Id_Key),
		ApiSecret:          os.Getenv(Spark_App_Secret_Key),
		ApiKey:             os.Getenv(Spark_ApiKey_Key),
		SparkDomainVersion: sparkDomainVersion,
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
	if cfg.SparkDomainVersion == "" {
		err = errors.New("请配置sparkUrl")
		return
	}

	return
}
