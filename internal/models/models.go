package models

type Message struct {
	createdAt int64
	messageId string
	text string
}

type Transaction struct {
	createdAt int
	messageId string
	amount float32
	category string
	tags []string
	comment string
}
