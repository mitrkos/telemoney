package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mitrkos/telemoney/gsheetclient"
	"github.com/mitrkos/telemoney/tgbot"
)

func main() {

	tgBotToken, err := tgbot.GetToken()
	if err != nil {
		panic(err)
	}
	tgBot, err := tgbot.Create(tgBotToken)
	if err != nil {
		panic(err)
	}
	tgBot.SetDebug()

	gSheetsClientToken, err := gsheetclient.GetToken()
	if err != nil {
		panic(err)
	}
	gSheetsClient, err := gsheetclient.Create(gSheetsClientToken)
	if err != nil {
		panic(err)
	}

	tgBot.SetUpdateHandlerMessage(makeHandleTgMessage(gSheetsClient))
	tgBot.ListenToUpdates()
}

func makeHandleTgMessage(gSheetsClient *gsheetclient.GSheetsClient) func(tgMessage *tgbotapi.Message) error {
	return func(tgMessage *tgbotapi.Message) error {
		gSheetsClient.WriteData(tgMessage.Text)
		return nil
	}
}
