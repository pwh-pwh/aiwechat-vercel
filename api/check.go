package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/config"
)

func Check(w http.ResponseWriter, req *http.Request) {
	botType, err := config.CheckBotConfig()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "BOT [%v] config check passed", botType)
}
