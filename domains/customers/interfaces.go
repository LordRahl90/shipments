package customers

import (
	"context"
	"shipments/domains/entities"
)

// ICustomerService interface for customer service implementation
type ICustomerService interface {
	Create(ctx context.Context, c *entities.Customer) error
	FindOrCreate(ctx context.Context, c *entities.Customer) error
	Find(ctx context.Context, id string) (*entities.Customer, error)
	FindByEmail(ctx context.Context, email string) (*entities.Customer, error)
	Update(ctx context.Context, c *entities.Customer) error
}
