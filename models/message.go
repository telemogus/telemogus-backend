package models

type MessageState int16

const (
	FailedToSend MessageState = iota
	Sent
	Received
	Read
)

type Message struct {
	Base
	ChatID  uint
	UserID  uint
	Content string       `json:"content" binding:"required"`
	State   MessageState `json:"state" binding:"required"`
}
