package subscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"sneaker-store/notification-service/internal/email"
	"sneaker-store/notification-service/internal/logger"
)

type Event struct {
	EventType  string          `json:"event_type"`
	OccurredAt string          `json:"occurred_at"`
	Payload    json.RawMessage `json:"payload"`
}

type OrderCreatedPayload struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type OrderStatusPayload struct {
	ID        string `json:"id"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	UserID    string `json:"user_id"`
}

type UserRegisteredPayload struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type Subscriber struct {
	conn         *nats.Conn
	emailSender  *email.Sender
	notifyEmail  string
}

func New(conn *nats.Conn, sender *email.Sender) *Subscriber {
	return &Subscriber{
		conn:        conn,
		emailSender: sender,
		notifyEmail: envOr("NOTIFY_EMAIL", "notifications@sneakerstore.com"),
	}
}

func (s *Subscriber) Start() error {
	subs := []struct {
		subject string
		handler func([]byte)
	}{
		{"products.created", s.handleProductCreated},
		{"orders.created", s.handleOrderCreated},
		{"orders.status_updated", s.handleOrderStatusUpdated},
		{"users.registered", s.handleUserRegistered},
	}

	for _, sub := range subs {
		sub := sub
		if _, err := s.conn.Subscribe(sub.subject, func(msg *nats.Msg) {
			sub.handler(msg.Data)
		}); err != nil {
			return fmt.Errorf("subscribe %s: %w", sub.subject, err)
		}
		log.Printf("subscribed to %s", sub.subject)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down notification service…")
	s.conn.Drain()
	return nil
}

func (s *Subscriber) handleProductCreated(data []byte) {
	var ev Event
	if err := json.Unmarshal(data, &ev); err != nil {
		log.Printf("parse products.created: %v", err)
		return
	}
	logger.LogEvent("products.created", ev)
}

func (s *Subscriber) handleOrderCreated(data []byte) {
	var ev Event
	if err := json.Unmarshal(data, &ev); err != nil {
		log.Printf("parse orders.created: %v", err)
		return
	}
	logger.LogEvent("orders.created", ev)

	var payload OrderCreatedPayload
	if err := json.Unmarshal(ev.Payload, &payload); err != nil {
		log.Printf("parse order payload: %v", err)
		return
	}

	subject := "Your Sneaker Store order has been placed!"
	body := fmt.Sprintf(
		"Hello!\n\nYour order #%s has been placed successfully.\nStatus: %s\n\nThank you for shopping with Sneaker Store!\n",
		payload.ID, payload.Status,
	)
	if err := s.emailSender.Send(s.notifyEmail, subject, body); err != nil {
		log.Printf("send order created email: %v", err)
	}
}

func (s *Subscriber) handleOrderStatusUpdated(data []byte) {
	var ev Event
	if err := json.Unmarshal(data, &ev); err != nil {
		log.Printf("parse orders.status_updated: %v", err)
		return
	}
	logger.LogEvent("orders.status_updated", ev)

	var payload OrderStatusPayload
	if err := json.Unmarshal(ev.Payload, &payload); err != nil {
		log.Printf("parse status payload: %v", err)
		return
	}

	subject := fmt.Sprintf("Order #%s status update: %s", payload.ID, payload.NewStatus)
	body := fmt.Sprintf(
		"Hello!\n\nYour order #%s status has been updated.\nFrom: %s → To: %s\n\nThank you for shopping with Sneaker Store!\n",
		payload.ID, payload.OldStatus, payload.NewStatus,
	)
	if err := s.emailSender.Send(s.notifyEmail, subject, body); err != nil {
		log.Printf("send status updated email: %v", err)
	}
}

func (s *Subscriber) handleUserRegistered(data []byte) {
	var ev Event
	if err := json.Unmarshal(data, &ev); err != nil {
		log.Printf("parse users.registered: %v", err)
		return
	}
	logger.LogEvent("users.registered", ev)

	var payload UserRegisteredPayload
	if err := json.Unmarshal(ev.Payload, &payload); err != nil {
		log.Printf("parse user payload: %v", err)
		return
	}

	subject := "Welcome to Sneaker Store!"
	body := fmt.Sprintf(
		"Hello %s!\n\nWelcome to Sneaker Store. Your account has been created successfully.\n\nHappy shopping!\n",
		payload.FullName,
	)
	if err := s.emailSender.Send(payload.Email, subject, body); err != nil {
		log.Printf("send welcome email: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func init() {
	_ = time.RFC3339
}
