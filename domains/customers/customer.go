package customers

import (
	"context"
	"errors"
	"fmt"

	"shipments/domains/customers/store"
	"shipments/domains/entities"

	"gorm.io/gorm"
)

type CustomerService struct {
	store store.ICustomerStore
}

// FindOrCreate checks if the customer exists by email and create if it doesnt exist
func (cs *CustomerService) FindOrCreate(ctx context.Context, c *entities.Customer) error {
	res, err := cs.FindByEmail(ctx, c.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if res != nil {
		c.ID = res.ID
		return nil
	}
	return cs.Create(ctx, c)
}

// New returns an instance of customer service
func New(store store.ICustomerStore) ICustomerService {
	return &CustomerService{
		store: store,
	}
}

// Create creates a new customer record
func (cs *CustomerService) Create(ctx context.Context, c *entities.Customer) error {
	if c.Email == "" {
		return fmt.Errorf("invalid email for customer")
	}
	dbEnt := c.ToDBCustomer()
	if err := cs.store.Create(ctx, dbEnt); err != nil {
		return err
	}
	c.ID = dbEnt.ID
	return nil
}

// Find finds and return customer/error with the given ID
func (cs *CustomerService) Find(ctx context.Context, id string) (*entities.Customer, error) {
	res, err := cs.store.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return entities.FromCustomerDBEntities(res), nil
}

// FindByEmail finds and return customer/error with the given email
func (cs *CustomerService) FindByEmail(ctx context.Context, email string) (*entities.Customer, error) {
	res, err := cs.store.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return entities.FromCustomerDBEntities(res), nil
}

// Update updates a customer's record (name)
func (cs *CustomerService) Update(ctx context.Context, c *entities.Customer) error {
	return cs.store.Update(ctx, c.ToDBCustomer())
}
