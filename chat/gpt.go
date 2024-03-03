package chat

import (
	"context"
	"fmt"
	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/sashabaranov/go-openai"
)

type SimpleGptChat struct {
	token string
	url   string
	SimpleChat
}

func (s *SimpleGptChat) toGptMsgList(msgList []db.Msg) []openai.ChatCompletionMessage {
	var result []openai.ChatCompletionMessage
	for _, item := range msgList {
		result = append(result, openai.ChatCompletionMessage{
			Role:    item.Role,
			Content: item.Msg,
		})
	}
	return result
}

func (s *SimpleGptChat) toChatMsList(msgList []openai.ChatCompletionMessage) []db.Msg {
	var result []db.Msg
	for _, item := range msgList {
		result = append(result, db.Msg{
			Role: item.Role,
			Msg:  item.Content,
		})
	}
	return result
}

func (s *SimpleGptChat) chat(userID, msg string) string {
	cfg := openai.DefaultConfig(s.token)
	cfg.BaseURL = s.url
	client := openai.NewClientWithConfig(cfg)
	var msgs = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		},
	}
	chatDb := db.ChatDbInstance
	if chatDb != nil {
		msgList, err := chatDb.GetMsgList(userID)
		if err != nil {
			list := s.toGptMsgList(msgList)
			msgs = append(list, msgs...)
		}
	}
	// log msgs
	fmt.Println(msgs)
	resp, err := client.CreateChatCompletion(context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: msgs,
		})
	if err != nil {
		return err.Error()
	}
	content := resp.Choices[0].Message.Content
	go func() {
		msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleSystem, Content: content})
		chatDb.SetMsgList(userID, s.toChatMsList(msgs))
	}()
	return content
}

func (s *SimpleGptChat) Chat(userID string, msg string) string {
	return WithTimeChat(userID, msg, s.chat)
}
