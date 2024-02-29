package chat

import (
	"context"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/sashabaranov/go-openai"
	"os"
	"time"
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
	token    string
	url      string
	cacheMsg map[string]string
}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	if s.cacheMsg[userID] != "" {
		r := s.cacheMsg[userID]
		delete(s.cacheMsg, userID)
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
	case <-time.After(4 * time.Second):
		go func() {
			s.cacheMsg[userID] = <-resChan
		}()
		return "响应内容过长，重新发送任意回复获取答复"
	}
}

func GetChatBot() BaseChat {
	err := config.CheckConfig()
	if err != nil {
		return &ErrorChat{
			errMsg: err.Error(),
		}
	}
	useType := config.UseType
	switch useType {
	case config.GPT:
		url := os.Getenv("GPT_URL")
		if url == "" {
			url = "https://api.openai.com/v1/"
		}
		return &SimpleGptChat{
			token:    os.Getenv("GPT_TOKEN"),
			url:      url,
			cacheMsg: make(map[string]string),
		}
	case config.ECHO:
		return &Echo{}
	}
	return &Echo{}
}
