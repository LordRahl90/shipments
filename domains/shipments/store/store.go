package store

import (
	"context"

	"shipments/domains/tracing"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

var (
	_ IShipmentStore = (*ShipmentStore)(nil)
)

// ShipmentStore implements the IShipmentStore
type ShipmentStore struct {
	DB *gorm.DB
}

// New returns a new instance of the shipment store interface
func New(db *gorm.DB) (IShipmentStore, error) {
	if err := db.AutoMigrate(&Shipment{}); err != nil {
		return nil, err
	}
	return &ShipmentStore{DB: db}, nil
}

// Create implements IShipmentStore
func (ss *ShipmentStore) Create(ctx context.Context, s *Shipment) error {
	ctx, span := tracing.Tracer().Start(ctx, "db:CreateShipment")
	defer span.End()
	s.ID = uuid.NewString()
	return ss.DB.WithContext(ctx).Create(s).Error
}

// Find implements IShipmentStore
func (ss *ShipmentStore) Find(ctx context.Context, id string) (*Shipment, error) {
	ctx, span := tracing.Tracer().Start(ctx, "db:FindShipment")
	span.SetAttributes(attribute.KeyValue{
		Key:   "shipment_id",
		Value: attribute.StringValue(id),
	})
	defer span.End()
	var s *Shipment
	err := ss.DB.WithContext(ctx).Where("id = ?", id).First(&s).Error
	return s, err
}

// FindCustomerShipments returns all the shipments for a customer.
// TODO: This should be paginated for scale
func (ss *ShipmentStore) FindCustomerShipments(ctx context.Context, customerID string) ([]*Shipment, error) {
	ctx, span := tracing.Tracer().Start(ctx, "db:FindCustomerShipments")
	span.SetAttributes(attribute.KeyValue{
		Key:   "customer_id",
		Value: attribute.StringValue(customerID),
	})
	defer span.End()
	var s []*Shipment
	err := ss.DB.WithContext(ctx).Where("customer_id = ?", customerID).Find(&s).Error
	return s, err
}
