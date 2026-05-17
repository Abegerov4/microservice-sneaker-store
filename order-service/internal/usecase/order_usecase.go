package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"sneaker-store/order-service/internal/model"
)


type OrderUseCase struct {
	repo          OrderRepository
	cache         OrderCache
	publisher     EventPublisher
	productClient ProductClient
}

func NewOrderUseCase(repo OrderRepository, cache OrderCache, pub EventPublisher, pc ProductClient) *OrderUseCase {
	return &OrderUseCase{repo: repo, cache: cache, publisher: pub, productClient: pc}
}

type CreateOrderInput struct {
	UserID          string
	Items           []OrderItemInput
	ShippingAddress string
}

type OrderItemInput struct {
	ProductID string
	Quantity  int
	Size      string
}

func (uc *OrderUseCase) Create(ctx context.Context, input CreateOrderInput) (*model.Order, error) {
	if input.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("at least one item is required")
	}

	var items []model.OrderItem
	var total float64

	for _, item := range input.Items {
		if item.ProductID == "" {
			return nil, fmt.Errorf("product_id is required for each item")
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be positive")
		}

		product, err := uc.productClient.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %w", item.ProductID, err)
		}

		itemTotal := product.Price * float64(item.Quantity)
		total += itemTotal

		items = append(items, model.OrderItem{
			ProductID:   item.ProductID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Price:       product.Price,
			Size:        item.Size,
		})
	}

	order := &model.Order{
		ID:              uuid.NewString(),
		UserID:          input.UserID,
		Items:           items,
		TotalAmount:     total,
		Status:          model.StatusPending,
		ShippingAddress: input.ShippingAddress,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := uc.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	if err := uc.cache.DeleteByUser(ctx, input.UserID); err != nil {
		log.Printf("cache delete by user failed: %v", err)
	}

	go func() {
		if err := uc.publisher.Publish(context.Background(), "orders.created", order); err != nil {
			log.Printf("publish orders.created failed: %v", err)
		}
	}()

	return order, nil
}

func (uc *OrderUseCase) GetByID(ctx context.Context, id string) (*model.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	if o, err := uc.cache.GetByID(ctx, id); err == nil && o != nil {
		return o, nil
	}

	o, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.SetByID(ctx, o); err != nil {
		log.Printf("cache set order failed: %v", err)
	}

	return o, nil
}

func (uc *OrderUseCase) List(ctx context.Context) ([]*model.Order, error) {
	return uc.repo.List(ctx)
}

func (uc *OrderUseCase) UpdateStatus(ctx context.Context, id string, newStatus model.OrderStatus) (*model.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	if !newStatus.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", newStatus)
	}

	order, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order.Status == model.StatusCancelled {
		return nil, fmt.Errorf("cannot update a cancelled order")
	}
	if order.Status == model.StatusDelivered && newStatus != model.StatusDelivered {
		return nil, fmt.Errorf("cannot change status from delivered")
	}

	oldStatus := order.Status
	order.Status = newStatus
	order.UpdatedAt = time.Now()

	if err := uc.repo.UpdateStatus(ctx, id, newStatus); err != nil {
		return nil, err
	}

	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete order failed: %v", err)
	}
	if err := uc.cache.DeleteByUser(ctx, order.UserID); err != nil {
		log.Printf("cache delete by user failed: %v", err)
	}

	go func() {
		payload := map[string]string{
			"id":         id,
			"old_status": string(oldStatus),
			"new_status": string(newStatus),
			"user_id":    order.UserID,
		}
		if err := uc.publisher.Publish(context.Background(), "orders.status_updated", payload); err != nil {
			log.Printf("publish orders.status_updated failed: %v", err)
		}
	}()

	return order, nil
}

func (uc *OrderUseCase) Cancel(ctx context.Context, id string) (*model.Order, error) {
	return uc.UpdateStatus(ctx, id, model.StatusCancelled)
}

func (uc *OrderUseCase) GetByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if orders, err := uc.cache.GetByUser(ctx, userID); err == nil && orders != nil {
		return orders, nil
	}

	orders, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.SetByUser(ctx, userID, orders); err != nil {
		log.Printf("cache set by user failed: %v", err)
	}

	return orders, nil
}

func (uc *OrderUseCase) GetByStatus(ctx context.Context, orderStatus model.OrderStatus) ([]*model.Order, error) {
	if !orderStatus.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", orderStatus)
	}
	return uc.repo.GetByStatus(ctx, orderStatus)
}

func (uc *OrderUseCase) GetStats(ctx context.Context) (*model.OrderStats, error) {
	return uc.repo.GetStats(ctx)
}

func (uc *OrderUseCase) GetItemsByOrderID(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}
	o, err := uc.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return o.Items, nil
}

func (uc *OrderUseCase) GetTotalRevenue(ctx context.Context) (float64, error) {
	return uc.repo.GetTotalRevenue(ctx)
}

func (uc *OrderUseCase) GetByDateRange(ctx context.Context, from, to time.Time) ([]*model.Order, error) {
	return uc.repo.GetByDateRange(ctx, from, to)
}

func (uc *OrderUseCase) CountByUserID(ctx context.Context, userID string) (int, error) {
	if userID == "" {
		return 0, fmt.Errorf("user_id is required")
	}
	return uc.repo.CountByUserID(ctx, userID)
}
