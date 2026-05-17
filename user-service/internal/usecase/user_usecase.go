package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"sneaker-store/user-service/internal/model"
)

type UserUseCase struct {
	repo      UserRepository
	cache     UserCache
	publisher EventPublisher
}

func NewUserUseCase(repo UserRepository, cache UserCache, pub EventPublisher) *UserUseCase {
	return &UserUseCase{repo: repo, cache: cache, publisher: pub}
}

func (uc *UserUseCase) Create(ctx context.Context, email, password, fullName, phone string) (*model.User, error) {
	return uc.CreateWithRole(ctx, email, password, fullName, phone, model.RoleUser)
}

func (uc *UserUseCase) CreateWithRole(ctx context.Context, email, password, fullName, phone, role string) (*model.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if role == "" {
		role = model.RoleUser
	}

	if _, err := uc.repo.GetByEmail(ctx, email); err == nil {
		return nil, fmt.Errorf("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &model.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		Phone:        phone,
		Role:         role,
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	go func() {
		payload := map[string]string{
			"id":        u.ID,
			"email":     u.Email,
			"full_name": u.FullName,
		}
		if err := uc.publisher.Publish(context.Background(), "users.registered", payload); err != nil {
			log.Printf("publish users.registered failed: %v", err)
		}
	}()

	return u, nil
}

// EnsureAdmin creates the admin account if it does not already exist.
func (uc *UserUseCase) EnsureAdmin(ctx context.Context, email, password, fullName string) error {
	if _, err := uc.repo.GetByEmail(ctx, email); err == nil {
		return nil // already exists
	}
	_, err := uc.CreateWithRole(ctx, email, password, fullName, "", model.RoleAdmin)
	return err
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	if u, err := uc.cache.GetByID(ctx, id); err == nil && u != nil {
		return u, nil
	}

	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := uc.cache.SetByID(ctx, u); err != nil {
		log.Printf("cache set user failed: %v", err)
	}

	return u, nil
}

func (uc *UserUseCase) Update(ctx context.Context, id, fullName, phone string) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if fullName != "" {
		u.FullName = fullName
	}
	if phone != "" {
		u.Phone = phone
	}
	u.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete user failed: %v", err)
	}

	return u, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}

	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete user failed: %v", err)
	}

	return nil
}

// Authenticate returns (user, error). Caller reads user.Role for JWT issuance.
func (uc *UserUseCase) Authenticate(ctx context.Context, email, password string) (*model.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	u, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !u.Active {
		return nil, fmt.Errorf("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return u, nil
}

func (uc *UserUseCase) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *UserUseCase) ChangePassword(ctx context.Context, id, oldPassword, newPassword string) error {
	if id == "" || oldPassword == "" || newPassword == "" {
		return fmt.Errorf("id, old_password and new_password are required")
	}

	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid old password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := uc.repo.UpdatePassword(ctx, id, string(hash)); err != nil {
		return err
	}

	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete user failed: %v", err)
	}
	return nil
}

func (uc *UserUseCase) List(ctx context.Context, page, limit int) ([]*model.User, int, error) {
	return uc.repo.List(ctx, page, limit)
}

func (uc *UserUseCase) Search(ctx context.Context, query string) ([]*model.User, error) {
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}
	return uc.repo.Search(ctx, query)
}

func (uc *UserUseCase) GetStats(ctx context.Context) (*model.UserStats, error) {
	return uc.repo.GetStats(ctx)
}

func (uc *UserUseCase) UpdateStatus(ctx context.Context, id string, active bool) (*model.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	if err := uc.repo.UpdateStatus(ctx, id, active); err != nil {
		return nil, err
	}
	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete user failed: %v", err)
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *UserUseCase) ResetPassword(ctx context.Context, id, newPassword string) error {
	if id == "" || newPassword == "" {
		return fmt.Errorf("id and new_password are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err := uc.repo.UpdatePassword(ctx, id, string(hash)); err != nil {
		return err
	}
	if err := uc.cache.DeleteByID(ctx, id); err != nil {
		log.Printf("cache delete user failed: %v", err)
	}
	return nil
}
