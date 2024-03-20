package gsheetclient

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
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
	AuthToken     string
	SpreadsheetID string
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

func (gsc *GSheetsClient) AppendDataToRange(appendRange *A1Range, dataRow []interface{}) error {
	row := &sheets.ValueRange{
		Values: [][]interface{}{dataRow},
	}

	response, err := gsc.service.Spreadsheets.Values.
		Append(gsc.config.SpreadsheetID, appendRange.String(), row).
		ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").
		Do()
	if err = parseGSheetAPIError(err, response.HTTPStatusCode); err != nil {
		slog.Error(
			"Append data to gseets failed",
			slog.Any("err", err),
			slog.Any("response", response),
			slog.Any("dataRow", dataRow),
		)
		return err
	}
	return nil
}

func (gsc *GSheetsClient) UpdateDataRange(updateRange *A1Range, dataRow []interface{}) error {
	row := &sheets.ValueRange{
		Values: [][]interface{}{dataRow},
	}

	response, err := gsc.service.Spreadsheets.Values.
		Update(gsc.config.SpreadsheetID, updateRange.String(), row).
		ValueInputOption("USER_ENTERED").
		Do()
	if err = parseGSheetAPIError(err, response.HTTPStatusCode); err != nil {
		slog.Error(
			"Update data to gseets failed",
			slog.Any("err", err), slog.Any("response", response), slog.Any("dataRow", dataRow), slog.Any("updateRange", updateRange))
		return err
	}
	return nil
}

func (gsc *GSheetsClient) ClearRange(deleteRange *A1Range) error {
	response, err := gsc.service.Spreadsheets.Values.Clear(gsc.config.SpreadsheetID, deleteRange.String(), &sheets.ClearValuesRequest{}).Do()
	if err = parseGSheetAPIError(err, response.HTTPStatusCode); err != nil {
		slog.Error("Clear data in gseets failed", slog.Any("err", err), slog.Any("response", response), slog.Any("deleteRange", deleteRange))
		return err
	}
	return nil
}

func (gsc *GSheetsClient) FindValueLocation(searchRange *A1Range, searchValue string) (*A1Location, error) {
	response, err := gsc.service.Spreadsheets.Values.Get(gsc.config.SpreadsheetID, searchRange.String()).Do()
	if err = parseGSheetAPIError(err, response.HTTPStatusCode); err != nil {
		slog.Error("Find data in gseets failed", slog.Any("err", err), slog.Any("response", response), slog.Any("searchRange", searchRange))
		return nil, err
	}

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
		return nil, nil //nolint:nilnil // fix later
	}

	searchValueColumnA1Index := toIntAlphabetic(searchRange.LeftTop.Column) + searchValueColumnIdx
	searchValueRowA1Index := searchRange.LeftTop.Row + searchValueRowIdx

	return &A1Location{
		Column: toStrAlphabetic(searchValueColumnA1Index),
		Row:    searchValueRowA1Index,
	}, nil
}

func parseGSheetAPIError(err error, httpStatusCode int) error {
	if err != nil {
		return err
	}
	if httpStatusCode != http.StatusOK {
		return fmt.Errorf("gsheet connection error: %d", httpStatusCode)
	}
	return nil
}

type A1Location struct {
	Column string // A, B, C, ...
	Row    int    // (0 - not set) 1, 2, 3 ...
}

type A1Range struct {
	SheetID     string
	LeftTop     *A1Location
	RightBottom *A1Location
}

func (l *A1Location) String() string {
	result := l.Column
	if l.Row != 0 {
		result += strconv.Itoa(l.Row)
	}
	return result
}

func (r *A1Range) String() string {
	result := r.SheetID
	if r.LeftTop != nil {
		result += "!" + r.LeftTop.String()

		if r.RightBottom != nil {
			result += ":" + r.RightBottom.String()
		}
	}
	return result
}

const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func toStrAlphabetic(i int) string {
	return abc[i-1 : i]
}

func toIntAlphabetic(symbol string) int {
	return strings.Index(abc, symbol) + 1
}
