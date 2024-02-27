package api

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "h1>Hello from Go!</h1>")
}
