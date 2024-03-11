package gsheetclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

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
	AuthToken          string
	SpreadsheetID      string
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

func (gsc *GSheetsClient) AppendDataRow(dataRow []interface{}) {
	row := &sheets.ValueRange{
		Values: [][]interface{}{dataRow},
	}

	response, err := gsc.service.Spreadsheets.Values.Append(gsc.config.SpreadsheetID, gsc.config.TransactionSheetID, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Do()
	if err != nil || response.HTTPStatusCode != 200 {
		slog.Error("Append data to gseets failed", slog.Any("err", err), slog.Any("response", response), slog.Any("dataRow", dataRow))
		return
	}
}

func (gsc *GSheetsClient) EditDataRow(dataRow []interface{}, updateRange *A1Range) {
	row := &sheets.ValueRange{
		Values: [][]interface{}{dataRow},
	}

	response, err := gsc.service.Spreadsheets.Values.Update(gsc.config.SpreadsheetID, updateRange.String(), row).ValueInputOption("USER_ENTERED").Do()
	if err != nil || response.HTTPStatusCode != 200 {
		slog.Error("Update data to gseets failed", slog.Any("err", err), slog.Any("response", response), slog.Any("dataRow", dataRow), slog.Any("updateRange", updateRange))
		return
	}
}

func (gsc *GSheetsClient) FindValueLocation(searchValue string, searchRange *A1Range) (*A1Location, error) {
	response, err := gsc.service.Spreadsheets.Values.Get(gsc.config.SpreadsheetID, searchRange.String()).Do()
	if err != nil {
		return nil, err
	}
	respJson, err := json.Marshal(response)
	slog.Info("gsheet get response", slog.Any("respJson", string(respJson)), slog.Any("err", err))

	searchValueColumnIdx := -1
	searchValueRowIdx := -1

	for rowIdx, row := range response.Values {
		for columnIdx, valueRaw := range row {
			if value, ok := valueRaw.(string); ok && value == searchValue {
				searchValueColumnIdx = columnIdx
				searchValueRowIdx = rowIdx
				break
			}
		}
		if searchValueColumnIdx != -1 && searchValueRowIdx != -1 {
			break
		}
	}

	if searchValueColumnIdx == -1 || searchValueRowIdx == -1 {
		return nil, nil
	}

	searchValueColumnA1Index := toIntAlphabetic(searchRange.LeftTop.Column) + searchValueColumnIdx
	searchValueRowA1Index := searchRange.LeftTop.Row + searchValueRowIdx

	return &A1Location{
		Column: toStrAlphabetic(searchValueColumnA1Index),
		Row:    searchValueRowA1Index,
	}, nil
}


func (gsc *GSheetsClient) ClearRange(deleteRange *A1Range) error {
	response, err := gsc.service.Spreadsheets.Values.Clear(gsc.config.SpreadsheetID, deleteRange.String(), &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		return err
	}
	respJson, err := json.Marshal(response)
	slog.Info("gsheet clear response", slog.Any("respJson", string(respJson)), slog.Any("err", err))

	return nil
}

type A1Location struct {
	Column string // A, B, C, ...
	Row    int    // (0 - not set) 1, 2, 3 ...
}

type A1Range struct {
	SheetId     string
	LeftTop     A1Location
	RightBottom A1Location
}

func (l *A1Location) String() string {
	result := l.Column
	if l.Row != 0 {
		result += fmt.Sprint(l.Row)
	}
	return result
}

func (r *A1Range) String() string {
	result := r.SheetId + "!"
	result += r.LeftTop.String()
	result += ":"
	result += r.RightBottom.String()
	return result
}

const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func toStrAlphabetic(i int) string {
	return abc[i-1 : i]
}

func toIntAlphabetic(symbol string) int {
	return strings.Index(abc, symbol) + 1
}

