package usecase

import (
	"context"

	"sneaker-store/user-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*model.User, int, error)
	Search(ctx context.Context, query string) ([]*model.User, error)
	GetStats(ctx context.Context) (*model.UserStats, error)
	UpdatePassword(ctx context.Context, id, newHash string) error
	UpdateStatus(ctx context.Context, id string, active bool) error
}

type UserCache interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	SetByID(ctx context.Context, u *model.User) error
	DeleteByID(ctx context.Context, id string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}
