package logging

import (
	"encoding/json"
	"os"
	"time"
)

type Logger struct{}

func New() *Logger { return &Logger{} }

func (l *Logger) Info(msg string, fields ...any)  { l.log("INFO", msg, fields) }
func (l *Logger) Warn(msg string, fields ...any)  { l.log("WARN", msg, fields) }
func (l *Logger) Error(msg string, fields ...any) { l.log("ERROR", msg, fields) }

func (l *Logger) log(level, msg string, fields []any) {
	entry := map[string]any{"ts": time.Now().UTC().Format(time.RFC3339), "level": level, "msg": msg}
	for i := 0; i+1 < len(fields); i += 2 {
		if k, ok := fields[i].(string); ok {
			entry[k] = fields[i+1]
		}
	}
	json.NewEncoder(os.Stderr).Encode(entry)
}
