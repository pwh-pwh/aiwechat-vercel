package config

import (
	"errors"
	"os"
	"sync"
)

const (
	GPT  = "gpt"
	ECHO = "echo"
)

var Cache sync.Map

func CheckBotConfig() (botType string, err error) {
	botType = GetBotType()
	switch botType {
	case Bot_Type_Gpt:
		err = CheckGptConfig()
	case Bot_Type_Spark:
		_, err = GetSparkConfig()
	}
	return
}

func CheckGptConfig() error {
	gptToken := os.Getenv("GPT_TOKEN")
	token := os.Getenv("TOKEN")
	botType := GetBotType()
	if token == "" {
		return errors.New("请配置微信TOKEN")
	}
	if gptToken == "" && botType == Bot_Type_Gpt {
		return errors.New("请配置ChatGPTToken")
	}
	return nil
}

const (
	Bot_Type_Echo  = "echo"
	Bot_Type_Gpt   = "gpt"
	Bot_Type_Spark = "spark"
)

func GetBotType() string {
	botType := os.Getenv("botType")
	switch botType {
	case Bot_Type_Gpt:
		return Bot_Type_Gpt
	case Bot_Type_Spark:
		return Bot_Type_Spark
	default:
		return Bot_Type_Echo
	}
}
