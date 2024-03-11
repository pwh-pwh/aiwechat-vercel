package config

import "os"

const (
	Gemini_Welcome_Reply_Key = "geminiWelcomeReply"
	Gemini_Key               = "geminiKey"
)

func GetGeminiWelcomeReply() string {
	return os.Getenv(Gemini_Welcome_Reply_Key)
}

func GetGeminiKey() string {
	return os.Getenv(Gemini_Key)
}
