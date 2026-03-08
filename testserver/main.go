package main

import (
	"fmt"
	"net/http"
	"testserver/handler"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello this is v1\n")
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/hi", handler.HiHandler)

	http.ListenAndServe(":8090", nil)
}
