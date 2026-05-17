package model

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusShipped, StatusDelivered, StatusCancelled:
		return true
	}
	return false
}

type OrderItem struct {
	ProductID   string
	ProductName string
	Quantity    int
	Price       float64
	Size        string
}

type Order struct {
	ID              string
	UserID          string
	Items           []OrderItem
	TotalAmount     float64
	Status          OrderStatus
	ShippingAddress string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderStats struct {
	TotalOrders     int
	PendingOrders   int
	ConfirmedOrders int
	ShippedOrders   int
	DeliveredOrders int
	CancelledOrders int
	TotalRevenue    float64
}
