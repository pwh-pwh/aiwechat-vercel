package config

import "os"

const (
	Claude_Welcome_Reply_Key = "claudeWelcomeReply"
	Claude_Key               = "claudeKey"
	Claude_Url               = "claudeUrl"
	Claude_Model             = "claudeModel"
)

func GetClaudeWelcomeReply() (r string) {
	r = os.Getenv(Claude_Welcome_Reply_Key)
	if r == "" {
		r = "我是Claude，开始聊天吧！"
	}
	return
}

func GetClaudeKey() string {
	return os.Getenv(Claude_Key)
}

func GetClaudeUrl() string {
	url := os.Getenv(Claude_Url)
	if url == "" {
		url = "https://api.anthropic.com"
	}
	return url
}

func GetClaudeModel() string {
	model := os.Getenv(Claude_Model)
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}
	return model
}