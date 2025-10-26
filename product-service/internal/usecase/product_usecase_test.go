package usecase_test

import (
	"context"
	"errors"
	"testing"

	"sneaker-store/product-service/internal/model"
	"sneaker-store/product-service/internal/usecase"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockProductRepo struct {
	products map[string]*model.Product
	err      error
}

func newMockRepo() *mockProductRepo {
	return &mockProductRepo{products: make(map[string]*model.Product)}
}

func (m *mockProductRepo) Create(_ context.Context, p *model.Product) error {
	if m.err != nil {
		return m.err
	}
	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepo) GetByID(_ context.Context, id string) (*model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	p, ok := m.products[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (m *mockProductRepo) List(_ context.Context) ([]*model.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	out := make([]*model.Product, 0, len(m.products))
	for _, p := range m.products {
		out = append(out, p)
	}
	return out, nil
}

func (m *mockProductRepo) Update(_ context.Context, p *model.Product) error {
	if m.err != nil {
		return m.err
	}
	m.products[p.ID] = p
	return nil
}

func (m *mockProductRepo) Delete(_ context.Context, id string) error {
	if m.err != nil {
		return m.err
	}
	delete(m.products, id)
	return nil
}

func (m *mockProductRepo) Search(_ context.Context, brand string, _, _ float64, _ string) ([]*model.Product, error) {
	var out []*model.Product
	for _, p := range m.products {
		if brand == "" || p.Brand == brand {
			out = append(out, p)
		}
	}
	return out, nil
}

func (m *mockProductRepo) UpdateStock(_ context.Context, id string, delta int) (int, error) {
	p, ok := m.products[id]
	if !ok {
		return 0, errors.New("not found")
	}
	p.Stock += delta
	return p.Stock, nil
}

func (m *mockProductRepo) GetByBrand(_ context.Context, brand string) ([]*model.Product, error) {
	var out []*model.Product
	for _, p := range m.products {
		if p.Brand == brand {
			out = append(out, p)
		}
	}
	return out, nil
}

func (m *mockProductRepo) GetLowStock(_ context.Context, threshold int) ([]*model.Product, error) {
	var out []*model.Product
	for _, p := range m.products {
		if p.Stock <= threshold {
			out = append(out, p)
		}
	}
	return out, nil
}

func (m *mockProductRepo) GetBrands(_ context.Context) ([]string, error) {
	seen := map[string]bool{}
	var brands []string
	for _, p := range m.products {
		if !seen[p.Brand] {
			brands = append(brands, p.Brand)
			seen[p.Brand] = true
		}
	}
	return brands, nil
}

func (m *mockProductRepo) GetStats(_ context.Context) (*model.ProductStats, error) {
	return &model.ProductStats{TotalProducts: len(m.products)}, nil
}

func (m *mockProductRepo) BulkDelete(_ context.Context, ids []string) (int, error) {
	count := 0
	for _, id := range ids {
		if _, ok := m.products[id]; ok {
			delete(m.products, id)
			count++
		}
	}
	return count, nil
}

type mockProductCache struct{}

func (m *mockProductCache) GetByID(_ context.Context, _ string) (*model.Product, error) {
	return nil, errors.New("miss")
}
func (m *mockProductCache) SetByID(_ context.Context, _ *model.Product) error  { return nil }
func (m *mockProductCache) DeleteByID(_ context.Context, _ string) error        { return nil }
func (m *mockProductCache) GetList(_ context.Context) ([]*model.Product, error) { return nil, errors.New("miss") }
func (m *mockProductCache) SetList(_ context.Context, _ []*model.Product) error { return nil }
func (m *mockProductCache) DeleteList(_ context.Context) error                  { return nil }

type mockProductPublisher struct{ published []string }

func (m *mockProductPublisher) Publish(_ context.Context, subject string, _ interface{}) error {
	m.published = append(m.published, subject)
	return nil
}

func newUC() (*usecase.ProductUseCase, *mockProductRepo, *mockProductPublisher) {
	repo := newMockRepo()
	pub := &mockProductPublisher{}
	uc := usecase.NewProductUseCase(repo, &mockProductCache{}, pub)
	return uc, repo, pub
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestCreateProduct_Success(t *testing.T) {
	uc, _, _ := newUC()

	p, err := uc.Create(context.Background(), &model.Product{Name: "Air Max 90", Brand: "Nike", Price: 130})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID == "" {
		t.Error("expected ID to be generated")
	}
	if p.Name != "Air Max 90" {
		t.Errorf("got name %q, want %q", p.Name, "Air Max 90")
	}
}

func TestCreateProduct_MissingName(t *testing.T) {
	uc, _, _ := newUC()

	_, err := uc.Create(context.Background(), &model.Product{Price: 100})
	if err == nil {
		t.Error("expected validation error for missing name")
	}
}

func TestCreateProduct_NegativePrice(t *testing.T) {
	uc, _, _ := newUC()

	_, err := uc.Create(context.Background(), &model.Product{Name: "Test", Price: -10})
	if err == nil {
		t.Error("expected validation error for negative price")
	}
}

func TestGetByID_FromRepo(t *testing.T) {
	uc, repo, _ := newUC()

	created, _ := uc.Create(context.Background(), &model.Product{Name: "Yeezy 350", Brand: "Adidas", Price: 220})
	_ = repo // repo already has the product

	got, err := uc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("got id %q, want %q", got.ID, created.ID)
	}
}

func TestGetByID_EmptyID(t *testing.T) {
	uc, _, _ := newUC()

	_, err := uc.GetByID(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty id")
	}
}

func TestUpdateProduct_Success(t *testing.T) {
	uc, _, _ := newUC()

	created, _ := uc.Create(context.Background(), &model.Product{Name: "Old Name", Brand: "Nike", Price: 100})

	updated, err := uc.Update(context.Background(), &model.Product{ID: created.ID, Name: "New Name", Price: 150})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("got name %q, want %q", updated.Name, "New Name")
	}
	if updated.Price != 150 {
		t.Errorf("got price %v, want 150", updated.Price)
	}
}

func TestUpdateProduct_MissingID(t *testing.T) {
	uc, _, _ := newUC()

	_, err := uc.Update(context.Background(), &model.Product{Name: "No ID"})
	if err == nil {
		t.Error("expected error for missing id")
	}
}

func TestDeleteProduct_Success(t *testing.T) {
	uc, repo, _ := newUC()

	created, _ := uc.Create(context.Background(), &model.Product{Name: "To Delete", Brand: "Nike", Price: 100})

	if err := uc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := repo.products[created.ID]; ok {
		t.Error("product should have been deleted from repo")
	}
}

func TestDeleteProduct_EmptyID(t *testing.T) {
	uc, _, _ := newUC()

	if err := uc.Delete(context.Background(), ""); err == nil {
		t.Error("expected error for empty id")
	}
}

func TestUpdateStock(t *testing.T) {
	uc, _, _ := newUC()

	created, _ := uc.Create(context.Background(), &model.Product{Name: "Jordan 1", Brand: "Jordan", Price: 180, Stock: 10})

	newStock, err := uc.UpdateStock(context.Background(), created.ID, -3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newStock != 7 {
		t.Errorf("got stock %d, want 7", newStock)
	}
}

func TestGetByBrand_Success(t *testing.T) {
	uc, _, _ := newUC()

	uc.Create(context.Background(), &model.Product{Name: "Air Max 90", Brand: "Nike", Price: 130})
	uc.Create(context.Background(), &model.Product{Name: "Air Force 1", Brand: "Nike", Price: 110})
	uc.Create(context.Background(), &model.Product{Name: "Yeezy 350", Brand: "Adidas", Price: 220})

	products, err := uc.GetByBrand(context.Background(), "Nike")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("got %d products for Nike, want 2", len(products))
	}
}

func TestGetByBrand_EmptyBrand(t *testing.T) {
	uc, _, _ := newUC()

	_, err := uc.GetByBrand(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty brand")
	}
}

func TestBulkDelete(t *testing.T) {
	uc, _, _ := newUC()

	p1, _ := uc.Create(context.Background(), &model.Product{Name: "P1", Brand: "Nike", Price: 100})
	p2, _ := uc.Create(context.Background(), &model.Product{Name: "P2", Brand: "Adidas", Price: 110})

	count, err := uc.BulkDelete(context.Background(), []string{p1.ID, p2.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("got count %d, want 2", count)
	}
}

func TestBulkDelete_EmptyIDs(t *testing.T) {
	uc, _, _ := newUC()

	_, err := uc.BulkDelete(context.Background(), []string{})
	if err == nil {
		t.Error("expected error for empty ids")
	}
}
