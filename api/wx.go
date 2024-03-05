package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pwh-pwh/aiwechat-vercel/chat"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func Wx(rw http.ResponseWriter, req *http.Request) {
	token := os.Getenv("TOKEN")
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "",
		AppSecret: "",
		Token:     token,
		Cache:     memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	bot := chat.GetChatBot("")
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//回复消息：演示回复用户发送的消息
		msgType := msg.MsgType
		replyMsg := ""
		if msgType == message.MsgTypeText {
			replyMsg = bot.Chat(string(msg.FromUserName), msg.Content)
		} else {
			replyMsg = bot.HandleMediaMsg(msg)
		}
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
