package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSource string) error {
	var err error
	DB, err = sql.Open("postgres", dataSource)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func SaveMessage(senderID, recipientID, content string) error {
	_, err := DB.Exec(`
		INSERT INTO messages (sender_id, recipient_id, content)
		VALUES ($1, $2, $3)
	`, senderID, recipientID, content)
	return err
}

type Message struct {
	ID          int
	SenderID    string
	RecipientID string
	Content     string
	SentAt      time.Time
}
