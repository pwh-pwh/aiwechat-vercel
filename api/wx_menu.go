package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/pwh-pwh/aiwechat-vercel/config"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	oaConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/menu"
)

const (
	Opt_Query  = "query"
	Opt_Create = "create"
	Opt_Delete = "delete"
)

func WxMenu(rpn http.ResponseWriter, req *http.Request) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &oaConfig.Config{
		AppID:     config.GetWxAppId(),
		AppSecret: config.GetWxAppSecret(),
		Token:     config.GetWxToken(),
		Cache:     memory,
	}
	oa := wc.GetOfficialAccount(cfg)

	wxMenu := oa.GetMenu()

	opt := req.URL.Query().Get("opt")
	if opt == "" {
		opt = Opt_Query
	}
	switch opt {
	case Opt_Query:
		menus, err := wxMenu.GetMenu()
		if err != nil {
			rpn.Write([]byte(err.Error()))
			return
		}
		json, _ := sonic.Marshal(menus)
		rpn.Write([]byte(json))
	case Opt_Create:
		var buttons []*menu.Button
		body, _ := io.ReadAll(req.Body)
		sonic.Unmarshal(body, &buttons)
		err := wxMenu.SetMenu(buttons)
		if err != nil {
			rpn.Write([]byte(err.Error()))
		} else {
			rpn.Write([]byte("create menu success"))
		}
	case Opt_Delete:
		menuIdStr := req.URL.Query().Get("menuId")
		if menuIdStr == "" {
			return
		}
		menuId, err := strconv.ParseInt(menuIdStr, 10, 64)
		if err != nil {
			rpn.Write([]byte(err.Error()))
			return
		}
		err = wxMenu.DeleteConditional(menuId)
		if err != nil {
			rpn.Write([]byte(err.Error()))
		} else {
			rpn.Write([]byte("delete menu success"))
		}
	default:
		rpn.Write([]byte("unknown opt"))
	}

}
