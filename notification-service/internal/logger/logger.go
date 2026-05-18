package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var std = log.New(os.Stdout, "", 0)

type LogEntry struct {
	Time    string      `json:"time"`
	Subject string      `json:"subject"`
	Event   interface{} `json:"event"`
}

func LogEvent(subject string, event interface{}) {
	entry := LogEntry{
		Time:    time.Now().UTC().Format(time.RFC3339),
		Subject: subject,
		Event:   event,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		std.Printf(`{"time":%q,"error":"marshal failed: %v"}`, time.Now().UTC().Format(time.RFC3339), err)
		return
	}
	std.Println(string(b))
}
