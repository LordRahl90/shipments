package customers

import (
	"context"
	"errors"
	"shipments/domains/customers/store"
)

var (
	errMockNotInitialized = errors.New("mock not initialized")

	_ store.ICustomerStore = (*MockCustomerStore)(nil)
)

// MockCustomerStore mock service for the store
type MockCustomerStore struct {
	CreateFunc      func(ctx context.Context, c *store.Customer) error
	FindFunc        func(ctx context.Context, id string) (*store.Customer, error)
	FindByEmailFunc func(ctx context.Context, email string) (*store.Customer, error)
	UpdateFunc      func(ctx context.Context, c *store.Customer) error
}

// Create implements store.ICustomerStore
func (m *MockCustomerStore) Create(ctx context.Context, c *store.Customer) error {
	if m.CreateFunc == nil {
		return errMockNotInitialized
	}
	return m.CreateFunc(ctx, c)
}

// Find implements store.ICustomerStore
func (m *MockCustomerStore) Find(ctx context.Context, id string) (*store.Customer, error) {
	if m.FindFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.FindFunc(ctx, id)
}

// FindByEmail implements store.ICustomerStore
func (m *MockCustomerStore) FindByEmail(ctx context.Context, email string) (*store.Customer, error) {
	if m.FindByEmailFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.FindByEmailFunc(ctx, email)
}

// Update implements store.ICustomerStore
func (m *MockCustomerStore) Update(ctx context.Context, c *store.Customer) error {
	if m.UpdateFunc == nil {
		return errMockNotInitialized
	}
	return m.UpdateFunc(ctx, c)
}
