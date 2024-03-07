package model

import (
	"github.com/pwh-pwh/aiwechat-vercel/chat"
	"github.com/sashabaranov/go-openai"
)

type ChatMsg interface {
	openai.ChatCompletionMessage | chat.QwenMessage | chat.SparkMessage
}
