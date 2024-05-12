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
	ChatId  uint         `json:"chatId"`
	UserId  uint         `json:"userId"`
	Content string       `json:"content" binding:"required"`
	State   MessageState `json:"state" binding:"required"`
}
