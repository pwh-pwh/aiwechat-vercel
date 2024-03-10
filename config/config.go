package config

import (
	"errors"
	"os"
	"slices"
	"sync"
)

const (
	GPT  = "gpt"
	ECHO = "echo"
)

var Cache sync.Map

func CheckBotConfig(botType string) (actualotType string, err error) {
	if botType == "" {
		botType = GetBotType()
	}
	actualotType = botType
	switch actualotType {
	case Bot_Type_Gpt:
		err = CheckGptConfig()
	case Bot_Type_Spark:
		_, err = GetSparkConfig()
	case Bot_Type_Qwen:
		_, err = GetQwenConfig()
	case Bot_Type_Gemini:
		err = CheckGeminiConfig()
	}
	return
}

func CheckAllBotConfig() (botType string, checkRes map[string]bool) {
	botType = GetBotType()
	checkRes = map[string]bool{
		Bot_Type_Echo:   true,
		Bot_Type_Gpt:    true,
		Bot_Type_Spark:  true,
		Bot_Type_Qwen:   true,
		Bot_Type_Gemini: true,
	}

	err := CheckGptConfig()
	if err != nil {
		checkRes[Bot_Type_Gpt] = false
	}
	_, err = GetSparkConfig()
	if err != nil {
		checkRes[Bot_Type_Spark] = false
	}
	_, err = GetQwenConfig()
	if err != nil {
		checkRes[Bot_Type_Qwen] = false
	}
	err = CheckGeminiConfig()
	if err != nil {
		checkRes[Bot_Type_Gemini] = false
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

func CheckGeminiConfig() error {
	key := os.Getenv("geminiKey")
	if key == "" {
		return errors.New("请配置geminiKey")
	}
	return nil
}

const (
	Bot_Type_Echo   = "echo"
	Bot_Type_Gpt    = "gpt"
	Bot_Type_Spark  = "spark"
	Bot_Type_Qwen   = "qwen"
	Bot_Type_Gemini = "gemini"
)

var (
	Support_Bots = []string{Bot_Type_Gpt, Bot_Type_Spark, Bot_Type_Qwen, Bot_Type_Gemini}
)

func GetBotType() string {
	botType := os.Getenv("botType")
	if slices.Contains(Support_Bots, botType) {
		return botType
	} else {
		return Bot_Type_Echo
	}
}
