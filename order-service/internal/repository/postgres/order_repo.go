package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"sneaker-store/order-service/internal/model"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *model.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO orders (id, user_id, total_amount, status, shipping_address, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		o.ID, o.UserID, o.TotalAmount, string(o.Status), o.ShippingAddress, o.CreatedAt, o.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	for _, item := range o.Items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (id, order_id, product_id, product_name, quantity, price, size)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			uuid.NewString(), o.ID, item.ProductID, item.ProductName, item.Quantity, item.Price, item.Size,
		)
		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*model.Order, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, total_amount, status, shipping_address, created_at, updated_at
		 FROM orders WHERE id = $1`, id)

	o := &model.Order{}
	var statusStr string
	if err := row.Scan(&o.ID, &o.UserID, &o.TotalAmount, &statusStr, &o.ShippingAddress, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	o.Status = model.OrderStatus(statusStr)

	items, err := r.getItems(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return o, nil
}

func (r *OrderRepository) getItems(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	rows, err := r.db.Query(ctx,
		`SELECT product_id, product_name, quantity, price, size FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.Quantity, &item.Price, &item.Size); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]*model.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, total_amount, status, shipping_address, created_at, updated_at
		 FROM orders ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	res, err := r.db.Exec(ctx,
		`UPDATE orders SET status=$2, updated_at=now() WHERE id=$1`, id, string(status))
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (r *OrderRepository) GetByStatus(ctx context.Context, orderStatus model.OrderStatus) ([]*model.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, total_amount, status, shipping_address, created_at, updated_at
		 FROM orders WHERE status = $1 ORDER BY created_at DESC`, string(orderStatus))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) GetStats(ctx context.Context) (*model.OrderStats, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status='pending'),
			COUNT(*) FILTER (WHERE status='confirmed'),
			COUNT(*) FILTER (WHERE status='shipped'),
			COUNT(*) FILTER (WHERE status='delivered'),
			COUNT(*) FILTER (WHERE status='cancelled'),
			COALESCE(SUM(total_amount),0)
		FROM orders`)
	s := &model.OrderStats{}
	if err := row.Scan(&s.TotalOrders, &s.PendingOrders, &s.ConfirmedOrders,
		&s.ShippedOrders, &s.DeliveredOrders, &s.CancelledOrders, &s.TotalRevenue); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *OrderRepository) GetTotalRevenue(ctx context.Context) (float64, error) {
	var total float64
	err := r.db.QueryRow(ctx, `SELECT COALESCE(SUM(total_amount),0) FROM orders WHERE status != 'cancelled'`).Scan(&total)
	return total, err
}

func (r *OrderRepository) GetByDateRange(ctx context.Context, from, to time.Time) ([]*model.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, total_amount, status, shipping_address, created_at, updated_at
		 FROM orders WHERE created_at >= $1 AND created_at <= $2 ORDER BY created_at DESC`,
		from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanOrders(ctx, rows)
}

func (r *OrderRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *OrderRepository) scanOrders(ctx context.Context, rows interface {
	Next() bool
	Scan(...interface{}) error
	Err() error
}) ([]*model.Order, error) {
	var orders []*model.Order
	for rows.Next() {
		o := &model.Order{}
		var statusStr string
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &statusStr, &o.ShippingAddress, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		o.Status = model.OrderStatus(statusStr)
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, o := range orders {
		items, err := r.getItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}
		o.Items = items
	}
	return orders, nil
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, total_amount, status, shipping_address, created_at, updated_at
		 FROM orders WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanOrders(ctx, rows)
}
