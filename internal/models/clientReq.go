package models

type Message struct {
	MessageID int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
}
