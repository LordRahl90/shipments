package shipments

import (
	"context"
	"shipments/domains/entities"
)

// IShipmentService service to handle shipments
type IShipmentService interface {
	Create(ctx context.Context, s *entities.Shipment) error
	Find(ctx context.Context, id string) (*entities.Shipment, error)
	FindCustomerShipments(ctx context.Context, customerID string) ([]*entities.Shipment, error)
}
