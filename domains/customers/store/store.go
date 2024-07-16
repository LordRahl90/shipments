package store

import (
	"context"

	"shipments/domains/tracing"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
)

var (
	_ ICustomerStore = (*CustomerStore)(nil)
)

// CustomerStore implements the ICustomerStore interface
type CustomerStore struct {
	DB *gorm.DB
}

// New returns a new instance of CustomerStore
func New(db *gorm.DB) (ICustomerStore, error) {
	if err := db.AutoMigrate(&Customer{}); err != nil {
		return nil, err
	}
	return &CustomerStore{
		DB: db,
	}, nil
}

// Create implements ICustomerStore
func (cs *CustomerStore) Create(ctx context.Context, c *Customer) error {
	ctx, span := tracing.Tracer().Start(ctx, "db:CreateCustomer")
	defer span.End()

	c.ID = uuid.NewString()
	return cs.DB.WithContext(ctx).Create(c).Error
}

// Find implements finds a user by the given ID
func (cs *CustomerStore) Find(ctx context.Context, id string) (*Customer, error) {
	ctx, span := tracing.Tracer().Start(ctx, "db:FindCustomer")
	span.SetAttributes(attribute.KeyValue{
		Key:   "customer_id",
		Value: attribute.StringValue(id),
	})
	defer span.End()
	var c *Customer
	err := cs.DB.WithContext(ctx).Where("id = ?", id).First(&c).Error
	return c, err
}

// FindByEmail finds a user by the email
func (cs *CustomerStore) FindByEmail(ctx context.Context, email string) (*Customer, error) {
	ctx, span := tracing.Tracer().Start(ctx, "db:FindCustomerByEmail")
	span.SetAttributes(attribute.KeyValue{
		Key:   "customer_email",
		Value: attribute.StringValue(email),
	})
	defer span.End()

	var c *Customer
	err := cs.DB.WithContext(ctx).Where("email = ?", email).First(&c).Error
	return c, err
}

// Update updates a customer's information, name mainly.
func (cs *CustomerStore) Update(ctx context.Context, c *Customer) error {
	ctx, span := tracing.Tracer().Start(ctx, "db:UpdateCustomer")
	defer span.End()
	return cs.DB.WithContext(ctx).Table("customers").Where("id = ?", c.ID).Update("name", c.Name).Error
}
