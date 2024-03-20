package apihandler

import (
	"errors"

	"github.com/mitrkos/telemoney/internal/model"
)

var ErrAPIOperationFailed = errors.New("api operation failed")

type MessageHandler interface {
	// inputs
	SetUpdateHandlerStartCommand(func())
	SetUpdateHandlerMessage(func(*model.MessageToHandle))
	SetUpdateHandlerEditedMessage(func(*model.MessageToHandle))
	SetUpdateHandlerRemoveMessageCommand(func(*model.MessageToHandle))

	// outputs
	SendMessage(*model.MessageToSend) error
	RemoveMessage(*model.MessageToInteract) error
	MarkMessageProcessedOK(*model.MessageToInteract) error
	MarkMessageProcessedFail(*model.MessageToInteract) error

	ListenToUpdates() error
}
