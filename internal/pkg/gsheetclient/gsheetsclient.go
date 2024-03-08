package gsheetclient

import (
	"context"
	"encoding/base64"
	"log/slog"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GSheetsClient struct {
	config *Config

	service *sheets.Service
}

type Config struct {
	AuthToken string
	SpreadsheetID string
	TransactionSheetID string
}


func New(config *Config) (*GSheetsClient, error) {
	credBytes, err := base64.StdEncoding.DecodeString(config.AuthToken)
	if err != nil {
		return nil, err
	}

	// authenticate and get configuration
	jwtConfig, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// create client with config and context
	httpClient := jwtConfig.Client(ctx)

	// create new service using client
	service, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	slog.Info("gsheets connected", slog.Any("jwtConfig", jwtConfig), slog.Any("service", service))

	gsc := &GSheetsClient{config: config, service: service}
	return gsc, nil
}

func (gsc *GSheetsClient) WriteDataRow(dataRow []interface{}) {
	row := &sheets.ValueRange{
		Values: [][]interface{}{dataRow},
	}

	ctx := context.Background()
	response, err := gsc.service.Spreadsheets.Values.Append(gsc.config.SpreadsheetID, gsc.config.TransactionSheetID, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || response.HTTPStatusCode != 200 {
		slog.ErrorContext(ctx, "Write data failed", slog.Any("err", err), slog.Any("response", response), slog.Any("dataRow", dataRow))
		return
	}
}
