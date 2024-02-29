package config

import (
	"errors"
	"os"
)

const (
	GPT  = "gpt"
	ECHO = "echo"
)

var UseType = ""

func CheckConfig() error {
	gptToken := os.Getenv("GPT_TOKEN")
	token := os.Getenv("TOKEN")
	_type := os.Getenv("TYPE")
	if _type != GPT && _type != ECHO {
		return errors.New("请配置bot类型")
	}
	UseType = _type
	if token == "" {
		return errors.New("请配置微信TOKEN")
	}
	if gptToken == "" && _type == "gpt" {
		return errors.New("请配置ChatGPTToken")
	}
	return nil
}
