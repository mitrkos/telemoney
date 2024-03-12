package telemoney

import (
	"log/slog"

	"github.com/mitrkos/telemoney/internal/app/telemoney/storage"
	"github.com/mitrkos/telemoney/internal/app/telemoney/storage/gsheetstorage"
	"github.com/mitrkos/telemoney/internal/model"
	"github.com/mitrkos/telemoney/internal/pkg/gsheetclient"
	parsing "github.com/mitrkos/telemoney/internal/pkg/parser"
	"github.com/mitrkos/telemoney/internal/pkg/tgbot"
)

func Start() error {
	config, err := readConfig()
	if err != nil {
		slog.Error("can't read the config", slog.Any("err", err))
		return err
	}

	tgConfig := tgbot.Config{
		AuthToken: config.TgAuthTokenTest,
	}
	if config.Env == "prod" {
		tgConfig.AuthToken = config.TgAuthToken
	}

	tgBot, err := tgbot.New(&tgConfig)
	if err != nil {
		slog.Error("can't connect to tg", slog.Any("err", err))
		return err
	}

	gsheetConfig := gsheetclient.Config{
		AuthToken:     config.GSheetsAuthToken,
		SpreadsheetID: config.SpreadsheetID,
	}
	gSheetsClient, err := gsheetclient.New(&gsheetConfig)
	if err != nil {
		slog.Error("can't connect to gsheets", slog.Any("err", err))
		return err
	}

	transactionSheetID := config.TransactionSheetIDTest
	if config.Env == "prod" {
		transactionSheetID = config.TransactionSheetID
	}
	transactionStorage := gsheetstorage.New(gSheetsClient, transactionSheetID)

	parser := parsing.New()

	tgBot.SetUpdateHandlerMessage(makeHandleTgMessage(parser, transactionStorage))
	err = tgBot.ListenToUpdates()
	if err != nil {
		slog.Error("problem with listening to tg", slog.Any("err", err))
		return err
	}

	return nil
}

func makeHandleTgMessage(parser *parsing.Parser, storage storage.TransactionRepository) func(msg *model.Message) error {
	return func(msg *model.Message) error {
		transaction, err := convertMessageIntoTransaction(parser, msg)
		if err != nil {
			return err
		}

		if transaction != nil {
			if msg.IsEdited {
				storage.Update(transaction)
			} else {
				storage.Insert(transaction)
			}
		}

		return nil
	}
}

func convertMessageIntoTransaction(parser *parsing.Parser, msg *model.Message) (*model.Transaction, error) {
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
