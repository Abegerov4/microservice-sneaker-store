package model

import "time"

type ChatMessage struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"` // "user" | "assistant"
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type SneakerRecommendation struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	Brand      string  `json:"brand"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
	Reason     string  `json:"reason"`
	MatchScore float64 `json:"match_score"`
}

type TrendingSneaker struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	TrendReason string  `json:"trend_reason"`
	TrendScore  int32   `json:"trend_score"`
}

type ProductInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Brand       string   `json:"brand"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Sizes       []string `json:"sizes"`
	Stock       int32    `json:"stock"`
	ImageURL    string   `json:"image_url"`
}
