package entities

import "shipments/domains/customers/store"

// Customer contains basic customer information
type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToDBCustomer converts service entities to db entiites
func (c *Customer) ToDBCustomer() *store.Customer {
	return &store.Customer{
		ID:    c.ID,
		Email: c.Email,
		Name:  c.Name,
	}
}

// FromCustomerDBEntities converts db entity to service entity
func FromCustomerDBEntities(c *store.Customer) *Customer {
	return &Customer{
		ID:    c.ID,
		Name:  c.Name,
		Email: c.Email,
	}
}
