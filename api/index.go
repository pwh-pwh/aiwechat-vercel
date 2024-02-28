package api

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "<h1>Hello Aiwechat-Vercel!</h1>")
}
