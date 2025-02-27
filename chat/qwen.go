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
	Config    *config.QwenConfig
	maxTokens int
}

type QwenRequest struct {
	Model      string        `json:"model"`
	Message    []QwenMessage `json:"messages"`
	Parameters Parameters    `json:"parameters"`
}

type Input struct {
	Messages []QwenMessage `json:"messages"`
}

type Parameters struct {
	ResultFormat      string   `json:"result_format"`
	Seed              int      `json:"seed"`
	MaxTokens         int      `json:"max_tokens"`
	TopP              float64  `json:"top_p"`
	TopK              float64  `json:"top_k"`
	RepetitionPenalty float64  `json:"repetition_penalty"`
	Temperature       float64  `json:"temperature"`
	Stop              string   `json:"stop"`
	EnableSearch      bool     `json:"enable_search"`
	IncrementalOutput bool     `json:"incremental_output"`
	Tools             []string `json:"tools"`
}

type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type QwenResponse struct {
	Choices   []Choices `json:"choices"`
	Usage     Usage     `json:"usage"`
	RequestID string    `json:"id"`
}

type Choices struct {
	Message      QwenMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	OutputTokens int `json:"output_tokens"`
	InputTokens  int `json:"input_tokens"`
}

func (chat *QwenChat) Chat(userId, message string) (res string) {
	r, flag := DoAction(userId, message)
	if flag {
		return r
	}
	return WithTimeChat(userId, message, chat.chat)
}

func (chat *QwenChat) getModel(userID string) string {
	if model, err := db.GetModel(userID, config.Bot_Type_Qwen); err == nil && model != "" {
		return model
	}
	return chat.Config.ModelVersion
}

func (chat *QwenChat) chat(userId string, message string) (res string) {
	var msgs = GetMsgListWithDb(config.Bot_Type_Qwen, userId, QwenMessage{
		Role:    QwenChatUser,
		Content: message,
	}, chat.toDbMsg, chat.toChatMsg)

	qwenReq := QwenRequest{
		Model:   chat.getModel(userId),
		Message: msgs,
	}
	// 如果设置了环境变量且合法，则增加maxTokens参数，否则不设置
	if chat.maxTokens > 0 {
		qwenReq.Parameters.MaxTokens = chat.maxTokens // 参数名称参考：https://help.aliyun.com/zh/dashscope/developer-reference/api-details
	}
	qwenReq.Parameters.TopP = 0.8              // 通义千问要求top_p ∈ (0,1)
	qwenReq.Parameters.RepetitionPenalty = 1.1 // 用于控制模型生成时的重复度，需要大于0。提高repetition_penalty时可以降低模型生成的重复度。1.0表示不做惩罚。默认为1.1。
	qwenReq.Parameters.Temperature = 0.85      // 取值范围：[0, 2)，系统默认值0.85。不建议取值为0，无意义。

	body, _ := sonic.Marshal(qwenReq)

	fmt.Println(string(body))
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

	res = qwenRpn.Choices[0].Message.Content

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
