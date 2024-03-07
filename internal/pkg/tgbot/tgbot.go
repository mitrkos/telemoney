package tgbot

import (
	"log/slog"
	"strconv"

	"github.com/mitrkos/telemoney/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TOKEN_FILE = "./local/tg/token.txt"
)

type TgBot struct {
	client *tgbotapi.BotAPI

	updateHandlerMessage func(msg *model.Message) error
}

func New(token string) (*TgBot, error) {
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	slog.Info("Tg connected", slog.Any("botApi", botApi.Self))

	return &TgBot{client: botApi}, nil
}

func (tg *TgBot) SetDebug() {
	tg.client.Debug = true
}

func (tg *TgBot) SetUpdateHandlerMessage(updateHandlerMessage func(msg *model.Message) error) {
	tg.updateHandlerMessage = updateHandlerMessage
}

func (tg *TgBot) ListenToUpdates() {
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0) // TODO: how to make offset
	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := tg.client.GetUpdatesChan(updateConfig)
	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message != nil {
			err := tg.updateHandlerMessage(convertTgMessageToMessage(update.Message))
			if err != nil {
				slog.Error("Error while handling a tg msg", slog.Any("err", err), slog.Any("msg", update.Message))
			}
		}
	}
}


func convertTgMessageToMessage(tgMsg *tgbotapi.Message) *model.Message {
	return &model.Message{
		CreatedAt: tgMsg.Date,
		MessageId: strconv.Itoa(tgMsg.MessageID),
		Text: tgMsg.Text,
	}
}
