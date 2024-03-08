package telemoney

import (
	"errors"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Env                    string // TODO: use enum
	SpreadsheetID          string
	TransactionSheetID     string
	TransactionSheetIDTest string

	TgAuthToken      string
	GSheetsAuthToken string
}

type ConfigFile struct {
	Env                    string
	SpreadsheetID          string
	TransactionSheetID     string
	TransactionSheetIDTest string
}

func readConfig() (*Config, error) {
	// sensitive params are set via env vars
	tgAuthToken, ok := os.LookupEnv("TG_BOT_TOKEN")
	if !ok {
		return nil, errors.New("TG_BOT_TOKEN is not set")
	}
	gSheetsAuthToken, ok := os.LookupEnv("GAUTH_TOKEN")
	if !ok {
		return nil, errors.New("GAUTH_TOKEN is not set")
	}

	// app params are set via ./config/telemoney.toml
	var configFile ConfigFile
	_, err := toml.DecodeFile("./config/telemoney.toml", &configFile)
	if err != nil {
		return nil, err
	}

	slog.Info("Config loaded", slog.Any("configFile", configFile))

	return &Config{
		Env:                    configFile.Env,
		SpreadsheetID:          configFile.SpreadsheetID,
		TransactionSheetID:     configFile.TransactionSheetID,
		TransactionSheetIDTest: configFile.TransactionSheetIDTest,

		TgAuthToken:      tgAuthToken,
		GSheetsAuthToken: gSheetsAuthToken,
	}, nil
}
