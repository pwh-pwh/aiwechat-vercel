package chat

import (
	"context"

	"os"

	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/sashabaranov/go-openai"
)

type SimpleGptChat struct {
	token     string
	url       string
	maxTokens int
	BaseChat
}

func (s *SimpleGptChat) toDbMsg(msg openai.ChatCompletionMessage) db.Msg {
	return db.Msg{
		Role: msg.Role,
		Msg:  msg.Content,
	}
}

func (s *SimpleGptChat) toChatMsg(msg db.Msg) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    msg.Role,
		Content: msg.Msg,
	}
}


func (s *SimpleGptChat) getModel(userID string) string {
	if model, err := db.GetModel(userID, config.Bot_Type_Gpt); err == nil && model != "" {
		return model
	} else if model = os.Getenv("gptModel"); model != "" {
		return model

	}
	return "gpt-4o"
}

func (s *SimpleGptChat) chat(userID, msg string) string {
	cfg := openai.DefaultConfig(s.token)
	cfg.BaseURL = s.url
	client := openai.NewClientWithConfig(cfg)

	var msgs = GetMsgListWithDb(config.Bot_Type_Gpt, userID, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: msg}, s.toDbMsg, s.toChatMsg)
	req := openai.ChatCompletionRequest{
		Model:    s.getModel(userID),
		Messages: msgs,
	}
	// 如果设置了环境变量且合法，则增加maxTokens参数，否则不设置
	if s.maxTokens > 0 {
		req.MaxTokens = s.maxTokens // 参数名称参考：https://github.com/sashabaranov/go-openai
	}
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return err.Error()
	}
	content := resp.Choices[0].Message.Content
	msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: content})
	SaveMsgListWithDb(config.Bot_Type_Gpt, userID, msgs, s.toDbMsg)
	return content
}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	r, flag := DoAction(userID, msg)
	if flag {
		return r
	}
	return WithTimeChat(userID, msg, s.chat)
}
