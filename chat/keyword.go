package chat

import (
	"fmt"
	"strings"

	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type KeywordChat struct {
	BaseChat
}

func (k *KeywordChat) Chat(userID string, msg string) string {
	replies, err := db.GetKeywordReplies()
	if err != nil {
		return "获取关键词回复失败"
	}

	for _, reply := range replies {
		if strings.Contains(msg, reply.Keyword) {
			return reply.Reply
		}
	}

	return "未找到匹配的关键词，请尝试其他内容"
}

func (k *KeywordChat) HandleMediaMsg(msg *message.MixMessage) string {
	return "关键词回复模式不支持处理多媒体消息"
}