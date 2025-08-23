package chat

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/qqhsx/aiwechat-vercel/config"
)

type SparkChat struct {
	Cfg *config.SparkConfig
}

type ChatResponse struct {
	Header  map[string]interface{} `json:"header"`
	Payload struct {
		Choices struct {
			Status int                      `json:"status"`
			Seq    int                      `json:"seq"`
			Text   []map[string]interface{} `json:"text"`
		} `json:"choices"`
	} `json:"payload"`
}

// 生成鉴权 URL
func (s *SparkChat) genAuthUrl(hostUrl string) (string, error) {
	u, err := url.Parse(hostUrl)
	if err != nil {
		return "", err
	}
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signatureOrigin := fmt.Sprintf("host: %s\ndate: %s\nGET %s HTTP/1.1", u.Host, date, u.Path)
	h := hmac.New(sha256.New, []byte(s.Cfg.ApiSecret))
	h.Write([]byte(signatureOrigin))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	authorizationOrigin := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line", signature="%s"`, s.Cfg.ApiKey, signature)
	authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))
	v := url.Values{}
	v.Add("authorization", authorization)
	v.Add("date", date)
	v.Add("host", u.Host)
	return hostUrl + "?" + v.Encode(), nil
}

// 核心聊天函数
func (s *SparkChat) Chat(ctx context.Context, uid string, query string) (string, error) {
	authUrl, err := s.genAuthUrl(s.Cfg.HostUrl)
	if err != nil {
		return "", fmt.Errorf("生成鉴权URL失败: %v", err)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, authUrl, nil)
	if err != nil {
		return "", fmt.Errorf("连接WebSocket失败: %v", err)
	}
	defer conn.Close()

	// 组装请求
	req := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": s.Cfg.AppId,
			"uid":    uid,
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      s.Cfg.SparkDomainVersion, // ⚠️ 使用 config 中的 domain
				"temperature": 0.5,
				"max_tokens":  2048,
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": []map[string]string{
					{"role": "user", "content": query},
				},
			},
		},
	}

	if err := conn.WriteJSON(req); err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}

	var sb strings.Builder
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var resp ChatResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			continue
		}

		// 解析返回内容
		for _, item := range resp.Payload.Choices.Text {
			var msg string
			if val, ok := item["content"]; ok && val != nil {
				if str, ok2 := val.(string); ok2 {
					msg = str
				}
			}
			if msg == "" {
				if val, ok := item["reasoning_content"]; ok && val != nil {
					if str, ok2 := val.(string); ok2 {
						msg = str
					}
				}
			}
			sb.WriteString(msg)
		}

		// status=2 表示一次对话完成
		if resp.Payload.Choices.Status == 2 {
			break
		}
	}

	return sb.String(), nil
}
