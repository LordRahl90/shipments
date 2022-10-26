package store

import "context"

// ICustomerStore interface for the customer store
type ICustomerStore interface {
	Create(ctx context.Context, c *Customer) error
	Find(ctx context.Context, id string) (*Customer, error)
	FindByEmail(ctx context.Context, email string) (*Customer, error)
	Update(ctx context.Context, c *Customer) error
}
