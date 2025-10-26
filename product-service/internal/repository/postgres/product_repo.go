package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sneaker-store/product-service/internal/model"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO products (id, name, brand, description, price, sizes, stock, image_url, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		p.ID, p.Name, p.Brand, p.Description, p.Price, p.Sizes, p.Stock, p.ImageURL, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*model.Product, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, brand, description, price, sizes, stock, image_url, created_at, updated_at
		 FROM products WHERE id = $1`, id)

	p := &model.Product{}
	err := row.Scan(&p.ID, &p.Name, &p.Brand, &p.Description, &p.Price, &p.Sizes, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return p, nil
}

func (r *ProductRepository) List(ctx context.Context) ([]*model.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, brand, description, price, sizes, stock, image_url, created_at, updated_at
		 FROM products ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Brand, &p.Description, &p.Price, &p.Sizes, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *model.Product) error {
	_, err := r.db.Exec(ctx,
		`UPDATE products SET name=$2, brand=$3, description=$4, price=$5, sizes=$6, stock=$7, image_url=$8, updated_at=$9
		 WHERE id=$1`,
		p.ID, p.Name, p.Brand, p.Description, p.Price, p.Sizes, p.Stock, p.ImageURL, p.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *ProductRepository) UpdateStock(ctx context.Context, id string, delta int) (int, error) {
	var newStock int
	err := r.db.QueryRow(ctx,
		`UPDATE products SET stock = stock + $2, updated_at = now() WHERE id = $1 RETURNING stock`,
		id, delta,
	).Scan(&newStock)
	if err != nil {
		return 0, fmt.Errorf("update stock: %w", err)
	}
	return newStock, nil
}

func (r *ProductRepository) GetByBrand(ctx context.Context, brand string) ([]*model.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, brand, description, price, sizes, stock, image_url, created_at, updated_at
		 FROM products WHERE brand ILIKE $1 ORDER BY name`,
		"%"+brand+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (r *ProductRepository) GetLowStock(ctx context.Context, threshold int) ([]*model.Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, brand, description, price, sizes, stock, image_url, created_at, updated_at
		 FROM products WHERE stock <= $1 ORDER BY stock ASC`,
		threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (r *ProductRepository) GetBrands(ctx context.Context) ([]string, error) {
	rows, err := r.db.Query(ctx, `SELECT DISTINCT brand FROM products ORDER BY brand`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var brands []string
	for rows.Next() {
		var b string
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		brands = append(brands, b)
	}
	return brands, nil
}

func (r *ProductRepository) GetStats(ctx context.Context) (*model.ProductStats, error) {
	row := r.db.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(DISTINCT brand), COALESCE(SUM(stock),0), COALESCE(AVG(price),0)
		 FROM products`)
	s := &model.ProductStats{}
	if err := row.Scan(&s.TotalProducts, &s.TotalBrands, &s.TotalStock, &s.AveragePrice); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *ProductRepository) BulkDelete(ctx context.Context, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	args := make([]interface{}, len(ids))
	placeholders := make([]string, len(ids))
	for i, id := range ids {
		args[i] = id
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query := fmt.Sprintf("DELETE FROM products WHERE id IN (%s)", joinStrings(placeholders, ","))
	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return int(res.RowsAffected()), nil
}

func joinStrings(ss []string, sep string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func scanProducts(rows interface{ Next() bool; Scan(...interface{}) error; Err() error }) ([]*model.Product, error) {
	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Brand, &p.Description, &p.Price, &p.Sizes, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Search(ctx context.Context, brand string, minPrice, maxPrice float64, size string) ([]*model.Product, error) {
	query := `SELECT id, name, brand, description, price, sizes, stock, image_url, created_at, updated_at
	          FROM products WHERE 1=1`
	args := []interface{}{}
	i := 1

	if brand != "" {
		query += fmt.Sprintf(" AND brand ILIKE $%d", i)
		args = append(args, "%"+brand+"%")
		i++
	}
	if minPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", i)
		args = append(args, minPrice)
		i++
	}
	if maxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", i)
		args = append(args, maxPrice)
		i++
	}
	if size != "" {
		query += fmt.Sprintf(" AND $%d = ANY(sizes)", i)
		args = append(args, size)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		p := &model.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Brand, &p.Description, &p.Price, &p.Sizes, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
