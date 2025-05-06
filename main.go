package main

import (
	"log"
	"messenger/ws"
	"net/http"
	"os"

	"github.com/Max/messenger-2/db"
)

func main() {
	http.HandleFunc("/ws", ws.HandleWebSocket)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	err := db.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
}
