package main

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/nats-io/nats.go"

	"sneaker-store/notification-service/internal/email"
	"sneaker-store/notification-service/internal/subscriber"
)

func main() {
	natsURL := envOr("NATS_URL", nats.DefaultURL)

	const maxRetries = 6
	var nc *nats.Conn
	var err error

	for i := 0; i < maxRetries; i++ {
		nc, err = nats.Connect(natsURL)
		if err == nil {
			break
		}
		wait := time.Duration(math.Pow(2, float64(i))) * time.Second
		log.Printf("nats connect attempt %d/%d failed: %v — retrying in %s", i+1, maxRetries, err, wait)
		time.Sleep(wait)
	}
	if err != nil {
		log.Fatalf("nats connect failed after %d retries: %v", maxRetries, err)
	}
	defer nc.Close()
	log.Println("connected to NATS")

	sender := email.NewSender()
	sub := subscriber.New(nc, sender)

	if err := sub.Start(); err != nil {
		log.Fatalf("subscriber error: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
