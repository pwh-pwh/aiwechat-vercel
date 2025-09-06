package config

import (
	"os"
	"strings"
)

const (
	Gemini_Welcome_Reply_Key = "geminiWelcomeReply"
	Gemini_Key               = "geminiKey"
	DefaultGeminiWelcome     = "我是gemini，开始聊天吧！"
)

// GetGeminiWelcomeReply returns the welcome message for Gemini bot
func GetGeminiWelcomeReply() string {
	if reply := os.Getenv(Gemini_Welcome_Reply_Key); reply != "" {
		return strings.TrimSpace(reply)
	}
	return DefaultGeminiWelcome
}

// GetGeminiKey returns the Gemini API key
func GetGeminiKey() string {
	return strings.TrimSpace(os.Getenv(Gemini_Key))
}

// IsGeminiConfigured checks if Gemini is properly configured
func IsGeminiConfigured() bool {
	return GetGeminiKey() != ""
}