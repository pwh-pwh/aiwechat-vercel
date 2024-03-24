package config

import "os"

const (
	Gpt_Welcome_Reply_Key = "gptWelcomeReply"
	Gpt_Token             = "GPT_TOKEN"
)

func GetGptWelcomeReply() (r string) {
	r = os.Getenv(Gpt_Welcome_Reply_Key)
	if r == "" {
		r = "gpt-3.5-turbo-16k-0613模型配置成功，我们开始聊天吧！"
	}
	return
}

func GetGptToken() string {
	return os.Getenv(Gpt_Token)
}
