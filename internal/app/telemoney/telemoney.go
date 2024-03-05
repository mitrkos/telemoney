package telemoney

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

	tgBot.SetUpdateHandlerMessage(makeHandleTgMessage(gSheetsClient))
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

func makeHandleTgMessage(gSheetsClient *gsheetclient.GSheetsClient) func(tgMessage *tgbotapi.Message) error {
	return func(tgMessage *tgbotapi.Message) error {
		gSheetsClient.WriteData(tgMessage.Text)
		return nil
	}
}
