package api

import (
	"fmt"
	"net/http"
)

func Check(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "check ok")
}
