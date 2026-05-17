package usecase_test

import (
	"context"
	"errors"
	"testing"

	"sneaker-store/user-service/internal/model"
	"sneaker-store/user-service/internal/usecase"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockUserRepo struct {
	users map[string]*model.User
	err   error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*model.User)}
}

func (m *mockUserRepo) Create(_ context.Context, u *model.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id string) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cp := *u
	return &cp, nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) Update(_ context.Context, u *model.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) Delete(_ context.Context, id string) error {
	if _, ok := m.users[id]; !ok {
		return errors.New("not found")
	}
	delete(m.users, id)
	return nil
}

func (m *mockUserRepo) List(_ context.Context, _, _ int) ([]*model.User, int, error) {
	out := make([]*model.User, 0, len(m.users))
	for _, u := range m.users {
		out = append(out, u)
	}
	return out, len(out), nil
}

func (m *mockUserRepo) Search(_ context.Context, query string) ([]*model.User, error) {
	var out []*model.User
	for _, u := range m.users {
		if u.Email == query || u.FullName == query {
			out = append(out, u)
		}
	}
	return out, nil
}

func (m *mockUserRepo) GetStats(_ context.Context) (*model.UserStats, error) {
	return &model.UserStats{TotalUsers: len(m.users)}, nil
}

func (m *mockUserRepo) UpdatePassword(_ context.Context, id, hash string) error {
	u, ok := m.users[id]
	if !ok {
		return errors.New("not found")
	}
	u.PasswordHash = hash
	return nil
}

func (m *mockUserRepo) UpdateStatus(_ context.Context, id string, active bool) error {
	u, ok := m.users[id]
	if !ok {
		return errors.New("not found")
	}
	u.Active = active
	return nil
}

type mockUserCache struct{}

func (m *mockUserCache) GetByID(_ context.Context, _ string) (*model.User, error) {
	return nil, errors.New("miss")
}
func (m *mockUserCache) SetByID(_ context.Context, _ *model.User) error { return nil }
func (m *mockUserCache) DeleteByID(_ context.Context, _ string) error   { return nil }

type mockUserPublisher struct{ subjects []string }

func (m *mockUserPublisher) Publish(_ context.Context, subject string, _ interface{}) error {
	m.subjects = append(m.subjects, subject)
	return nil
}

func newUserUC() (*usecase.UserUseCase, *mockUserRepo, *mockUserPublisher) {
	repo := newMockUserRepo()
	pub := &mockUserPublisher{}
	uc := usecase.NewUserUseCase(repo, &mockUserCache{}, pub)
	return uc, repo, pub
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestCreateUser_Success(t *testing.T) {
	uc, _, _ := newUserUC()

	u, err := uc.Create(context.Background(), "test@example.com", "password123", "Test User", "+7777")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == "" {
		t.Error("expected ID to be generated")
	}
	if u.Email != "test@example.com" {
		t.Errorf("got email %q, want test@example.com", u.Email)
	}
	if u.PasswordHash == "password123" {
		t.Error("password should be hashed, not stored as plain text")
	}
}

func TestCreateUser_MissingEmail(t *testing.T) {
	uc, _, _ := newUserUC()

	_, err := uc.Create(context.Background(), "", "password123", "Test User", "")
	if err == nil {
		t.Error("expected error for missing email")
	}
}

func TestCreateUser_MissingPassword(t *testing.T) {
	uc, _, _ := newUserUC()

	_, err := uc.Create(context.Background(), "test@example.com", "", "Test User", "")
	if err == nil {
		t.Error("expected error for missing password")
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	uc, _, _ := newUserUC()

	uc.Create(context.Background(), "dup@example.com", "pass1", "User One", "")

	_, err := uc.Create(context.Background(), "dup@example.com", "pass2", "User Two", "")
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}

func TestAuthenticate_Success(t *testing.T) {
	uc, _, _ := newUserUC()

	uc.Create(context.Background(), "auth@example.com", "secret123", "Auth User", "")

	u, err := uc.Authenticate(context.Background(), "auth@example.com", "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Email != "auth@example.com" {
		t.Errorf("got email %q, want auth@example.com", u.Email)
	}
}

func TestAuthenticate_WrongPassword(t *testing.T) {
	uc, _, _ := newUserUC()

	uc.Create(context.Background(), "auth@example.com", "correct_pass", "Auth User", "")

	_, err := uc.Authenticate(context.Background(), "auth@example.com", "wrong_pass")
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestAuthenticate_UserNotFound(t *testing.T) {
	uc, _, _ := newUserUC()

	_, err := uc.Authenticate(context.Background(), "nobody@example.com", "pass")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestAuthenticate_EmptyCredentials(t *testing.T) {
	uc, _, _ := newUserUC()

	_, err := uc.Authenticate(context.Background(), "", "")
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestAuthenticate_DeactivatedUser(t *testing.T) {
	uc, repo, _ := newUserUC()

	u, _ := uc.Create(context.Background(), "inactive@example.com", "pass123", "Inactive", "")
	repo.users[u.ID].Active = false

	_, err := uc.Authenticate(context.Background(), "inactive@example.com", "pass123")
	if err == nil {
		t.Error("expected error for deactivated account")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	uc, _, _ := newUserUC()

	u, _ := uc.Create(context.Background(), "update@example.com", "pass123", "Old Name", "")

	updated, err := uc.Update(context.Background(), u.ID, "New Name", "+9999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.FullName != "New Name" {
		t.Errorf("got name %q, want New Name", updated.FullName)
	}
	if updated.Phone != "+9999" {
		t.Errorf("got phone %q, want +9999", updated.Phone)
	}
}

func TestUpdateUser_EmptyID(t *testing.T) {
	uc, _, _ := newUserUC()

	_, err := uc.Update(context.Background(), "", "New Name", "")
	if err == nil {
		t.Error("expected error for empty id")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	uc, repo, _ := newUserUC()

	u, _ := uc.Create(context.Background(), "delete@example.com", "pass123", "Delete Me", "")

	if err := uc.Delete(context.Background(), u.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := repo.users[u.ID]; ok {
		t.Error("user should have been deleted")
	}
}

func TestChangePassword_Success(t *testing.T) {
	uc, _, _ := newUserUC()

	u, _ := uc.Create(context.Background(), "change@example.com", "oldpass", "Change User", "")

	if err := uc.ChangePassword(context.Background(), u.ID, "oldpass", "newpass"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	authed, err := uc.Authenticate(context.Background(), "change@example.com", "newpass")
	if err != nil {
		t.Fatalf("authentication with new password failed: %v", err)
	}
	if authed.ID != u.ID {
		t.Error("authenticated user id mismatch")
	}
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	uc, _, _ := newUserUC()

	u, _ := uc.Create(context.Background(), "change2@example.com", "correctpass", "User", "")

	err := uc.ChangePassword(context.Background(), u.ID, "wrongpass", "newpass")
	if err == nil {
		t.Error("expected error for wrong old password")
	}
}

func TestUpdateStatus_Deactivate(t *testing.T) {
	uc, repo, _ := newUserUC()

	u, _ := uc.Create(context.Background(), "status@example.com", "pass", "Status User", "")

	_, err := uc.UpdateStatus(context.Background(), u.ID, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.users[u.ID].Active {
		t.Error("user should be deactivated")
	}
}

func TestGetByID_EmptyID(t *testing.T) {
	uc, _, _ := newUserUC()

	_, err := uc.GetByID(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty id")
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	uc, _, _ := newUserUC()

	_, err := uc.Search(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty query")
	}
}
