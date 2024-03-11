package chat

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/pwh-pwh/aiwechat-vercel/db"
)

func TestGeminiChat(t *testing.T) {
	godotenv.Load("../conf/.env")
	db.ChatDbInstance, _ = db.GetChatDb()

	chat := &GeminiChat{
		BaseChat: SimpleChat{},
		key:      config.GetGeminiKey(),
	}

	res := chat.Chat("testUser", "用10个字描述你的能力")

	fmt.Println(res)
}
