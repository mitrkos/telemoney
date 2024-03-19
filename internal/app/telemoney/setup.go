package telemoney

import (
	"log/slog"

	"github.com/mitrkos/telemoney/internal/app/telemoney/apihandler"
	"github.com/mitrkos/telemoney/internal/app/telemoney/apihandler/tgbothandler"
	"github.com/mitrkos/telemoney/internal/app/telemoney/storage"
	"github.com/mitrkos/telemoney/internal/app/telemoney/storage/gsheetstorage"
	"github.com/mitrkos/telemoney/internal/pkg/gsheetclient"
	parsing "github.com/mitrkos/telemoney/internal/pkg/parser"
	"github.com/mitrkos/telemoney/internal/pkg/tgbot"
)

type TelemoneyDependencies struct {
	Config             *Config
	Api                apihandler.MessageHandler
	TransactionStorage storage.TransactionStorage
	Parser             *parsing.Parser
}

func PrepareDependencies() (*TelemoneyDependencies, error) {
	config, err := readConfig()
	if err != nil {
		slog.Error("can't read the config", slog.Any("err", err))
		return nil, err
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
		return nil, err
	}
	tgBotHandler := tgbothandler.New(tgBot)

	gsheetConfig := gsheetclient.Config{
		AuthToken:     config.GSheetsAuthToken,
		SpreadsheetID: config.SpreadsheetID,
	}
	gSheetsClient, err := gsheetclient.New(&gsheetConfig)
	if err != nil {
		slog.Error("can't connect to gsheets", slog.Any("err", err))
		return nil, err
	}

	transactionSheetID := config.TransactionSheetIDTest
	if config.Env == "prod" {
		transactionSheetID = config.TransactionSheetID
	}
	transactionStorage := gsheetstorage.New(gSheetsClient, transactionSheetID)

	parser := parsing.New()

	return &TelemoneyDependencies{
		Config:             config,
		Api:                tgBotHandler,
		TransactionStorage: transactionStorage,
		Parser:             parser,
	}, nil
}
