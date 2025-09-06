package config

import (
	"os"
	"strings"
)

const (
	Gpt_Welcome_Reply_Key = "gptWelcomeReply"
	Gpt_Token             = "GPT_TOKEN"
	DefaultGptWelcome     = "我是gpt，开始聊天吧！"
)

// GetGptWelcomeReply returns the welcome message for GPT bot
func GetGptWelcomeReply() string {
	if reply := os.Getenv(Gpt_Welcome_Reply_Key); reply != "" {
		return strings.TrimSpace(reply)
	}
	return DefaultGptWelcome
}

// GetGptToken returns the GPT API token
func GetGptToken() string {
	return strings.TrimSpace(os.Getenv(Gpt_Token))
}

// IsGptConfigured checks if GPT is properly configured
func IsGptConfigured() bool {
	return GetGptToken() != ""
}