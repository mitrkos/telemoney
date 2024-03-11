package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/mitrkos/telemoney/internal/model"
	"github.com/mitrkos/telemoney/internal/pkg/tgbot"
)

func main() {
	tgToken, ok := os.LookupEnv("TELEMONEY_TG_BOT_TOKEN_TEST")
	if !ok {
		panic("Env var with tg token not found")
	}

	tgConfig := tgbot.Config{
		AuthToken: tgToken,
	}

	tgBot, err := tgbot.New(&tgConfig)
	if err != nil {
		slog.Error("can't connect to tg", slog.Any("err", err))
		panic(err)
	}

	tgBot.SetUpdateHandlerMessage(func(msg *model.Message) error {
		if msg.Text == "test error" {
			return errors.New("test error")
		}
		return nil
	})

	if err := tgBot.ListenToUpdates(); err != nil {
		panic(err)
	}
}
