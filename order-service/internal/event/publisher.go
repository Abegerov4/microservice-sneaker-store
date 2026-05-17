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

func NewNATSPublisherConn(conn *nats.Conn) *NATSPublisher {
	return &NATSPublisher{conn: conn}
}

func (p *NATSPublisher) Publish(_ context.Context, subject string, data interface{}) error {
	if p.conn == nil {
		return fmt.Errorf("nats not connected")
	}
	ev := Event{
		EventType:  subject,
		OccurredAt: time.Now().UTC().Format(time.RFC3339),
		Payload:    data,
	}
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if err := p.conn.Publish(subject, b); err != nil {
		log.Printf("nats publish %s failed: %v", subject, err)
		return err
	}
	return nil
}
