package entities

import "time"

type Message struct {
	SenderName string    `json:"sender_name" db:"sender_name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	Content    string    `json:"content" db:"content"`
}
