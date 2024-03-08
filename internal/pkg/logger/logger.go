package logger

import(
	"os"
	"log/slog"
	"github.com/lmittmann/tint"
)

func SetLogger() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{}),
	))
}