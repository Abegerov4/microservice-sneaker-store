package model

import "time"

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	Active       bool
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserStats struct {
	TotalUsers  int
	ActiveUsers int
}
