package chat

import (
	_ "errors"
	"fmt"
	"github.com/pwh-pwh/aiwechat-vercel/client"
	"os"
	"strconv"
	"strings"
	"time"
	// "bytes"
	// "io"
	// "mime/multipart"
    // "net/http"

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

	config.Wx_Command_Prompt:    SetPrompt,
	config.Wx_Command_RmPrompt:  RmPrompt,
	config.Wx_Command_GetPrompt: GetPrompt,

	config.Wx_Command_SetModel: SetModel,
	config.Wx_Command_GetModel: GetModel,
	config.Wx_Command_Clear:    ClearMsg,

	config.Wx_Todo_List: GetTodoList,
	config.Wx_Todo_Add:  AddTodo,
	config.Wx_Todo_Del:  DelTodo,

	config.Wx_Coin: GetCoin,
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
	// case message.MsgTypeVoice:
	// 	voiceURL := msg.voiceURL
	// 	text, err = ConvertVoiceToText(voiceURL)
	// 	if err != nil {
	// 		return "Sorry, I couldn't understand that voice message."
	// 	}
	// 	return s.Chat(userID, text)
	default:
		return "未支持的类型"
	}
}

// func ConvertVoiceToText(voiceURL string) (string, error) {
// 	// Open the file
//     file, err := os.Open(filePath)
//     if err != nil {
//         return "", err
//     }
//     defer file.Close()

// 	// Prepare the multipart request
//     body := &bytes.Buffer{}
//     writer := multipart.NewWriter(body)
//     part, err := writer.CreateFormFile("file", filePath)
//     if err != nil {
//         return "", err
//     }
//     _, err = io.Copy(part, file)
//     if err != nil {
//         return "", err
//     }
//     writer.Close()

// 	// Create the HTTP request
//     req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", body)
//     if err != nil {
//         return "", err
//     }
//     req.Header.Set("Authorization", "Bearer "+config.GetGptToken())
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
//     // Add any necessary headers here, such as Authorization

//     // Send the request
//     client := &http.Client{}
//     resp, err := client.Do(req)
//     if err != nil {
//         return "", err
//     }
//     defer resp.Body.Close()

// 	respBody, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return "", err
//     }

// 	return string(respBody), nil
// }

func SwitchUserBot(userId string, botType string) string {
	if _, err := config.CheckBotConfig(botType); err != nil {
		return err.Error()
	}
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

func RmPrompt(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	db.RemovePrompt(userId, botType)
	return fmt.Sprintf("%s 删除prompt成功", botType)
}

func GetPrompt(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	prompt, err := db.GetPrompt(userId, botType)
	if err != nil {
		return fmt.Sprintf("%s 当前未设置prompt", botType)
	}
	return fmt.Sprintf("%s 获取prompt成功，prompt：%s", botType, prompt)
}

func GetTodoList(param string, userId string) string {
	list, err := db.GetTodoList(userId)
	if err != nil {
		return err.Error()
	}
	return list
}

func AddTodo(param, userId string) string {
	err := db.AddTodoList(userId, param)
	if err != nil {
		return err.Error()
	}
	return "添加成功"
}

func DelTodo(param, userId string) string {
	index, err := strconv.Atoi(param)
	if err != nil {
		return "传入索引必须为数字"
	}
	err = db.DelTodoList(userId, index)
	if err != nil {
		return err.Error()
	}
	return "删除todo成功"
}

func GetCoin(param, userId string) string {
	coinPrice, err := client.GetCoinPrice(param)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("代币对:%s 价格:%s", coinPrice.Symbol, coinPrice.Price)
}

func SetModel(param, userId string) string {
	botType := config.GetUserBotType(userId)
	if botType == config.Bot_Type_Gpt || botType == config.Bot_Type_Gemini || botType == config.Bot_Type_Qwen {
		if err := db.SetModel(userId, botType, param); err != nil {
			return fmt.Sprintf("%s 设置model失败", botType)
		}
		return fmt.Sprintf("%s 设置model成功", botType)
	}
	return fmt.Sprintf("%s 不支持设置model", botType)
}

func GetModel(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	model, err := db.GetModel(userId, botType)
	if err != nil || model == "" {
		return fmt.Sprintf("%s 当前未设置model", botType)
	}
	return fmt.Sprintf("%s 获取model成功，model：%s", botType, model)
}

func ClearMsg(param string, userId string) string {
	botType := config.GetUserBotType(userId)
	db.DeleteMsgList(botType, userId)
	return fmt.Sprintf("%s 清除消息成功", botType)
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
	maxTokens := config.GetMaxTokens()

	switch botType {
	case config.Bot_Type_Gpt:
		url := os.Getenv("GPT_URL")
		if url == "" {
			url = "https://api.openai.com/v1/"
		}
		return &SimpleGptChat{
			token:     config.GetGptToken(),
			url:       url,
			maxTokens: maxTokens,
			BaseChat:  SimpleChat{},
		}
	case config.Bot_Type_Gemini:
		return &GeminiChat{
			BaseChat:  SimpleChat{},
			key:       config.GetGeminiKey(),
			maxTokens: maxTokens,
		}
	case config.Bot_Type_Spark:
		config, _ := config.GetSparkConfig()
		return &SparkChat{
			BaseChat:  SimpleChat{},
			Config:    config,
			maxTokens: maxTokens,
		}
	case config.Bot_Type_Qwen:
		config, _ := config.GetQwenConfig()
		return &QwenChat{
			BaseChat:  SimpleChat{},
			Config:    config,
			maxTokens: maxTokens,
		}
	default:
		return &Echo{}
	}
}

type ChatMsg interface {
	openai.ChatCompletionMessage | QwenMessage | SparkMessage | *genai.Content
}

func GetMsgListWithDb[T ChatMsg](botType, userId string, msg T, f func(msg T) db.Msg, f2 func(msg db.Msg) T) []T {
	var dbList []db.Msg
	isSupportPrompt := config.IsSupportPrompt(botType)
	if isSupportPrompt {
		prompt, err := db.GetPrompt(userId, botType)
		if err == nil && prompt != "" {
			dbList = append(dbList, db.Msg{
				Role: "system",
				Msg:  prompt,
			})
		}
	}
	if db.ChatDbInstance != nil {
		list, err := db.ChatDbInstance.GetMsgList(botType, userId)
		if err == nil {
			// check is contain system prompt
			if len(list) > 0 {
				if list[0].Role == "system" {
					list = list[1:]
				}
			}
			dbList = append(dbList, list...)
		}
	}
	dbList = append(dbList, f(msg))
	r := make([]T, 0)
	for _, msg := range dbList {
		r = append(r, f2(msg))
	}
	return r
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
