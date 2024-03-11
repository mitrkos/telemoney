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
		var tgMessage *telego.Message
		isEdited := false

		if update.Message != nil {
			tgMessage = update.Message
		}
		if update.EditedMessage != nil {
			tgMessage = update.EditedMessage
			isEdited = true
		}

		if tgMessage == nil {
			continue
		}
		
		err := tg.updateHandlerMessage(convertTgMessageToMessage(tgMessage, isEdited))
		if err == nil {
			tg.bot.SetMessageReaction(
				&telego.SetMessageReactionParams{
					ChatID: tgMessage.Chat.ChatID(),
					MessageID: tgMessage.MessageID,
					Reaction: setReactionSuccessEmoji(),
					IsBig: true,
				},
			)
		} else {
			slog.Error("Error while handling a tg msg", slog.Any("err", err), slog.Any("msg", tgMessage))
			tg.bot.SetMessageReaction(
				&telego.SetMessageReactionParams{
					ChatID: tgMessage.Chat.ChatID(),
					MessageID: tgMessage.MessageID,
					Reaction: setReactionUnknownMessageEmoji(),
					IsBig: true,
				},
			)
		}
	}

	return nil
}

func convertTgMessageToMessage(tgMsg *telego.Message, isEdited bool) *model.Message {
	return &model.Message{
		CreatedAt: tgMsg.Date,
		MessageId: strconv.Itoa(tgMsg.MessageID),
		IsEdited: isEdited,
		Text:      tgMsg.Text,
	}
}


func setReactionUnknownMessageEmoji() []telego.ReactionType {
	return []telego.ReactionType{&telego.ReactionTypeEmoji{
		Type: telego.ReactionEmoji,
		Emoji: "🤷‍♂",
	}}
}

func setReactionSuccessEmoji() []telego.ReactionType {
	return []telego.ReactionType{&telego.ReactionTypeEmoji{
		Type: telego.ReactionEmoji,
		Emoji: "👌",
	}}
}