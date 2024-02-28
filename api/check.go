package api

import (
	"fmt"
	"net/http"
)

func check(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "check ok")
}
