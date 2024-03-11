package chat

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"

	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/sashabaranov/go-openai"

	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

var actionMap = map[string]func(param, userId string) string{
	config.Wx_Command_Help: func(param, userId string) string {
		return config.GetWxHelpReply()
	},
	config.Wx_Command_Gpt: func(param, userId string) string {
		return SwitchUserBot(userId, config.Bot_Type_Gpt)
	},
	config.Wx_Command_Spark: func(param, userId string) string {
		return SwitchUserBot(userId, config.Bot_Type_Spark)
	},
	config.Wx_Command_Qwen: func(param, userId string) string {
		return SwitchUserBot(userId, config.Bot_Type_Qwen)
	},
	config.Wx_Command_Gemini: func(param, userId string) string {
		return SwitchUserBot(userId, config.Bot_Type_Gemini)
	},
	config.Wx_Command_Prompt: SetPrompt,
}

func DoAction(userId, msg string) (r string, flag bool) {
	action, param, flag := isAction(msg)
	if flag {
		f := actionMap[action]
		r = f(param, userId)
	}
	return
}

func isAction(msg string) (string, string, bool) {
	for key := range actionMap {
		if strings.HasPrefix(msg, key) {
			return msg[:len(key)], strings.TrimSpace(msg[len(key):]), true
		}
	}
	return "", "", false
}

type BaseChat interface {
	Chat(userID string, msg string) string
	HandleMediaMsg(msg *message.MixMessage) string
}
type SimpleChat struct {
}

func (s SimpleChat) Chat(userID string, msg string) string {
	panic("implement me")
}

func (s SimpleChat) HandleMediaMsg(msg *message.MixMessage) string {
	switch msg.MsgType {
	case message.MsgTypeImage:
		return msg.PicURL
	case message.MsgTypeEvent:
		if msg.Event == message.EventSubscribe {
			subText := config.GetWxSubscribeReply() + config.GetWxHelpReply()
			if subText == "" {
				subText = "哇，又有帅哥美女关注我啦😄"
			}
			return subText
		} else if msg.Event == message.EventClick {
			switch msg.EventKey {
			case config.GetWxEventKeyChatGpt():
				return SwitchUserBot(string(msg.FromUserName), config.Bot_Type_Gpt)
			case config.GetWxEventKeyChatSpark():
				return SwitchUserBot(string(msg.FromUserName), config.Bot_Type_Spark)
			case config.GetWxEventKeyChatQwen():
				return SwitchUserBot(string(msg.FromUserName), config.Bot_Type_Qwen)
			default:
				return fmt.Sprintf("unkown event key=%v", msg.EventKey)
			}
		} else {
			return "不支持的类型"
		}
	default:
		return "未支持的类型"
	}
}

func SwitchUserBot(userId string, botType string) string {
	db.SetValue(fmt.Sprintf("%v:%v", config.Bot_Type_Key, userId), botType, 0)
	return config.GetBotWelcomeReply(botType)
}

func SetPrompt(param, userId string) string {
	botType := config.GetUserBotType(userId)
	switch botType {
	case config.Bot_Type_Gpt:
		db.SetPrompt(userId, botType, param)
	case config.Bot_Type_Qwen:
		db.SetPrompt(userId, botType, param)
	case config.Bot_Type_Spark:
		db.SetPrompt(userId, botType, param)
	default:
		return fmt.Sprintf("%s 不支持设置system prompt", botType)
	}
	return fmt.Sprintf("%s 设置prompt成功", botType)
}

// 加入超时控制
func WithTimeChat(userID, msg string, f func(userID, msg string) string) string {
	if _, ok := config.Cache.Load(userID + msg); ok {
		rAny, _ := config.Cache.Load(userID + msg)
		r := rAny.(string)
		config.Cache.Delete(userID + msg)
		return r
	}
	resChan := make(chan string)
	go func() {
		resChan <- f(userID, msg)
	}()
	select {
	case res := <-resChan:
		return res
	case <-time.After(5 * time.Second):
		config.Cache.Store(userID+msg, <-resChan)
		return ""
	}
}

type ErrorChat struct {
	errMsg string
}

func (e *ErrorChat) HandleMediaMsg(msg *message.MixMessage) string {
	return e.errMsg
}

func (e *ErrorChat) Chat(userID string, msg string) string {
	return e.errMsg
}

func GetChatBot(botType string) BaseChat {
	if botType == "" {
		botType = config.GetBotType()
	}
	var err error
	botType, err = config.CheckBotConfig(botType)
	if err != nil {
		return &ErrorChat{
			errMsg: err.Error(),
		}
	}

	switch botType {
	case config.Bot_Type_Gpt:
		url := os.Getenv("GPT_URL")
		if url == "" {
			url = "https://api.openai.com/v1/"
		}
		return &SimpleGptChat{
			token:    config.GetGptToken(),
			url:      url,
			BaseChat: SimpleChat{},
		}
	case config.Bot_Type_Gemini:
		return &GeminiChat{
			BaseChat: SimpleChat{},
			key:      config.GetGeminiKey(),
		}
	case config.Bot_Type_Spark:
		config, _ := config.GetSparkConfig()
		return &SparkChat{
			BaseChat: SimpleChat{},
			Config:   config,
		}
	case config.Bot_Type_Qwen:
		config, _ := config.GetQwenConfig()
		return &QwenChat{
			BaseChat: SimpleChat{},
			Config:   config,
		}
	default:
		return &Echo{}
	}
}

type ChatMsg interface {
	openai.ChatCompletionMessage | QwenMessage | SparkMessage | *genai.Content
}

func GetMsgListWithDb[T ChatMsg](botType, userId string, msg T, f func(msg T) db.Msg, f2 func(msg db.Msg) T) []T {
	if db.ChatDbInstance != nil {
		list, err := db.ChatDbInstance.GetMsgList(botType, userId)
		isSupportPrompt := config.IsSupportPrompt(botType)
		if isSupportPrompt {
			prompt, err := db.GetPrompt(userId, botType)
			if err == nil {
				list = append([]db.Msg{
					{
						Role: "system",
						Msg:  prompt,
					},
				}, list...)
			}
		}
		if err == nil {
			list = append(list, f(msg))
			r := make([]T, 0)
			for _, msg := range list {
				r = append(r, f2(msg))
			}
			return r
		}
	}
	return []T{msg}
}

func SaveMsgListWithDb[T ChatMsg](botType, userId string, msgList []T, f func(msg T) db.Msg) {
	if db.ChatDbInstance != nil {
		go func() {
			list := make([]db.Msg, 0)
			for _, msg := range msgList {
				list = append(list, f(msg))
			}
			db.ChatDbInstance.SetMsgList(botType, userId, list)
		}()
	}
}
