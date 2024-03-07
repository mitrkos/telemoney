package telemoney

import (
	parsing "github.com/mitrkos/telemoney/internal/pkg/parser"
	"github.com/mitrkos/telemoney/internal/model"

	"github.com/mitrkos/telemoney/internal/pkg/gsheetclient"
	"github.com/mitrkos/telemoney/internal/pkg/tgbot"
	"github.com/mitrkos/telemoney/internal/utils"
)

type Config struct {
	gsheets_token string
	tg_token      string
}

func Start() error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	tgBot, err := tgbot.New(config.tg_token)
	if err != nil {
		return err
	}
	tgBot.SetDebug()

	gSheetsClient, err := gsheetclient.New(config.gsheets_token)
	if err != nil {
		return err
	}

	parser := parsing.New()

	tgBot.SetUpdateHandlerMessage(makeHandleTgMessage(parser, gSheetsClient))
	tgBot.ListenToUpdates()

	return nil
}

func getConfig() (*Config, error) {
	gsheetToken, err := utils.GetTokenFromFile(gsheetclient.TOKEN_FILE)
	if err != nil {
		return nil, err
	}

	tgToken, err := utils.GetTokenFromFile(tgbot.TOKEN_FILE)
	if err != nil {
		return nil, err
	}

	return &Config{
		gsheets_token: gsheetToken,
		tg_token:      tgToken,
	}, nil
}

func makeHandleTgMessage(parser *parsing.Parser, gSheetsClient *gsheetclient.GSheetsClient) func(msg *model.Message) error {
	return func(msg *model.Message) error {
		transaction, err := convertMessageIntoTransaction(parser, msg)
		if err != nil {
			return err
		}

		if transaction != nil {
			gSheetsClient.WriteRow([]interface{}{transaction.CreatedAt, transaction.MessageId, transaction.Amount, transaction.Category, transaction.Tags, transaction.Comment})
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
		Amount: userInputData.Amount,
		Category: userInputData.Category,
		Tags: userInputData.Tags,
		Comment: userInputData.Comment,
	}, nil
}
