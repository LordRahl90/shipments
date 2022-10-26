package store

import (
	"context"

	"github.com/google/uuid"
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
	s.ID = uuid.NewString()
	return ss.DB.Create(s).Error
}

// Details implements IShipmentStore
func (ss *ShipmentStore) Find(ctx context.Context, id string) (*Shipment, error) {
	var s *Shipment
	err := ss.DB.Where("id = ?", id).First(&s).Error
	return s, err
}

// FindCustomerShipments returns all the shipments for a customer.
// TODO: This should be paginated for scale
func (ss *ShipmentStore) FindCustomerShipments(ctx context.Context, customerID string) ([]*Shipment, error) {
	var s []*Shipment
	err := ss.DB.Where("customer_id = ?", customerID).Find(&s).Error
	return s, err
}
