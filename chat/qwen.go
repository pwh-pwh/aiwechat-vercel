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
	var msgs = GetMsgListWithDb(config.Bot_Type_Qwen, userId, QwenMessage{
		Role:    QwenChatUser,
		Content: message,
	}, chat.toDbMsg, chat.toChatMsg)

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

	msgs = append(msgs, QwenMessage{
		Role:    QwenChatBot,
		Content: res,
	})
	SaveMsgListWithDb(config.Bot_Type_Qwen, userId, msgs, chat.toDbMsg)
	return
}

func (s *QwenChat) toDbMsg(msg QwenMessage) db.Msg {
	return db.Msg{
		Role: msg.Role,
		Msg:  msg.Content,
	}
}

func (s *QwenChat) toChatMsg(msg db.Msg) QwenMessage {
	return QwenMessage{
		Role:    msg.Role,
		Content: msg.Msg,
	}
}
