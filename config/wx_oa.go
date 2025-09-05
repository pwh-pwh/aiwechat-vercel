package config

import (
	"os"
	"strings"

	"github.com/pwh-pwh/aiwechat-vercel/db"
)

const (
	Wx_Token_key           = "WX_TOKEN"
	Wx_App_Id_key          = "WX_APP_ID"
	Wx_App_Secret_key      = "WX_APP_SECRET"
	Wx_Subscribe_Reply_key = "WX_SUBSCRIBE_REPLY"
	Wx_Help_Reply_key      = "WX_HELP_REPLY"

	Wx_Event_Key_Chat_Gpt_key   = "AI_CHAT_GPT"
	Wx_Event_Key_Chat_Spark_key = "AI_CHAT_SPARK"
	Wx_Event_Key_Chat_Qwen_key  = "AI_CHAT_QWEN"

	Wx_Command_Help      = "/help"
	Wx_Command_Gpt       = "/gpt"
	Wx_Command_Spark     = "/spark"
	Wx_Command_Qwen      = "/qwen"
	Wx_Command_Gemini    = "/gemini"
	Wx_Command_Prompt    = "/prompt"
	Wx_Command_RmPrompt  = "/cpt"
	Wx_Command_GetPrompt = "/getpt"
	Wx_Command_SetModel  = "/setmodel"
	Wx_Command_GetModel  = "/getmodel"
	Wx_Command_Clear     = "/clear"

	Wx_Todo_Add  = "/ta"
	Wx_Todo_Del  = "/td"
	Wx_Todo_List = "/tl"

	Wx_Coin = "/cb"

	Wx_Command_AddMe = "/addme"
)

const (
	Wx_Subscribe_Password_Key = "WX_SUBSCRIBE_PASSWORD"
	Wx_User_Verified_Prefix    = "verified_user_"
)

var (
	Wx_Commands = []string{Wx_Command_Help, Wx_Command_Gpt, Wx_Command_Spark, Wx_Command_Qwen, Wx_Command_Gemini, Wx_Command_AddMe}
)

func GetWxToken() string {
	return os.Getenv(Wx_Token_key)
}
func GetWxAppId() string {
	return os.Getenv(Wx_App_Id_key)
}
func GetWxAppSecret() string {
	return os.Getenv(Wx_App_Secret_key)
}
func GetWxSubscribeReply() string {
	subscribeMsg := os.Getenv(Wx_Subscribe_Reply_key)
	return strings.ReplaceAll(subscribeMsg, "\\n", "\n")
}
func GetWxHelpReply() string {
	helpMsg := os.Getenv(Wx_Help_Reply_key)
	if helpMsg == "" {
		helpMsg = "输入以下命令进行对话\n/help：查看帮助\n/gpt：与GPT对话\n/spark：与星火对话\n/qwen：与通义千问对话\n/gemini：与gemini对话\n" +
			"/addme 密码：验证身份并使用完整功能\n/prompt 你的prompt: 设置system prompt\n/getpt: 获取当前设置prompt\n/cpt: 清除当前设置prompt\n" +
			"/setmodel model: 设置自定义model\n/setmodel: 重置model为默认值\n/getmodel: 获取当前model\n" +
			"/clear:清除历史对话\n" + "/ta 代办事项1:设置todo\n" + "/tl:获取代办列表\n" + "/td 2:删除索引代办事件\n" + "/cb 代币对:查询价格"
	}
	return strings.ReplaceAll(helpMsg, "\\n", "\n")
}
func GetWxEventKeyChatGpt() string {
	return os.Getenv(Wx_Event_Key_Chat_Gpt_key)
}
func GetWxEventKeyChatSpark() string {
	return os.Getenv(Wx_Event_Key_Chat_Spark_key)
}
func GetWxEventKeyChatQwen() string {
	return os.Getenv(Wx_Event_Key_Chat_Qwen_key)
}

func GetBotWelcomeReply(botType string) string {
	switch botType {
	case Bot_Type_Gpt:
		return GetGptWelcomeReply()
	case Bot_Type_Gemini:
		return GetGeminiWelcomeReply()
	case Bot_Type_Spark:
		return GetSparkWelcomeReply()
	case Bot_Type_Qwen:
		return GetQwenWelcomeReply()
	}

	return botType
}

func GetSubscribePassword() string {
	return os.Getenv(Wx_Subscribe_Password_Key)
}

func IsUserVerified(userId string) bool {
	verifiedKey := Wx_User_Verified_Prefix + userId
	// Check if user is verified in Redis
	if result, err := db.GetValue(verifiedKey); err == nil && result != "" {
		return result == "true"
	}
	return false
}

func VerifyUser(userId string) {
	verifiedKey := Wx_User_Verified_Prefix + userId
	db.SetValue(verifiedKey, "true", 0)
}

func GetDevMessage() string {
	return "功能还在开发中"
}
