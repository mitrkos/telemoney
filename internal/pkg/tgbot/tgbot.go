package tgbot

import (
	"log/slog"
	"strconv"

	"github.com/mitrkos/telemoney/internal/model"
	"github.com/mymmrac/telego"
)

type TgBot struct {
	config *Config

	bot *telego.Bot

	updateHandlerMessage func(msg *model.Message) error
}

type Config struct {
	AuthToken string
}

func New(config *Config) (*TgBot, error) {
	bot, err := telego.NewBot(config.AuthToken, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}

	return &TgBot{
		config: config,
		bot: bot,
	}, nil
}

func (tg *TgBot) SetUpdateHandlerMessage(updateHandlerMessage func(msg *model.Message) error) {
	tg.updateHandlerMessage = updateHandlerMessage
}

func (tg *TgBot) ListenToUpdates() error {
	// Get updates channel
	// (more on configuration in examples/updates_long_polling/main.go)
	updates, err := tg.bot.UpdatesViaLongPolling(nil)
	if err != nil {
		return err
	}

	// Stop reviving updates from update channel
	defer tg.bot.StopLongPolling()

	// Loop through all updates when they came
	for update := range updates {
		if update.Message != nil {
			err := tg.updateHandlerMessage(convertTgMessageToMessage(update.Message))
			if err != nil {
				slog.Error("Error while handling a tg msg", slog.Any("err", err), slog.Any("msg", update.Message))
			}	
		}
	}

	return nil
}

func convertTgMessageToMessage(tgMsg *telego.Message) *model.Message {
	return &model.Message{
		CreatedAt: tgMsg.Date,
		MessageId: strconv.Itoa(tgMsg.MessageID),
		Text:      tgMsg.Text,
	}
}
