package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"sneaker-store/user-service/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

const selectCols = `id, email, password_hash, full_name, phone, active, role, created_at, updated_at`

func scanUser(row interface{ Scan(...interface{}) error }) (*model.User, error) {
	u := &model.User{}
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone, &u.Active, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	role := u.Role
	if role == "" {
		role = model.RoleUser
	}
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, full_name, phone, active, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		u.ID, u.Email, u.PasswordHash, u.FullName, u.Phone, true, role, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+selectCols+` FROM users WHERE id = $1`, id)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+selectCols+` FROM users WHERE email = $1`, email)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET full_name=$2, phone=$3, updated_at=$4 WHERE id=$1`,
		u.ID, u.FullName, u.Phone, u.UpdatedAt,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, page, limit int) ([]*model.User, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT `+selectCols+` FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *UserRepository) Search(ctx context.Context, query string) ([]*model.User, error) {
	pattern := "%" + query + "%"
	rows, err := r.db.Query(ctx,
		`SELECT `+selectCols+` FROM users WHERE full_name ILIKE $1 OR email ILIKE $1 ORDER BY full_name`,
		pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) GetStats(ctx context.Context) (*model.UserStats, error) {
	s := &model.UserStats{}
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*), COUNT(*) FILTER (WHERE active = true) FROM users`).
		Scan(&s.TotalUsers, &s.ActiveUsers)
	return s, err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id, newHash string) error {
	res, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash=$2, updated_at=now() WHERE id=$1`, id, newHash)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id string, active bool) error {
	res, err := r.db.Exec(ctx,
		`UPDATE users SET active=$2, updated_at=now() WHERE id=$1`, id, active)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}
