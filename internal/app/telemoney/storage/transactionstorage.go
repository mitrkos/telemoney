package storage

import (
	"errors"

	"github.com/mitrkos/telemoney/internal/model"
)

var ErrTransactionNotFound = errors.New("transaction not found")
var ErrOperationFailed = errors.New("operation failed")

type TransactionStorage interface {
	Insert(*model.Transaction) error
	Update(*model.Transaction) error
	DeleteByMessageId(string) error
}
