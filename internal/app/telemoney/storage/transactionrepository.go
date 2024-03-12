package storage

import "github.com/mitrkos/telemoney/internal/model"


type TransactionRepository interface {
	Insert(*model.Transaction) error
	Update(*model.Transaction) error
	DeleteByMessageId(string) error
}
