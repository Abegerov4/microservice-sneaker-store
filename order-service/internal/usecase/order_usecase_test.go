package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"sneaker-store/order-service/internal/model"
	"sneaker-store/order-service/internal/usecase"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockOrderRepo struct {
	orders map[string]*model.Order
	err    error
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{orders: make(map[string]*model.Order)}
}

func (m *mockOrderRepo) Create(_ context.Context, o *model.Order) error {
	if m.err != nil {
		return m.err
	}
	m.orders[o.ID] = o
	return nil
}

func (m *mockOrderRepo) GetByID(_ context.Context, id string) (*model.Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	o, ok := m.orders[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cp := *o
	return &cp, nil
}

func (m *mockOrderRepo) List(_ context.Context) ([]*model.Order, error) {
	out := make([]*model.Order, 0, len(m.orders))
	for _, o := range m.orders {
		out = append(out, o)
	}
	return out, nil
}

func (m *mockOrderRepo) UpdateStatus(_ context.Context, id string, status model.OrderStatus) error {
	o, ok := m.orders[id]
	if !ok {
		return errors.New("not found")
	}
	o.Status = status
	return nil
}

func (m *mockOrderRepo) GetByUserID(_ context.Context, userID string) ([]*model.Order, error) {
	var out []*model.Order
	for _, o := range m.orders {
		if o.UserID == userID {
			out = append(out, o)
		}
	}
	return out, nil
}

func (m *mockOrderRepo) GetByStatus(_ context.Context, status model.OrderStatus) ([]*model.Order, error) {
	var out []*model.Order
	for _, o := range m.orders {
		if o.Status == status {
			out = append(out, o)
		}
	}
	return out, nil
}

func (m *mockOrderRepo) GetStats(_ context.Context) (*model.OrderStats, error) {
	return &model.OrderStats{TotalOrders: len(m.orders)}, nil
}

func (m *mockOrderRepo) GetTotalRevenue(_ context.Context) (float64, error) {
	var total float64
	for _, o := range m.orders {
		total += o.TotalAmount
	}
	return total, nil
}

func (m *mockOrderRepo) GetByDateRange(_ context.Context, from, to time.Time) ([]*model.Order, error) {
	var out []*model.Order
	for _, o := range m.orders {
		if (o.CreatedAt.Equal(from) || o.CreatedAt.After(from)) && o.CreatedAt.Before(to) {
			out = append(out, o)
		}
	}
	return out, nil
}

func (m *mockOrderRepo) CountByUserID(_ context.Context, userID string) (int, error) {
	count := 0
	for _, o := range m.orders {
		if o.UserID == userID {
			count++
		}
	}
	return count, nil
}

type mockOrderCache struct{}

func (m *mockOrderCache) GetByID(_ context.Context, _ string) (*model.Order, error) {
	return nil, errors.New("miss")
}
func (m *mockOrderCache) SetByID(_ context.Context, _ *model.Order) error { return nil }
func (m *mockOrderCache) DeleteByID(_ context.Context, _ string) error    { return nil }
func (m *mockOrderCache) GetByUser(_ context.Context, _ string) ([]*model.Order, error) {
	return nil, errors.New("miss")
}
func (m *mockOrderCache) SetByUser(_ context.Context, _ string, _ []*model.Order) error {
	return nil
}
func (m *mockOrderCache) DeleteByUser(_ context.Context, _ string) error { return nil }

type mockOrderPublisher struct{ subjects []string }

func (m *mockOrderPublisher) Publish(_ context.Context, subject string, _ interface{}) error {
	m.subjects = append(m.subjects, subject)
	return nil
}

type mockProductClient struct {
	product *usecase.ProductInfo
	err     error
}

func (m *mockProductClient) GetProduct(_ context.Context, _ string) (*usecase.ProductInfo, error) {
	return m.product, m.err
}

func newOrderUC() (*usecase.OrderUseCase, *mockOrderRepo, *mockOrderPublisher, *mockProductClient) {
	repo := newMockOrderRepo()
	pub := &mockOrderPublisher{}
	pc := &mockProductClient{product: &usecase.ProductInfo{ID: "prod-1", Name: "Air Max 90", Price: 130}}
	uc := usecase.NewOrderUseCase(repo, &mockOrderCache{}, pub, pc)
	return uc, repo, pub, pc
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestCreateOrder_Success(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	order, err := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 2, Size: "42"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID == "" {
		t.Error("expected order ID to be generated")
	}
	if order.TotalAmount != 260 {
		t.Errorf("got total %v, want 260", order.TotalAmount)
	}
	if order.Status != model.StatusPending {
		t.Errorf("got status %q, want pending", order.Status)
	}
}

func TestCreateOrder_MissingUserID(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	_, err := uc.Create(context.Background(), usecase.CreateOrderInput{
		Items: []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})
	if err == nil {
		t.Error("expected error for missing user_id")
	}
}

func TestCreateOrder_EmptyItems(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	_, err := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{},
	})
	if err == nil {
		t.Error("expected error for empty items")
	}
}

func TestCreateOrder_ZeroQuantity(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	_, err := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 0}},
	})
	if err == nil {
		t.Error("expected error for zero quantity")
	}
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	uc, _, _, pc := newOrderUC()
	pc.err = errors.New("not found")

	_, err := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "missing", Quantity: 1}},
	})
	if err == nil {
		t.Error("expected error when product not found")
	}
}

func TestUpdateStatus_Success(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	order, _ := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})

	updated, err := uc.UpdateStatus(context.Background(), order.ID, model.StatusConfirmed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != model.StatusConfirmed {
		t.Errorf("got status %q, want confirmed", updated.Status)
	}
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	order, _ := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})

	_, err := uc.UpdateStatus(context.Background(), order.ID, "invalid_status")
	if err == nil {
		t.Error("expected error for invalid status")
	}
}

func TestUpdateStatus_CannotUpdateCancelled(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	order, _ := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})

	uc.Cancel(context.Background(), order.ID)

	_, err := uc.UpdateStatus(context.Background(), order.ID, model.StatusConfirmed)
	if err == nil {
		t.Error("expected error when updating a cancelled order")
	}
}

func TestCancelOrder(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	order, _ := uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})

	cancelled, err := uc.Cancel(context.Background(), order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cancelled.Status != model.StatusCancelled {
		t.Errorf("got status %q, want cancelled", cancelled.Status)
	}
}

func TestGetByUserID(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})
	uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 2}},
	})
	uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-2",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})

	orders, err := uc.GetByUserID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("got %d orders for user-1, want 2", len(orders))
	}
}

func TestGetByUserID_EmptyID(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	_, err := uc.GetByUserID(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty user_id")
	}
}

func TestCountByUserID(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})
	uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})

	count, err := uc.CountByUserID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("got count %d, want 2", count)
	}
}

func TestGetTotalRevenue(t *testing.T) {
	uc, _, _, _ := newOrderUC()

	uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-1",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 2}},
	})
	uc.Create(context.Background(), usecase.CreateOrderInput{
		UserID: "user-2",
		Items:  []usecase.OrderItemInput{{ProductID: "prod-1", Quantity: 1}},
	})

	revenue, err := uc.GetTotalRevenue(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if revenue != 390 {
		t.Errorf("got revenue %v, want 390", revenue)
	}
}
