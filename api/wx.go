package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/chat"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func Wx(rw http.ResponseWriter, req *http.Request) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "",
		AppSecret: "",
		Token:     config.GetWxToken(),
		Cache:     memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//回复消息：演示回复用户发送的消息
		replyMsg := handleWxMessage(msg)
		text := message.NewText(replyMsg)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

func handleWxMessage(msg *message.MixMessage) (replyMsg string) {
	msgType := msg.MsgType
	msgContent := msg.Content
	userId := string(msg.FromUserName)

	// Check if user is authenticated (only if ADDME_PASSWORD is set)
	if config.GetAddMePassword() != "" && !config.IsUserAuthenticated(userId) {
		if msgType == message.MsgTypeText {
			// Only allow /addme command for non-authenticated users
			if msgContent == "/addme" || len(msgContent) > len("/addme") && msgContent[:len("/addme")] == "/addme" {
				bot := chat.GetChatBot(config.GetUserBotType(userId))
				replyMsg = bot.Chat(userId, msgContent)
			} else {
				replyMsg = "功能还在开发中"
			}
		} else {
			replyMsg = "功能还在开发中"
		}
		return
	}

	bot := chat.GetChatBot(config.GetUserBotType(userId))
	if msgType == message.MsgTypeText {
		replyMsg = bot.Chat(userId, msgContent)
	} else {
		replyMsg = bot.HandleMediaMsg(msg)
	}

	return
}
