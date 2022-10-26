package store

import "context"

// IShipmentStore interface for the shipment store
type IShipmentStore interface {
	Create(ctx context.Context, s *Shipment) error
	Find(ctx context.Context, id string) (*Shipment, error)
	FindCustomerShipments(ctx context.Context, customerID string) ([]*Shipment, error)
}
