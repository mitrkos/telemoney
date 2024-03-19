package tgbothandler

import (
	"github.com/mitrkos/telemoney/internal/app/telemoney/apihandler"
	"github.com/mitrkos/telemoney/internal/model"
	"github.com/mitrkos/telemoney/internal/pkg/tgbot"
)

type TgBotMessageHandler struct {
	tgbot *tgbot.TgBot
}

func New(tgbot *tgbot.TgBot) *TgBotMessageHandler {
	return &TgBotMessageHandler{
		tgbot: tgbot,
	}
}

func (tgh *TgBotMessageHandler) ListenToUpdates() error {
	return tgh.tgbot.ListenToUpdatesUsingHandlers()
}

func (tgh *TgBotMessageHandler) SetUpdateHandlerStartCommand(handler func()) {
	tgh.tgbot.SetUpdateHandlerStartCommand(handler)
}

func (tgh *TgBotMessageHandler) SetUpdateHandlerRemoveMessageCommand(handler func(*model.MessageToHandle)) {
	tgh.tgbot.SetUpdateHandlerRemoveMessageCommand(handler)
}

func (tgh *TgBotMessageHandler) SetUpdateHandlerMessage(handler func(*model.MessageToHandle)) {
	tgh.tgbot.SetUpdateHandlerMessage(handler)
}

func (tgh *TgBotMessageHandler) SetUpdateHandlerEditedMessage(handler func(*model.MessageToHandle)) {
	tgh.tgbot.SetUpdateHandlerEditedMessage(handler)
}

func (tgh *TgBotMessageHandler) SendMessage(msg *model.MessageToSend) error {
	err := tgh.tgbot.SendMessage(msg)
	if err != nil {
		return apihandler.ErrApiOperationFailed
	}
	return nil
}

func (tgh *TgBotMessageHandler) RemoveMessage(msg *model.MessageToInteract) error {
	err := tgh.tgbot.RemoveMessage(msg)
	if err != nil {
		return apihandler.ErrApiOperationFailed
	}
	return nil
}

func (tgh *TgBotMessageHandler) MarkMessageProcessedOK(msg *model.MessageToInteract) error {
	err := tgh.tgbot.SetMessageReaction(&tgbot.ReactionForMessage{
		Msg:      msg,
		Reaction: tgbot.MakeReactionOkEmoji(),
	})
	if err != nil {
		return apihandler.ErrApiOperationFailed
	}
	return nil
}

func (tgh *TgBotMessageHandler) MarkMessageProcessedFail(msg *model.MessageToInteract) error {
	err := tgh.tgbot.SetMessageReaction(&tgbot.ReactionForMessage{
		Msg:      msg,
		Reaction: tgbot.MakeReactionShruggingEmoji(),
	})
	if err != nil {
		return apihandler.ErrApiOperationFailed
	}
	return nil
}
