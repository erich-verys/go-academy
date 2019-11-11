package server

import (
	"fmt"
	"net/http"
)

func Run() {
	fmt.Println("starting server at :8888")
	http.HandleFunc("/status", status)
	http.ListenAndServe(":8888", nil)
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request received")
	w.WriteHeader(200)
	w.Write([]byte(`{"success":true}`))
}
