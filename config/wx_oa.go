package config

import (
	"os"
	"strings"
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
)

var (
	Wx_Commands = []string{Wx_Command_Help, Wx_Command_Gpt, Wx_Command_Spark, Wx_Command_Qwen, Wx_Command_Gemini}
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
		helpMsg = "选择功能吧 🎯\n📖 查看帮助（/help）\n🤖 与 GPT 对话（/gpt）\n🚀 与星火对话（/spark）\n🐦 与通义千问对话（/qwen）\n🌟 与 gemini 对话（/gemini）\n" +
			"✍️ 设置 system prompt（/prompt 你的 prompt）\n📄 获取当前设置 prompt（/getpt）\n🧹 清除当前设置 prompt（/cpt）\n" +
			"🛠️ 设置自定义 model（/setmodel model）\n🔧 重置 model 为默认值（/setmodel）\n📋 获取当前 model（/getmodel）\n" +
			"🗑️ 清除历史对话（/clear）\n" + "✅ 设置 todo（/ta 代办事项 1）\n" + "📃 获取代办列表（/tl）\n" + "❌ 删除索引代办事件（/td 2）\n" + "💰 查询价格（/cb 代币对）"
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
