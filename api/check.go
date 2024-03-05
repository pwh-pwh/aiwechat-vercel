package api

import (
	"fmt"
	"net/http"

	"github.com/pwh-pwh/aiwechat-vercel/config"
)

func Check(rw http.ResponseWriter, req *http.Request) {
	botType, checkRes := config.CheckAllBotConfig()
	var res string
	for bot, status := range checkRes {
		if res == "" {
			res = fmt.Sprintf("%v: %v", bot, status)
		} else {
			res = fmt.Sprintf("%v\n%v: %v", res, bot, status)
		}
	}
	res = fmt.Sprintf("%v\nDEFAULT BOT: %v", res, botType)
	fmt.Fprint(rw, res)
}
