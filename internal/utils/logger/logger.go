package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

func Init(mode string) {
	var handler slog.Handler

	switch mode {
	case "prod", "production":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	log = slog.New(handler)
	slog.SetDefault(log)
}

func Info(msg string, args ...any)  { log.Info(msg, args...) }
func Debug(msg string, args ...any) { log.Debug(msg, args...) }
func Warn(msg string, args ...any)  { log.Warn(msg, args...) }
func Error(msg string, args ...any) { log.Error(msg, args...) }
