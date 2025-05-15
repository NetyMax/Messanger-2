package main

import "time"

type Message struct {
	ID          int
	SenderID    string
	RecipientID string
	Content     string
	SentAt      time.Time
}
