package chat

import (
	"strings"

	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type KeywordChat struct {
	BaseChat
}

func (k *KeywordChat) Chat(userID string, msg string) string {
	// 检查是否为指令，如果是则交给DoAction处理
	r, flag := DoAction(userID, msg)
	if flag {
		return r
	}

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
	if msg.MsgType == message.MsgTypeEvent {
		// 将事件消息委托给通用的 SimpleChat 处理
		simpleChat := SimpleChat{}
		return simpleChat.HandleMediaMsg(msg)
	}
	return "关键词回复模式不支持处理多媒体消息"
}
