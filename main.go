package main

import (
	"log"
	"net/http"
)

func main() {

	// http.HandleFunc("/", rootHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		NewWebSocket(w, r)
	})

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
