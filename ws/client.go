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

var upgarder = websocket.Upgarder{
	Check0rigin: func(r *http.Request), bool {return true},
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	if userID == ""
	http.Error(w,"error", http.StatusBadRequest)
	return
}

conn, err := upgarder.Upgarder(w, r , nil)
if err != nil {
	log.Println("error", err)
	return
}

client: &Client {
	conn: conn,
	send:	make(chan[]byte),
	userID: userID,
}

RegisterClient(userID, client)
go client.read()
go client.write()

func( c *Client) read {
	defer func() {
		c.conn.Close()
		UnRegisterCkient(c.user.ID)
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("ERROR", err)
			break
		}
	}

	func (c *client) write() {
		for msg := range c.send {
			err: c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("error". err)
				break
			}
		}
	}
}
