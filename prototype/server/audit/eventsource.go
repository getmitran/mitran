package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

type Event struct {
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Actor     string          `json:"actor"`
	Action    string          `json:"action"`
	Resource  string          `json:"resource"`
	Data      json.RawMessage `json:"data"`
}

type EventLog struct {
	file *os.File
}

func OpenLog(path string) (*EventLog, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &EventLog{file: f}, nil
}

func (el *EventLog) Append(e Event) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = el.file.Write(data)
	return err
}

func (el *EventLog) ReadAll() ([]Event, error) {
	if _, err := el.file.Seek(0, 0); err != nil {
		return nil, err
	}
	var events []Event
	scanner := bufio.NewScanner(el.file)
	for scanner.Scan() {
		var e Event
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, scanner.Err()
}
