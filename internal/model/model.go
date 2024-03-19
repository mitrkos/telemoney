package model

type MessageToHandle struct {
	CreatedAt int64
	MessageID string
	ChatID    string
	Text      string
}

type MessageToSend struct {
	Text   string
	ChatID string
}

type MessageToInteract struct {
	MessageID string
	ChatID    string
}

type Transaction struct {
	CreatedAt int64
	MessageID string
	Amount    float64
	Category  string
	Tags      []string
	Comment   *string
}
