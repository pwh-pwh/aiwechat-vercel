package api

import (
	"fmt"
	"github.com/pwh-pwh/aiwechat-vercel/chat"
	"golang.org/x/text/encoding/charmap"
	"net/http"
)

func Chat(rw http.ResponseWriter, req *http.Request) {
	msg := req.URL.Query().Get("msg")
	botType := req.URL.Query().Get("botType")
	if msg == "" {
		msg = "用10个字介绍你自己"
	}
	bot := chat.GetChatBot(botType)
	rpn := bot.Chat("admin", msg)
	encoder := charmap.Windows1252.NewEncoder()
	s, e := encoder.String(rpn)
	if e != nil {
		fmt.Fprint(rw, e.Error())
		return
	}
	fmt.Fprint(rw, s)
}
