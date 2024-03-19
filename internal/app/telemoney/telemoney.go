package telemoney

import (
	"errors"
	"log/slog"

	"github.com/mitrkos/telemoney/internal/app/telemoney/apihandler"
	"github.com/mitrkos/telemoney/internal/app/telemoney/storage"
	"github.com/mitrkos/telemoney/internal/model"
	parsing "github.com/mitrkos/telemoney/internal/pkg/parser"
)

type Telemoney struct {
	config             *Config
	api                apihandler.MessageHandler
	transactionStorage storage.TransactionStorage
	parser             *parsing.Parser
}

func New(config *Config, api apihandler.MessageHandler, storage storage.TransactionStorage, parser *parsing.Parser) *Telemoney {
	t := Telemoney{
		config:             config,
		api:                api,
		transactionStorage: storage,
		parser:             parser,
	}

	t.api.SetUpdateHandlerStartCommand(t.handleStartCommand)
	t.api.SetUpdateHandlerRemoveMessageCommand(t.handleRemoveMessageCommand)
	t.api.SetUpdateHandlerEditedMessage(t.handleEditedMessage)
	t.api.SetUpdateHandlerMessage(t.handleMessage)

	return &t
}

func (t *Telemoney) Start() error {
	err := t.api.ListenToUpdates()

	if err != nil {
		slog.Error("problem with listening to tg", slog.Any("err", err))
		return err
	}
	return nil
}

func (t *Telemoney) handleStartCommand() {
	// TODO
}

func (t *Telemoney) handleRemoveMessageCommand(msg *model.MessageToHandle) {
	if msg == nil {
		// TODO: info msg
		// t.api.SendMessage(&model.MessageToSend{
		// 	ChatID: msg.ChatID,
		// 	Text:   "/remove command should be a reply to a message",
		// })
		// t.markMessageHandledFailure(msg)
		return
	}

	err := t.transactionStorage.DeleteByMessageId(msg.MessageID)
	if err != nil {
		return
	}

	_ = t.api.RemoveMessage(&model.MessageToInteract{
		ChatID:    msg.ChatID,
		MessageID: msg.MessageID,
	})
}

func (t *Telemoney) handleEditedMessage(msg *model.MessageToHandle) {
	transaction, err := convertMessageIntoTransaction(t.parser, msg)
	if err != nil {
		t.markMessageHandledFailure(msg)
		return
	}

	err = t.transactionStorage.Update(transaction)
	if errors.Is(err, storage.ErrTransactionNotFound) {
		err = t.transactionStorage.Insert(transaction)
	}

	if err != nil {
		t.markMessageHandledFailure(msg)
		return
	}

	t.markMessageHandleSuccess(msg)
}

func (t *Telemoney) handleMessage(msg *model.MessageToHandle) {
	transaction, err := convertMessageIntoTransaction(t.parser, msg)
	if err != nil {
		t.markMessageHandledFailure(msg)
		return
	}

	err = t.transactionStorage.Insert(transaction)
	if err != nil {
		t.markMessageHandledFailure(msg)
		return
	}

	t.markMessageHandleSuccess(msg)
}

func (t *Telemoney) markMessageHandleSuccess(msg *model.MessageToHandle) {
	_ = t.api.MarkMessageProcessedOK(&model.MessageToInteract{
		ChatID:    msg.ChatID,
		MessageID: msg.MessageID,
	})
}

func (t *Telemoney) markMessageHandledFailure(msg *model.MessageToHandle) {
	_ = t.api.MarkMessageProcessedFail(&model.MessageToInteract{
		ChatID:    msg.ChatID,
		MessageID: msg.MessageID,
	})
}

func convertMessageIntoTransaction(parser *parsing.Parser, msg *model.MessageToHandle) (*model.Transaction, error) {
	if msg == nil {
		return nil, nil
	}

	userInputData, err := parser.ParseTransactionUserInputDataFromText(msg.Text)
	if err != nil || userInputData == nil {
		return nil, err
	}

	return &model.Transaction{
		CreatedAt: msg.CreatedAt,
		MessageID: msg.MessageID,
		Amount:    userInputData.Amount,
		Category:  userInputData.Category,
		Tags:      userInputData.Tags,
		Comment:   userInputData.Comment,
	}, nil
}
