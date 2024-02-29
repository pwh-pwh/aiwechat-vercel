package config

import (
	"errors"
	"os"
)

func CheckConfig() error {
	gptToken := os.Getenv("GPT_TOKEN")
	token := os.Getenv("TOKEN")
	if token == "" {
		return errors.New("请配置微信TOKEN")
	}
	if gptToken == "" {
		return errors.New("请配置ChatGPTToken")
	}
	return nil
}
