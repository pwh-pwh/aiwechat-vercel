package config

import "os"

const (
	Gpt_Welcome_Reply_Key = "gptWelcomeReply"
)

var (
	Gpt_Welcome_Reply = os.Getenv(Gpt_Welcome_Reply_Key)
)
