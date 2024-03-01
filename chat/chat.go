package chat

import (
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"os"
	"time"
)

type BaseChat interface {
	Chat(userID string, msg string) string
}

// 加入超时控制
func WithTimeChat(userID, msg string, f func(userID, msg string) string) string {
	if _, ok := config.Cache.Load(userID); ok {
		rAny, _ := config.Cache.Load(userID)
		r := rAny.(string)
		config.Cache.Delete(userID)
		return r
	}
	resChan := make(chan string)
	go func() {
		resChan <- f(userID, msg)
	}()
	select {
	case res := <-resChan:
		return res
	case <-time.After(5 * time.Second):
		config.Cache.Store(userID, <-resChan)
		return ""
	}
}

type WithTimeoutChat interface {
	ChatWithTimeOut()
}

type SimpleWithTimeout struct {
	BaseChat
}

type ErrorChat struct {
	errMsg string
}

func (e *ErrorChat) Chat(userID string, msg string) string {
	return e.errMsg
}

func GetChatBot() BaseChat {
	botType, err := config.CheckBotConfig()
	if err != nil {
		return &ErrorChat{
			errMsg: err.Error(),
		}
	}

	switch botType {
	case config.Bot_Type_Gpt:
		url := os.Getenv("GPT_URL")
		if url == "" {
			url = "https://api.openai.com/v1/"
		}
		return &SimpleGptChat{
			token: os.Getenv("GPT_TOKEN"),
			url:   url,
		}
	case config.Bot_Type_Spark:
		return &SparkChat{}
	default:
		return &Echo{}
	}

}
