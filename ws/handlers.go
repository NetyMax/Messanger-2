package ws

import (
	"encoding/json"

	"github.com/Max/messenger-2/db"
)

type Message struct {
	Type    string `json:"type"`    // "message"
	To      string `json:"to"`      // получатель
	Content string `json:"content"` // текст сообщения
}

func HandleMessage(from string, data []byte) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil || msg.Type != "message" {
		return
	}

	_ = db.SaveMessage(from, msg.To, msg.Content)

	response, _ := json.Marshal(struct {
		From    string `json:"from"`
		Content string `json:"content"`
	}{
		From:    from,
		Content: msg.Content,
	})

	SendToUser(msg.To, response)
}
