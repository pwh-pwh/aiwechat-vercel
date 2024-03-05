package chat

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/pwh-pwh/aiwechat-vercel/db"
)

const (
	QwenChatUser = "user"
	QwenChatBot  = "assistant"
)

type QwenChat struct {
	BaseChat
	Config *config.QwenConfig
}

type QwenRequest struct {
	Model string `json:"model"`
	Input Input  `json:"input"`
}

type Input struct {
	Messages []QwenMessage `json:"messages"`
}

type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenResponse struct {
	Output    Output `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestID string `json:"request_id"`
}

type Output struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	OutputTokens int `json:"output_tokens"`
	InputTokens  int `json:"input_tokens"`
}

func (chat *QwenChat) Chat(userId, message string) (res string) {
	return WithTimeChat(userId, message, chat.chat)
}

func (chat *QwenChat) chat(userId string, message string) (res string) {
	chatDb := db.ChatDbInstance
	var msgs = []QwenMessage{
		{
			Role:    QwenChatUser,
			Content: message,
		},
	}
	if chatDb != nil {
		msgList, err := chatDb.GetMsgList(config.Bot_Type_Qwen, userId)
		if err == nil {
			list := toQwenMsgList(msgList)
			msgs = append(list, msgs...)
		}
	}
	qwenReq := QwenRequest{
		Model: chat.Config.ModelVersion,
		Input: Input{Messages: msgs},
	}

	body, _ := sonic.Marshal(qwenReq)
	req, err := http.NewRequest("POST", chat.Config.HostUrl, bytes.NewReader(body))
	if err != nil {
		res = fmt.Sprintf("NewRequest failed,err:%v", err.Error())
		return
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+chat.Config.ApiKey)
	client := http.Client{}
	// 发送请求
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		res = fmt.Sprintf("client.Do failed,err:%v", err.Error())
	}
	defer resp.Body.Close()

	var rpnBody []byte
	rpnBody, err = io.ReadAll(resp.Body)
	if err != nil {
		res = fmt.Sprintf("read http response failed,error=%v", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		res = string(rpnBody)
		return
	}

	// 读取响应
	var qwenRpn QwenResponse
	err = sonic.Unmarshal(rpnBody, &qwenRpn)
	if err != nil {
		res = fmt.Sprintf("Unmarshal response body failed,err:%v", err.Error())
		return
	}

	res = qwenRpn.Output.Text

	if chatDb != nil {
		go func() {
			msgs = append(msgs, QwenMessage{
				Role:    QwenChatBot,
				Content: res,
			})
			chatDb.SetMsgList(config.Bot_Type_Qwen, userId, toDbMsgList(msgs))
		}()
	}

	return
}

func toDbMsgList(msgList []QwenMessage) []db.Msg {
	var messages []db.Msg
	for _, msg := range msgList {
		messages = append(messages, db.Msg{
			Role: msg.Role,
			Msg:  msg.Content,
		})
	}
	return messages
}

func toQwenMsgList(msgList []db.Msg) []QwenMessage {
	var messages []QwenMessage
	for _, msg := range msgList {
		messages = append(messages, QwenMessage{
			Role:    msg.Role,
			Content: msg.Msg,
		})
	}
	return messages
}
