package customers

import (
	"context"
	"os"
	"testing"

	"shipments/domains/customers/store"
	"shipments/domains/entities"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var mockDB = make(map[string]*store.Customer)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()

	code = m.Run()
}

func TestCreateCustomer(t *testing.T) {
	store := &MockCustomerStore{
		CreateFunc: func(ctx context.Context, c *store.Customer) error {
			c.ID = uuid.NewString()
			mockDB[c.ID] = c
			return nil
		},
	}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Name:  gofakeit.FirstName(),
		Email: gofakeit.Email(),
	}

	err := svc.Create(ctx, c)
	require.NoError(t, err)
	require.NotEmpty(t, c.ID)
}

func TestCreateCustomerNoEmail(t *testing.T) {
	store := &MockCustomerStore{
		CreateFunc: func(ctx context.Context, c *store.Customer) error {
			c.ID = uuid.NewString()
			mockDB[c.ID] = c
			return nil
		},
	}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Name: gofakeit.LastName(),
	}
	err := svc.Create(ctx, c)
	require.EqualError(t, err, "invalid email for customer")
	require.Empty(t, c.ID)
}

func TestCreateCustomerWithNoMockInitialized(t *testing.T) {
	store := &MockCustomerStore{}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Email: gofakeit.Email(),
		Name:  gofakeit.LastName(),
	}
	err := svc.Create(ctx, c)
	require.EqualError(t, err, "mock not initialized")
	require.Empty(t, c.ID)
}

func TestFindCustomer(t *testing.T) {
	store := &MockCustomerStore{
		CreateFunc: func(ctx context.Context, c *store.Customer) error {
			c.ID = uuid.NewString()
			mockDB[c.ID] = c
			return nil
		},
		FindFunc: func(ctx context.Context, id string) (*store.Customer, error) {
			v, ok := mockDB[id]
			if !ok {
				return nil, gorm.ErrRecordNotFound
			}
			return v, nil
		},
	}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Name:  gofakeit.FirstName(),
		Email: gofakeit.Email(),
	}

	err := svc.Create(ctx, c)
	require.NoError(t, err)
	require.NotEmpty(t, c.ID)

	res, err := svc.Find(ctx, c.ID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.Name, c.Name)
	assert.Equal(t, res.Email, c.Email)
}

func TestFindCustomersStoreError(t *testing.T) {
	store := &MockCustomerStore{
		FindFunc: func(ctx context.Context, id string) (*store.Customer, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Name:  gofakeit.FirstName(),
		Email: gofakeit.Email(),
	}

	res, err := svc.Find(ctx, c.ID)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindCustomerByEmail(t *testing.T) {
	store := &MockCustomerStore{
		CreateFunc: func(ctx context.Context, c *store.Customer) error {
			c.ID = uuid.NewString()
			mockDB[c.ID] = c
			return nil
		},
		FindByEmailFunc: func(ctx context.Context, email string) (*store.Customer, error) {
			for _, v := range mockDB {
				if v.Email == email {
					return v, nil
				}
			}
			return nil, gorm.ErrRecordNotFound
		},
	}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Name:  gofakeit.FirstName(),
		Email: gofakeit.Email(),
	}

	err := svc.Create(ctx, c)
	require.NoError(t, err)
	require.NotEmpty(t, c.ID)

	res, err := svc.FindByEmail(ctx, c.Email)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.Name, c.Name)
	assert.Equal(t, res.Email, c.Email)
}

func TestFindByEmailCustomersStoreError(t *testing.T) {
	store := &MockCustomerStore{
		FindByEmailFunc: func(ctx context.Context, id string) (*store.Customer, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Name:  gofakeit.FirstName(),
		Email: gofakeit.Email(),
	}

	res, err := svc.FindByEmail(ctx, c.Email)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestUpdatCustomerRecord(t *testing.T) {
	store := &MockCustomerStore{
		CreateFunc: func(ctx context.Context, c *store.Customer) error {
			c.ID = uuid.NewString()
			mockDB[c.ID] = c
			return nil
		},
		UpdateFunc: func(ctx context.Context, c *store.Customer) error {
			v, ok := mockDB[c.ID]
			if !ok {
				return gorm.ErrRecordNotFound
			}
			v.Name = c.Name
			mockDB[v.ID] = v
			return nil
		},
		FindFunc: func(ctx context.Context, id string) (*store.Customer, error) {
			return mockDB[id], nil
		},
	}
	ctx := context.Background()
	svc := New(store)
	c := &entities.Customer{
		Name:  gofakeit.FirstName(),
		Email: gofakeit.Email(),
	}

	err := svc.Create(ctx, c)
	require.NoError(t, err)
	require.NotEmpty(t, c.ID)

	newName := gofakeit.LastName()
	c.Name = newName

	err = svc.Update(ctx, c)
	require.NoError(t, err)

	oldID := c.ID
	c.ID = uuid.NewString()
	require.EqualError(t, svc.Update(ctx, c), gorm.ErrRecordNotFound.Error())

	res, err := svc.Find(ctx, oldID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.Email, c.Email)
	assert.Equal(t, newName, res.Name)
}

func TestUpdateCustomerNoMock(t *testing.T) {
	store := &MockCustomerStore{}
	svc := New(store)

	c := &entities.Customer{
		Email: gofakeit.Email(),
		Name:  gofakeit.Name(),
	}
	ctx := context.Background()
	err := svc.Update(ctx, c)
	require.EqualError(t, err, errMockNotInitialized.Error())
}

func TestFindOrCreate(t *testing.T) {
	ctx := context.Background()
	localMockDB := make(map[string]*store.Customer)
	store := &MockCustomerStore{
		CreateFunc: func(ctx context.Context, c *store.Customer) error {
			c.ID = uuid.NewString()
			localMockDB[c.ID] = c
			return nil
		},
		FindByEmailFunc: func(ctx context.Context, email string) (*store.Customer, error) {
			for _, v := range localMockDB {
				if v.Email == email {
					return v, nil
				}
			}
			return nil, gorm.ErrRecordNotFound

		},
	}
	svc := New(store)
	c := &entities.Customer{
		Email: gofakeit.Email(),
		Name:  gofakeit.Name(),
	}
	err := svc.FindOrCreate(ctx, c)
	require.NoError(t, err)
	require.NotEmpty(t, c.ID)
	assert.Len(t, localMockDB, 1)

	err = svc.FindOrCreate(ctx, c)
	require.NoError(t, err)
	require.NotEmpty(t, c.ID)

	assert.Len(t, localMockDB, 1) //no new record should be
}
