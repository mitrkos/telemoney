package gsheetclient

import (
	"context"
	"encoding/base64"
	"log/slog"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	TOKEN_FILE = "./local/gauth/telemoney-b63c1d5ddf79.txt"

	// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>
	// https://docs.google.com/spreadsheets/d/1DNP3yNOA03Qd52u6HPAw4uGQLSpQac2o5JaaI-9JjGs/edit#gid=0
	SPREADSHEET_ID = "1DNP3yNOA03Qd52u6HPAw4uGQLSpQac2o5JaaI-9JjGs"
	SHEET_ID       = 0
	// TRANSACTION_SPREADSHEET_NAME = "transaction"
	TRANSACTION_SPREADSHEET_NAME = "transaction_test"
)

type GSheetsClient struct {
	service *sheets.Service
}

func New(token string) (*GSheetsClient, error) {
	credBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// create client with config and context
	httpClient := config.Client(ctx)

	// create new service using client
	service, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	slog.Info("gsheets connected", slog.Any("config", config), slog.Any("service", service))

	gsc := &GSheetsClient{service: service}
	return gsc, nil
}

func (gsc *GSheetsClient) WriteDataRow(dataRow []interface{}) {
	row := &sheets.ValueRange{
		Values: [][]interface{}{dataRow},
	}

	ctx := context.Background()
	response, err := gsc.service.Spreadsheets.Values.Append(SPREADSHEET_ID, TRANSACTION_SPREADSHEET_NAME, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || response.HTTPStatusCode != 200 {
		slog.ErrorContext(ctx, "Write data failed", slog.Any("err", err), slog.Any("response", response), slog.Any("dataRow", dataRow))
		return
	}
}
