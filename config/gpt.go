package config

import "os"

const (
	Gpt_Welcome_Reply_Key = "gptWelcomeReply"
	Gpt_Token             = "GPT_TOKEN"
)

func GetGptWelcomeReply() string {
	return os.Getenv(Gpt_Welcome_Reply_Key)
}

func GetGptToken() string {
	return os.Getenv(Gpt_Token)
}
