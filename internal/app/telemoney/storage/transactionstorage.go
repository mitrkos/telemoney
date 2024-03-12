package storage

import "github.com/mitrkos/telemoney/internal/model"


type TransactionStorage interface {
	Insert(*model.Transaction) error
	Update(*model.Transaction) error
	DeleteByMessageId(string) error
}
