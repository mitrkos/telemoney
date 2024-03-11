package model

type Message struct {
	CreatedAt int64
	MessageId string
	IsEdited bool
	Text      string
}

type Transaction struct {
	CreatedAt int64
	MessageId string
	Amount    float64
	Category  string
	Tags      []string
	Comment   *string
}
