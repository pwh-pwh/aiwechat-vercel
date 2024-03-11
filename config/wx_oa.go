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

	Wx_Command_Help  = "/help"
	Wx_Command_Gpt   = "/gpt"
	Wx_Command_Spark = "/spark"
	Wx_Command_Qwen  = "/qwen"
)

var (
	Wx_Token           = os.Getenv(Wx_Token_key)
	Wx_App_Id          = os.Getenv(Wx_App_Id_key)
	Wx_App_Secret      = os.Getenv(Wx_App_Secret_key)
	Wx_Subscribe_Reply = os.Getenv(Wx_Subscribe_Reply_key)
	Wx_Help_Reply      = strings.ReplaceAll(os.Getenv(Wx_Help_Reply_key), "\\n", "\n")

	Wx_Event_Key_Chat_Gpt   = os.Getenv(Wx_Event_Key_Chat_Gpt_key)
	Wx_Event_Key_Chat_Spark = os.Getenv(Wx_Event_Key_Chat_Spark_key)
	Wx_Event_Key_Chat_Qwen  = os.Getenv(Wx_Event_Key_Chat_Qwen_key)

	Wx_Commands = []string{Wx_Command_Help, Wx_Command_Gpt, Wx_Command_Spark, Wx_Command_Qwen}
)
