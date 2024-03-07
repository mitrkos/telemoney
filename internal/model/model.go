package model

type Message struct {
	CreatedAt int
	MessageId string
	Text      string
}

type Transaction struct {
	CreatedAt int
	MessageId string
	Amount    float64
	Category  string
	Tags      []string
	Comment   *string
}
