package telemoney

import (
	"log/slog"
	"strings"

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
		AuthToken:          config.GSheetsAuthToken,
		SpreadsheetID:      config.SpreadsheetID,
		TransactionSheetID: config.TransactionSheetIDTest,
	}
	if config.Env == "prod" {
		gsheetConfig.TransactionSheetID = config.TransactionSheetID
	}
	gSheetsClient, err := gsheetclient.New(&gsheetConfig)
	if err != nil {
		slog.Error("can't connect to gsheets", slog.Any("err", err))
		return err
	}

	parser := parsing.New()

	tgBot.SetUpdateHandlerMessage(makeHandleTgMessage(parser, gSheetsClient))
	err = tgBot.ListenToUpdates()
	if err != nil {
		slog.Error("problem with listening to tg", slog.Any("err", err))
		return err
	}

	return nil
}

func makeHandleTgMessage(parser *parsing.Parser, gSheetsClient *gsheetclient.GSheetsClient) func(msg *model.Message) error {
	return func(msg *model.Message) error {
		if msg.IsEdited {
			// skip edited for now
			return nil
		}

		transaction, err := convertMessageIntoTransaction(parser, msg)
		if err != nil {
			return err
		}

		if transaction != nil {
			gSheetsClient.AppendDataRow(convertTransactionToDataRow(transaction))
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
		MessageId: msg.MessageId,
		Amount:    userInputData.Amount,
		Category:  userInputData.Category,
		Tags:      userInputData.Tags,
		Comment:   userInputData.Comment,
	}, nil
}

func convertTransactionToDataRow(transaction *model.Transaction) []interface{} {
	dataRow := make([]interface{}, 6)

	dataRow[0] = transaction.CreatedAt
	dataRow[1] = transaction.MessageId
	dataRow[2] = transaction.Amount
	dataRow[3] = transaction.Category
	if len(transaction.Tags) > 0 {
		tagsStr := strings.Join(transaction.Tags[:], ",")
		dataRow[4] = tagsStr
	}
	if transaction.Comment != nil {
		dataRow[5] = *transaction.Comment
	}

	return dataRow
}
