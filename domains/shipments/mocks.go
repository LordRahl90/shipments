package shipments

import (
	"context"
	"errors"

	"shipments/domains/shipments/store"
)

var (
	errMockNotInitialized = errors.New("mock not initialized")

	_ store.IShipmentStore = (*MockShipmentStore)(nil)
)

// MockShipmentStore mock for shipment store
type MockShipmentStore struct {
	CreateFunc                func(ctx context.Context, s *store.Shipment) error
	FindFunc                  func(ctx context.Context, id string) (*store.Shipment, error)
	FindCustomerShipmentsFunc func(ctx context.Context, customerID string) ([]*store.Shipment, error)
}

// Create implements store.IShipmentStore
func (m *MockShipmentStore) Create(ctx context.Context, s *store.Shipment) error {
	if m.CreateFunc == nil {
		return errMockNotInitialized
	}
	return m.CreateFunc(ctx, s)
}

// Find implements store.IShipmentStore
func (m *MockShipmentStore) Find(ctx context.Context, id string) (*store.Shipment, error) {
	if m.FindFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.FindFunc(ctx, id)
}

// FindCustomerShipments implements store.IShipmentStore
func (m *MockShipmentStore) FindCustomerShipments(ctx context.Context, customerID string) ([]*store.Shipment, error) {
	if m.FindCustomerShipmentsFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.FindCustomerShipmentsFunc(ctx, customerID)
}
