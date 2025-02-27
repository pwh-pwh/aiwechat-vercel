package chat

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/pwh-pwh/aiwechat-vercel/db"
)

func TestQwenChat(t *testing.T) {
	godotenv.Load("D:\\Workspace\\GO\\aiwechat-vercel\\conf\\.env")
	db.ChatDbInstance, _ = db.GetChatDb()

	config, _ := config.GetQwenConfig()
	chat := &QwenChat{
		BaseChat: SimpleChat{},
		Config:   config,
	}

	res := chat.Chat("testUser", "用10个字描述你的能力")

	fmt.Println(res)
}
