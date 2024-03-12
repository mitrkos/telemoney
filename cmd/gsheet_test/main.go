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
		TransactionSheetID: "transaction_test",
	}
	gSheetsClient, err := gsheetclient.New(&gsheetConfig)
	if err != nil {
		slog.Error("can't connect to gsheets", slog.Any("err", err))
		panic(err)
	}

	location, err := gSheetsClient.FindValueLocation(&gsheetclient.A1Range{
		SheetId: gsheetConfig.TransactionSheetID,
		LeftTop: gsheetclient.A1Location{
			Column: "B",
			Row: 1,
		},
		RightBottom: gsheetclient.A1Location{
			Column: "B",
		},
	}, "103")
	slog.Info("found location", slog.Any("location", location))
	if err != nil {
		panic(err)
	}


	newDataRow := make([]interface{}, 6)
	newDataRow[0] = 1234567
	newDataRow[1] = "103"
	newDataRow[2] = 123.5
	newDataRow[3] = "test_update"
	newDataRow[4] = "sazda, asd"
	newDataRow[5] = "Let's test"


	gSheetsClient.UpdateDataRange(&gsheetclient.A1Range{
		SheetId: gsheetConfig.TransactionSheetID,
		LeftTop: gsheetclient.A1Location{
			Column: "A",
			Row: location.Row,
		},
		RightBottom: gsheetclient.A1Location{
			Column: "F",
			Row: location.Row,
		},
	}, newDataRow)
	if err != nil {
		panic(err)
	}
}
