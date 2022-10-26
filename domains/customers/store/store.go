package store

import (
	"context"

	"github.com/google/uuid"
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
	c.ID = uuid.NewString()
	return cs.DB.Create(c).Error
}

// Find implements finds a user by the given ID
func (cs *CustomerStore) Find(ctx context.Context, id string) (*Customer, error) {
	var c *Customer
	err := cs.DB.Where("id = ?", id).First(&c).Error
	return c, err
}

// FindByEmail finds a user by the email
func (cs *CustomerStore) FindByEmail(ctx context.Context, email string) (*Customer, error) {
	var c *Customer
	err := cs.DB.Where("email = ?", email).First(&c).Error
	return c, err
}

// Update updates a customer's information, name mainly.
func (cs *CustomerStore) Update(ctx context.Context, c *Customer) error {
	return cs.DB.Table("customers").Where("id = ?", c.ID).Update("name", c.Name).Error
}
