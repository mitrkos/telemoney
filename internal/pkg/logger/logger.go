package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func SetLogger() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{}),
	))
}
