package shipments

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"shipments/domains/tracing"

	"shipments/domains/core"
	"shipments/domains/entities"
	"shipments/domains/shipments/store"
)

var _ IShipmentService = (*ShipmentService)(nil)

// ShipmentService implements IShipmentService
type ShipmentService struct {
	store store.IShipmentStore
}

// New returns a new implementation of IShipmentService
func New(store store.IShipmentStore) IShipmentService {
	return &ShipmentService{
		store: store,
	}
}

// Create creates a new shipment record
func (ss *ShipmentService) Create(ctx context.Context, s *entities.Shipment) error {
	ctx, span := tracing.Tracer().Start(ctx, "svc:CreateShipment")
	defer span.End()
	p, err := core.PriceFromSize(ctx, s.Weight, s.Origin, s.Destination)
	if err != nil {
		return err
	}
	s.Price = p
	dbEnt := s.ToDBShipment()
	if err := ss.store.Create(ctx, dbEnt); err != nil {
		return err
	}
	s.ID = dbEnt.ID
	return nil
}

// Find gets a single record from the store with the given ID
func (ss *ShipmentService) Find(ctx context.Context, id string) (*entities.Shipment, error) {
	ctx, span := tracing.Tracer().Start(ctx, "svc:FindShipment")
	span.SetAttributes(attribute.KeyValue{
		Key:   "shipment_id",
		Value: attribute.StringValue(id),
	})
	res, err := ss.store.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return entities.FromDBShipmentEntity(res), nil
}

// FindCustomerShipments finds all the customer records from the store
func (ss *ShipmentService) FindCustomerShipments(ctx context.Context, customerID string) ([]*entities.Shipment, error) {
	ctx, span := tracing.Tracer().Start(ctx, "svc:FindCustomerShipments")
	span.SetAttributes(attribute.KeyValue{
		Key:   "customer_id",
		Value: attribute.StringValue(customerID),
	})
	var result []*entities.Shipment

	res, err := ss.store.FindCustomerShipments(ctx, customerID)
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		result = append(result, entities.FromDBShipmentEntity(v))
	}

	return result, nil
}
