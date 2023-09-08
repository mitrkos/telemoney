package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	// //Append value to the sheet.
	// row := &sheets.ValueRange{
	// 	Values: [][]interface{}{{"1", "ABC", "abc@gmail.com"}},
	// }

	// response2, err := srv.Spreadsheets.Values.Append(SPREADSHEET_ID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
	// if err != nil || response2.HTTPStatusCode != 200 {
	// 	fmt.Print(err)
	// 	return
	// }

	token, err := os.ReadFile("./local/tg/token.txt")
	if err != nil {
        panic(err)
    }

    bot, err := tgbotapi.NewBotAPI(string(token))
    if err != nil {
        panic(err)
    }

    bot.Debug = true

	fmt.Printf("%+v\n", bot.Self)

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
    // to make sure Telegram knows we've handled previous values and we don't
    // need them repeated.
    updateConfig := tgbotapi.NewUpdate(0)

    // Tell Telegram we should wait up to 30 seconds on each request for an
    // update. This way we can get information just as quickly as making many
    // frequent requests without having to send nearly as many.
    updateConfig.Timeout = 30

    // Start polling Telegram for updates.
    updates := bot.GetUpdatesChan(updateConfig)

    // Let's go through each update that we're getting from Telegram.
    for update := range updates {
        // Telegram can send many types of updates depending on what your Bot
        // is up to. We only want to look at messages for now, so we can
        // discard any other updates.
        if update.Message == nil {
            continue
        }

        // Now that we know we've gotten a new message, we can construct a
        // reply! We'll take the Chat ID and Text from the incoming message
        // and use it to create a new message.
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
        // We'll also say that this message is a reply to the previous message.
        // For any other specifications than Chat ID or Text, you'll need to
        // set fields on the `MessageConfig`.
        msg.ReplyToMessageID = update.Message.MessageID

		//Append value to the sheet.
		row := &sheets.ValueRange{
			Values: [][]interface{}{{update.Message.Text}},
		}

		response2, err := srv.Spreadsheets.Values.Append(SPREADSHEET_ID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
		if err != nil || response2.HTTPStatusCode != 200 {
			fmt.Print(err)
			return
		}

        // Okay, we're sending our message off! We don't care about the message
        // we just sent, so we'll discard it.
        // if _, err := bot.Send(msg); err != nil {
        //     // Note that panics are a bad way to handle errors. Telegram can
        //     // have service outages`` or network errors, you should retry sending
        //     // messages or more gracefully handle failures.
        //     panic(err)
        // }
    }
}
