package chat

import (
	"fmt"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"

	"github.com/pwh-pwh/aiwechat-vercel/db"
	"github.com/sashabaranov/go-openai"

	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

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
