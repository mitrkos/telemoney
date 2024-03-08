package telemoney

import (
	"errors"
	"log/slog"

	"github.com/spf13/viper"
)

type Config struct {
	Env                    string // TODO: use enum
	SpreadsheetID          string
	TransactionSheetID     string
	TransactionSheetIDTest string

	TgAuthToken      string
	GSheetsAuthToken string
}

func readConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.SetConfigName("telemoney")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	err = viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	config := Config{
		Env:                    viper.GetString("env"),
		SpreadsheetID:          viper.GetString("gsheets.spreadsheet_id"),
		TransactionSheetID:     viper.GetString("gsheets.transaction_sheet_id"),
		TransactionSheetIDTest: viper.GetString("gsheets.transaction_sheet_id_test"),

		TgAuthToken:      viper.GetString("TELEMONEY_TG_BOT_TOKEN"),
		GSheetsAuthToken: viper.GetString("TELEMONEY_GAUTH_TOKEN"),
	}

	if config.Env == "" || config.SpreadsheetID == "" || config.TransactionSheetID == "" || config.TransactionSheetIDTest == "" || config.TgAuthToken == "" || config.GSheetsAuthToken == "" {
		slog.Error("Config parsing failed", slog.Any("parsedConfig", config))
		return nil, errors.New("Config is not complete")
	}

	return &config, nil
}
