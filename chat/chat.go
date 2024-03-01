package chat

import (
	"context"
	"os"
	"time"

	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/sashabaranov/go-openai"
)

type BaseChat interface {
	Chat(userID string, msg string) string
}

type ErrorChat struct {
	errMsg string
}

func (e *ErrorChat) Chat(userID string, msg string) string {
	return e.errMsg
}

type Echo struct{}

func (e *Echo) Chat(userID string, msg string) string {
	return msg
}

type SimpleGptChat struct {
	token string
	url   string
}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	if _, ok := config.Cache.Load(userID); ok {
		rAny, _ := config.Cache.Load(userID)
		r := rAny.(string)
		config.Cache.Delete(userID)
		return r
	}
	cfg := openai.DefaultConfig(s.token)
	cfg.BaseURL = s.url
	client := openai.NewClientWithConfig(cfg)
	resChan := make(chan string)
	go func() {
		resp, err := client.CreateChatCompletion(context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: msg,
					},
				},
			})
		if err != nil {
			resChan <- err.Error()
			return
		}
		resChan <- resp.Choices[0].Message.Content
	}()

	select {
	case res := <-resChan:
		return res
	case <-time.After(5 * time.Second):
		config.Cache.Store(userID, <-resChan)
		return ""
	}
}

func GetChatBot() BaseChat {
	botType := config.GetBotType()
	switch botType {
	case config.Bot_Type_Gpt:
		err := config.CheckGptConfig()
		if err != nil {
			return &ErrorChat{
				errMsg: err.Error(),
			}
		}
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
