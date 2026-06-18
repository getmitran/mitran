package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Logger *slog.Logger

func Init(level string) {
	var l slog.Level
	switch strings.ToLower(level) {
	case "debug":
		l = slog.LevelDebug
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}

func Info(msg string, args ...any)  { Logger.Info(msg, args...) }
func Error(msg string, args ...any) { Logger.Error(msg, args...) }
func Debug(msg string, args ...any) { Logger.Debug(msg, args...) }
func Warn(msg string, args ...any)  { Logger.Warn(msg, args...) }

func init() {
	Init(os.Getenv("MITRAN_LOG_LEVEL"))
}
