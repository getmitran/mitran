package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Level string

const (
	Info  Level = "INFO"
	Warn  Level = "WARN"
	Error Level = "ERROR"
	Debug Level = "DEBUG"
)

type Entry struct {
	Timestamp string                 `json:"ts"`
	Level     Level                  `json:"level"`
	Msg       string                 `json:"msg"`
	Component string                 `json:"component,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

type Logger struct {
	component string
	level     Level
}

func New(component string) *Logger { return &Logger{component: component, level: Info} }

func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
	e := Entry{Timestamp: time.Now().Format(time.RFC3339), Level: level, Msg: msg, Component: l.component, Fields: fields}
	data, _ := json.Marshal(e)
	fmt.Fprintln(os.Stdout, string(data))
}

func (l *Logger) Info(msg string, fields ...map[string]interface{})  { l.log(Info, msg, merge(fields)) }
func (l *Logger) Warn(msg string, fields ...map[string]interface{})  { l.log(Warn, msg, merge(fields)) }
func (l *Logger) Error(msg string, fields ...map[string]interface{}) { l.log(Error, msg, merge(fields)) }
func (l *Logger) Debug(msg string, fields ...map[string]interface{}) { l.log(Debug, msg, merge(fields)) }

func merge(fields []map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}
	return fields[0]
}
