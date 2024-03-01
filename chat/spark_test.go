package chat

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestSparkChat(t *testing.T) {
	godotenv.Load("../conf/.env")
	chat := SparkChat{}

	res := chat.Chat("testUser", "怎么通过接口配置发布菜单呢")

	fmt.Println(res)
}
