package api

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"os"
)

func Chat(w http.ResponseWriter, req *http.Request) {
	token := os.Getenv("GPT_TOKEN")
	gptUrl := os.Getenv("GPT_URL")
	cfg := openai.DefaultConfig(token)
	if gptUrl != "" {
		cfg.BaseURL = gptUrl
	}
	client := openai.NewClientWithConfig(cfg)
	msg := req.URL.Query().Get("msg")
	if msg == "" {
		msg = "介绍你自己"
	}
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
		fmt.Fprintf(w, "gptclient err:%E", err)
		return
	}
	fmt.Fprintf(w, resp.Choices[0].Message.Content)
}
