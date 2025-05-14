package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	if userID == "" {
		http.Error(w, "error", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error", err)
		return
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte),
		userID: userID,
	}

	RegisterClient(userID, client)
	go client.read()
	go client.write()
}

func (c *Client) read() {
	defer func() {
		c.conn.Close()
		UnRegisterCkient(c.userID)
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("error", err)
			break
		}
		HandleMessage(c.userID, message)
	}
}

func (c *Client) write() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
