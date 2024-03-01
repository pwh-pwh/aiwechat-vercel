package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/config"
)

func Check(w http.ResponseWriter, req *http.Request) {
	err := config.CheckBotConfig()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "配置成功")
}
