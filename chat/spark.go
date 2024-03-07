package chat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pwh-pwh/aiwechat-vercel/db"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/pwh-pwh/aiwechat-vercel/config"
)

type SparkChat struct {
	BaseChat
	Config *config.SparkConfig
}

type SparkResponse struct {
	Header  *SparkResponseHeader `json:"header"`
	Payload map[string]any       `json:"payload"`
}

type SparkResponseHeader struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Status  int    `json:"status"`
}

func (header *SparkResponseHeader) IsSuccess() bool {
	return header.Code == 0
}

func (header *SparkResponseHeader) IsFailed() bool {
	return !header.IsSuccess()
}

func (header *SparkResponseHeader) ToString() string {
	buf, _ := sonic.Marshal(header)
	return string(buf)
}

func (chat *SparkChat) Chat(userId, message string) (res string) {
	return WithTimeChat(userId, message, chat.chat)
}

func (chat *SparkChat) chat(userId string, message string) (res string) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := dialer.Dial(assembleAuthUrl1(chat.Config.HostUrl, chat.Config.ApiKey, chat.Config.ApiSecret), nil)
	if err != nil {
		res = readResp(resp) + err.Error()
		return
	} else if resp.StatusCode != 101 {
		res = readResp(resp)
		return
	}
	/*var msgs = []SparkMessage{
		{
			Role:    "user",
			Content: message,
		},
	}
	chatDb := db.ChatDbInstance
	if chatDb != nil {
		msgList, err := chatDb.GetMsgList(config.Bot_Type_Spark, userId)
		if err == nil {
			list := toSparkMsgList(msgList)
			msgs = append(list, msgs...)
		}
	}*/
	var msgs = GetMsgListWithDb(config.Bot_Type_Spark, userId, SparkMessage{
		Role:    "user",
		Content: message,
	}, chat.toDbMsg, chat.toChatMsg)

	go func() {
		data := generateRequestBody(chat.Config.AppId, chat.Config.SparkDomainVersion, msgs)
		conn.WriteJSON(data)
	}()

	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read message error:", err)
			break
		}

		fmt.Println(string(msg))

		var rpn SparkResponse
		err = sonic.Unmarshal(msg, &rpn)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}
		if rpn.Header.IsFailed() {
			res = rpn.Header.ToString()
			return
		}
		//解析数据
		choices := rpn.Payload["choices"].(map[string]interface{})
		status := choices["status"].(float64)
		fmt.Println(status)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			res += content
		} else {
			fmt.Println("收到最终结果")
			res += content
			usage := rpn.Payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			totalTokens := temp["total_tokens"].(float64)
			fmt.Println("total_tokens:", totalTokens)
			conn.Close()
			break
		}

	}
	msgs = append(msgs, SparkMessage{
		Role:    "assistant",
		Content: res,
	})
	SaveMsgListWithDb(config.Bot_Type_Spark, userId, msgs, chat.toDbMsg)
	return
}

func (s *SparkChat) toDbMsg(msg SparkMessage) db.Msg {
	return db.Msg{
		Role: msg.Role,
		Msg:  msg.Content,
	}
}

func (s *SparkChat) toChatMsg(msg db.Msg) SparkMessage {
	return SparkMessage{
		Role:    msg.Role,
		Content: msg.Msg,
	}
}

func toSparkMsgList(msgList []db.Msg) []SparkMessage {
	var messages []SparkMessage
	for _, msg := range msgList {
		messages = append(messages, SparkMessage{
			Role:    msg.Role,
			Content: msg.Msg,
		})
	}
	return messages
}

func toMsgList(msgList []SparkMessage) []db.Msg {
	var messages []db.Msg
	for _, msg := range msgList {
		messages = append(messages, db.Msg{
			Role: msg.Role,
			Msg:  msg.Content,
		})
	}
	return messages
}

// 生成参数
func generateRequestBody(appid string, domain string, messages []SparkMessage) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名
	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": appid, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      domain,       // 根据实际情况修改返回的数据结构和字段名
				"temperature": float64(0.8), // 根据实际情况修改返回的数据结构和字段名
				"top_k":       int64(6),     // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  int64(2048),  // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",    // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac api_key=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

type SparkMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
