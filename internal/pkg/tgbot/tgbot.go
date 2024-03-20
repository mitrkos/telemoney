package tgbot

import (
	"log/slog"
	"strconv"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"

	"github.com/mitrkos/telemoney/internal/model"
)

type TgMessageReaction = []telego.ReactionType

type ReactionForMessage struct {
	Msg      *model.MessageToInteract
	Reaction TgMessageReaction
}

type TgBot struct {
	config *Config

	bot *telego.Bot

	updateHandlerStartCommand         func()
	updateHandlerRemoveMessageCommand func(msg *model.MessageToHandle)
	updateHandlerMessage              func(msg *model.MessageToHandle)
	updateHandlerEditedMessage        func(msg *model.MessageToHandle)
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
		bot:    bot,
	}, nil
}

func (tg *TgBot) SetUpdateHandlerStartCommand(handler func()) {
	tg.updateHandlerStartCommand = handler
}

func (tg *TgBot) SetUpdateHandlerRemoveMessageCommand(handler func(*model.MessageToHandle)) {
	tg.updateHandlerRemoveMessageCommand = handler
}

func (tg *TgBot) SetUpdateHandlerMessage(handler func(*model.MessageToHandle)) {
	tg.updateHandlerMessage = handler
}

func (tg *TgBot) SetUpdateHandlerEditedMessage(handler func(*model.MessageToHandle)) {
	tg.updateHandlerEditedMessage = handler
}

func (tg *TgBot) ListenToUpdates() error {
	// Get updates channel
	// (more on configuration in examples/updates_long_polling/main.go)
	updates, err := tg.bot.UpdatesViaLongPolling(nil)
	if err != nil {
		return err
	}

	handler, _ := telegohandler.NewBotHandler(tg.bot, updates)

	// Stop handling updates
	defer handler.Stop()
	// Stop reviving updates from update channel
	defer tg.bot.StopLongPolling()

	handler.Handle(func(_ *telego.Bot, _ telego.Update) {
		if tg.updateHandlerStartCommand == nil {
			return
		}

		tg.updateHandlerStartCommand()
	}, telegohandler.CommandEqual("start"))

	handler.Handle(func(_ *telego.Bot, update telego.Update) {
		if tg.updateHandlerRemoveMessageCommand == nil {
			return
		}

		tg.updateHandlerRemoveMessageCommand(convertTGMessageToMessage(update.Message.ReplyToMessage))
	}, telegohandler.CommandEqual("remove"))

	handler.Handle(func(_ *telego.Bot, update telego.Update) {
		if tg.updateHandlerEditedMessage == nil {
			return
		}
		if update.EditedMessage == nil {
			return
		}

		tg.updateHandlerMessage(convertTGMessageToMessage(update.EditedMessage))
	}, telegohandler.AnyEditedMessageWithText())

	handler.Handle(func(_ *telego.Bot, update telego.Update) {
		if tg.updateHandlerEditedMessage == nil {
			return
		}
		if update.Message == nil {
			return
		}

		tg.updateHandlerMessage(convertTGMessageToMessage(update.Message))
	}, telegohandler.AnyMessageWithText())

	// Start handling updates
	handler.Start()

	return nil
}

func (tg *TgBot) SendMessage(msg *model.MessageToSend) error {
	tgChatID, err := convertChatIDToTGChatID(msg.ChatID)
	if err != nil {
		return err
	}

	_, err = tg.bot.SendMessage(telegoutil.Messagef(telegoutil.ID(tgChatID), msg.Text))
	if err != nil {
		slog.Error("sending msg to tg failed", slog.Any("err", err), slog.Any("msg", msg))
		return err
	}

	return nil
}

func (tg *TgBot) RemoveMessage(msg *model.MessageToInteract) error {
	tgChatID, err := convertChatIDToTGChatID(msg.ChatID)
	if err != nil {
		return err
	}
	tgMessageID, err := convertMessageIDToTGMessageID(msg.MessageID)
	if err != nil {
		return err
	}

	err = tg.bot.DeleteMessage(&telego.DeleteMessageParams{
		ChatID:    telegoutil.ID(tgChatID),
		MessageID: tgMessageID,
	})
	if err != nil {
		slog.Error("deleting msg in tg  failed", slog.Any("err", err), slog.Any("msg", msg))
		return err
	}

	return nil
}

func (tg *TgBot) SetMessageReaction(reactionForMsg *ReactionForMessage) error {
	tgChatID, err := convertChatIDToTGChatID(reactionForMsg.Msg.ChatID)
	if err != nil {
		return err
	}
	tgMessageID, err := convertMessageIDToTGMessageID(reactionForMsg.Msg.MessageID)
	if err != nil {
		return err
	}

	err = tg.bot.SetMessageReaction(
		&telego.SetMessageReactionParams{
			ChatID:    telegoutil.ID(tgChatID),
			MessageID: tgMessageID,
			Reaction:  reactionForMsg.Reaction,
			IsBig:     true,
		},
	)
	if err != nil {
		slog.Error("setting msg reaction failed", slog.Any("err", err), slog.Any("reactionForMsg", reactionForMsg))
		return err
	}

	return nil
}

func convertTGMessageToMessage(tgMsg *telego.Message) *model.MessageToHandle {
	if tgMsg == nil {
		return nil
	}

	return &model.MessageToHandle{
		CreatedAt: tgMsg.Date,
		MessageID: strconv.Itoa(tgMsg.MessageID),
		ChatID:    strconv.FormatInt(tgMsg.Chat.ID, 10),
		Text:      tgMsg.Text,
	}
}

func convertChatIDToTGChatID(chatID string) (int64, error) {
	tgChatID, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		slog.Error("can't convert ChatID to tg ChatID", slog.Any("chatID", chatID))
	}

	return tgChatID, err
}

func convertMessageIDToTGMessageID(messageID string) (int, error) {
	tgMessageID, err := strconv.Atoi(messageID)
	if err != nil {
		slog.Error("can't convert MessageID to tg MessageID", slog.Any("messageID", tgMessageID))
	}

	return tgMessageID, err
}

func MakeReactionShruggingEmoji() TgMessageReaction {
	return []telego.ReactionType{&telego.ReactionTypeEmoji{
		Type:  telego.ReactionEmoji,
		Emoji: "🤷‍♂",
	}}
}

func MakeReactionOkEmoji() TgMessageReaction {
	return []telego.ReactionType{&telego.ReactionTypeEmoji{
		Type:  telego.ReactionEmoji,
		Emoji: "👌",
	}}
}
