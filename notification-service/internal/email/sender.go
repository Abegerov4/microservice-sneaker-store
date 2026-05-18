package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type Sender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSender() *Sender {
	return &Sender{
		host:     envOr("SMTP_HOST", "smtp.gmail.com"),
		port:     envOr("SMTP_PORT", "587"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     envOr("SMTP_FROM", os.Getenv("SMTP_USERNAME")),
	}
}

func (s *Sender) Send(to, subject, body string) error {
	if s.username == "" || s.password == "" {
		log.Printf("[email] credentials not configured — skipping send to %s: %s", to, subject)
		return nil
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := []byte(fmt.Sprintf(
		"From: Sneaker Store <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.from, to, subject, body,
	))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	if err := smtp.SendMail(addr, auth, s.from, []string{to}, msg); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	log.Printf("[email] sent to %s: %s", to, subject)
	return nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
