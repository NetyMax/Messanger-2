package main

import (
	"log"
	"net/http"

	"github.com/Max/messenger-2/ws"
)

func main() {
	http.HandleFunc("/ws", ws.HandleWebSocket)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
