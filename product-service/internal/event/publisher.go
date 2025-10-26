package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Event struct {
	EventType  string      `json:"event_type"`
	OccurredAt string      `json:"occurred_at"`
	Payload    interface{} `json:"payload"`
}

type NATSPublisher struct {
	conn *nats.Conn
}

func NewNATSPublisher(url string) (*NATSPublisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}
	log.Println("connected to NATS")
	return &NATSPublisher{conn: conn}, nil
}

func NewNATSPublisherConn(conn *nats.Conn) *NATSPublisher {
	return &NATSPublisher{conn: conn}
}

func (p *NATSPublisher) Publish(_ context.Context, subject string, data interface{}) error {
	ev := Event{
		EventType:  subject,
		OccurredAt: time.Now().UTC().Format(time.RFC3339),
		Payload:    data,
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return p.conn.Publish(subject, b)
}
