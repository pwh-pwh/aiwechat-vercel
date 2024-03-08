package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/pwh-pwh/aiwechat-vercel/db"
)

const (
	GPT  = "gpt"
	ECHO = "echo"

	Gpt_Token = "GPT_TOKEN"

	Bot_Type_Key   = "botType"
	Bot_Type_Echo  = "echo"
	Bot_Type_Gpt   = "gpt"
	Bot_Type_Spark = "spark"
	Bot_Type_Qwen  = "qwen"
)

var (
	Cache sync.Map

	Bot_Type = os.Getenv(Bot_Type_Key)
)

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
	}
	return
}

func CheckAllBotConfig() (botType string, checkRes map[string]bool) {
	botType = GetBotType()
	checkRes = map[string]bool{
		Bot_Type_Echo:  true,
		Bot_Type_Gpt:   true,
		Bot_Type_Spark: true,
		Bot_Type_Qwen:  true,
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

	return
}

func CheckGptConfig() error {
	gptToken := os.Getenv(Gpt_Token)
	token := os.Getenv(Wx_Token)
	botType := GetBotType()
	if token == "" {
		return errors.New("请配置微信TOKEN")
	}
	if gptToken == "" && botType == Bot_Type_Gpt {
		return errors.New("请配置ChatGPTToken")
	}
	return nil
}

var (
	Support_Bots = []string{Bot_Type_Gpt, Bot_Type_Spark, Bot_Type_Qwen}
)

func GetBotType() string {
	botType := Bot_Type
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
