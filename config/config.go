package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"

	"github.com/pwh-pwh/aiwechat-vercel/db"
)

const (
	GPT  = "gpt"
	ECHO = "echo"

	Bot_Type_Key    = "botType"
	Bot_Type_Echo   = "echo"
	Bot_Type_Gpt    = "gpt"
	Bot_Type_Spark  = "spark"
	Bot_Type_Qwen   = "qwen"
	Bot_Type_Gemini = "gemini"
)

var (
	Cache sync.Map

	Support_Bots = []string{Bot_Type_Gpt, Bot_Type_Spark, Bot_Type_Qwen, Bot_Type_Gemini}
)

func IsSupportPrompt(botType string) bool {
	return botType == Bot_Type_Gpt || botType == Bot_Type_Qwen || botType == Bot_Type_Spark
}

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
	gptToken := GetGptToken()
	token := GetWxToken()
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
	key := GetGeminiKey()
	if key == "" {
		return errors.New("请配置geminiKey")
	}
	return nil
}

func GetBotType() string {
	botType := os.Getenv(Bot_Type_Key)
	if slices.Contains(Support_Bots, botType) {
		return botType
	} else {
		return Bot_Type_Echo
	}
}

func GetUserBotType(userId string) (bot string) {
	bot, err := db.GetValue(fmt.Sprintf("%v:%v", Bot_Type_Key, userId))
	if err != nil {
		bot = GetBotType()
	}
	if !slices.Contains(Support_Bots, bot) {
		bot = GetBotType()
	}
	return
}

func GetMaxTokens() int {
	// 不设置或者设置不合法，均返回0，模型将使用默认值或者不设置
	maxTokensStr := os.Getenv("maxOutput")
	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		return 0
	}
	return maxTokens
}

func GetDefaultSystemPrompt() string {
	return os.Getenv("defaultSystemPrompt")
}
