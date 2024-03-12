package main

import (
	"log/slog"
	"os"

	"github.com/mitrkos/telemoney/internal/pkg/gsheetclient"
	"github.com/mitrkos/telemoney/internal/pkg/logger"
)

func main() {
	logger.SetLogger()

	authToken, ok := os.LookupEnv("TELEMONEY_GAUTH_TOKEN")
	if !ok {
		panic("Env var with tg token not found")
	}

	gsheetConfig := gsheetclient.Config{
		AuthToken:          authToken,
		SpreadsheetID:      "1DNP3yNOA03Qd52u6HPAw4uGQLSpQac2o5JaaI-9JjGs",
	}
	_, err := gsheetclient.New(&gsheetConfig)
	if err != nil {
		slog.Error("can't connect to gsheets", slog.Any("err", err))
		panic(err)
	}
}
