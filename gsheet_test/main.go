package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	// https://docs.google.com/spreadsheets/d/<SPREADSHEETID>/edit#gid=<SHEETID>
	// https://docs.google.com/spreadsheets/d/1DNP3yNOA03Qd52u6HPAw4uGQLSpQac2o5JaaI-9JjGs/edit#gid=0
	SPREADSHEET_ID = "1DNP3yNOA03Qd52u6HPAw4uGQLSpQac2o5JaaI-9JjGs"
	SHEET_ID       = 0
)

func main() {
	auth64, err := os.ReadFile("./local/gauth/telemoney-b63c1d5ddf79.txt")
	if err != nil {
		panic(err)
	}

	// create api context
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	// credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))

	credBytes, err := base64.StdEncoding.DecodeString(string(auth64))
	if err != nil {
		fmt.Print(err)
		return
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(config.Email)
	fmt.Println(config.Scopes)

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(srv)

	// Convert sheet ID to sheet name.
	response1, err := srv.Spreadsheets.Get(SPREADSHEET_ID).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || response1.HTTPStatusCode != 200 {
		fmt.Println(config.Scopes)
		fmt.Print(err)
		return
	}

	sheetName := ""
	for _, v := range response1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(SHEET_ID) {
			sheetName = prop.Title
			break
		}
	}

	//Append value to the sheet.
	row := &sheets.ValueRange{
		Values: [][]interface{}{{"1", "ABC", "abc@gmail.com"}},
	}

	response2, err := srv.Spreadsheets.Values.Append(SPREADSHEET_ID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	if err != nil || response2.HTTPStatusCode != 200 {
		fmt.Print(err)
		return
	}

	// The A1 notation of cells range to update.
	range2 := "A1:C1"

	// prepare data for update cells
	row = &sheets.ValueRange{
		Values: [][]interface{}{{"2", "XYZ", "xyz@gmail.com"}},
	}

	// update cells in given range
	_, err = srv.Spreadsheets.Values.Update(SPREADSHEET_ID, range2, row).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		fmt.Print(err)
		return
	}

	records := []string{"1", "ABC", "abc@gmail.com"}

	// create the batch request
	batchUpdateRequest := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				AppendCells: &sheets.AppendCellsRequest{
					Fields:  "*",                   // for adding data in all cells
					Rows:    prepareCells(records), // get formatted cells row
					SheetId: int64(SHEET_ID),       // use sheetID here
				},
			},
		},
	}

	// execute the request using spreadsheetId
	res, err := srv.Spreadsheets.BatchUpdate(SPREADSHEET_ID, &batchUpdateRequest).Context(ctx).Do()
	if err != nil || res.HTTPStatusCode != 200 {
		fmt.Print(err)
		return
	}
}

func prepareCells(records []string) []*sheets.RowData {
	// init cells array
	cells := []*sheets.CellData{}

	bgWhite := &sheets.Color{ // green background
		Alpha: 1,
		Blue:  0,
		Red:   0,
		Green: 1,
	}

	// prepare cell for each records and append it to cells
	for i := range records {
		data := &sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: &(records[i]), // add value
			},
			UserEnteredFormat: &sheets.CellFormat{ // add background color
				BackgroundColor: bgWhite,
			},
		}
		cells = append(cells, data)
	}

	// prepare row from cells
	return []*sheets.RowData{
		{Values: cells},
	}
}
