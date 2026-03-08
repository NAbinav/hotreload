package handler

import (
	"fmt"
	"net/http"
)

func HiHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hi guy\n")
}
