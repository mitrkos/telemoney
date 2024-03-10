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
	TgAuthTokenTest      string
	GSheetsAuthToken string
}

func readConfig() (*Config, error) {
	viper.BindEnv("env", "TELEMONEY_ENV")
	viper.BindEnv("tg.auth_token", "TELEMONEY_TG_BOT_TOKEN")
	viper.BindEnv("tg.auth_token_test", "TELEMONEY_TG_BOT_TOKEN_TEST")
	viper.BindEnv("gsheets.auth_token", "TELEMONEY_GAUTH_TOKEN")

	viper.SetConfigName("telemoney")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := Config{
		Env:                    viper.GetString("env"),
		SpreadsheetID:          viper.GetString("gsheets.spreadsheet_id"),
		TransactionSheetID:     viper.GetString("gsheets.transaction_sheet_id"),
		TransactionSheetIDTest: viper.GetString("gsheets.transaction_sheet_id_test"),

		TgAuthToken:      viper.GetString("tg.auth_token"),
		TgAuthTokenTest:      viper.GetString("tg.auth_token_test"),
		GSheetsAuthToken: viper.GetString("gsheets.auth_token"),
	}

	if config.Env == "" || config.SpreadsheetID == "" || config.TransactionSheetID == "" || config.TransactionSheetIDTest == "" || config.TgAuthToken == "" || config.GSheetsAuthToken == "" {
		slog.Error("Config parsing failed", slog.Any("parsedConfig", config))
		return nil, errors.New("Config is not complete")
	}

	return &config, nil
}
