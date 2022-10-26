package shipments

import (
	"context"
	"fmt"
	"os"
	"testing"

	"shipments/domains/entities"
	"shipments/domains/shipments/store"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	mockDB = make(map[string]*store.Shipment)
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()

	code = m.Run()
}

func TestCreateShipment(t *testing.T) {
	store := MockShipmentStore{
		CreateFunc: func(ctx context.Context, s *store.Shipment) error {
			s.ID = uuid.NewString()
			mockDB[s.ID] = s
			return nil
		},
	}
	svc := New(&store)
	ctx := context.Background()
	customerID := uuid.NewString()

	s := &entities.Shipment{
		CustomerID:  customerID,
		Origin:      "us",
		Destination: "se",
		Weight:      45.0,
	}
	err := svc.Create(ctx, s)
	require.NoError(t, err)
	require.NotEmpty(t, s.ID)
	assert.Equal(t, float64(1250), s.Price)
}

func TestCreateShipmentStoreError(t *testing.T) {
	store := MockShipmentStore{
		CreateFunc: func(ctx context.Context, s *store.Shipment) error {
			return fmt.Errorf("undefined error")
		},
	}
	svc := New(&store)
	ctx := context.Background()
	customerID := uuid.NewString()

	s := &entities.Shipment{
		CustomerID:  customerID,
		Origin:      "us",
		Destination: "se",
		Weight:      45.0,
	}
	err := svc.Create(ctx, s)
	require.EqualError(t, err, "undefined error")
	require.Empty(t, s.ID)
}

func TestCreateShipmentInvalidMock(t *testing.T) {
	store := MockShipmentStore{}
	svc := New(&store)
	ctx := context.Background()
	customerID := uuid.NewString()

	s := &entities.Shipment{
		CustomerID:  customerID,
		Origin:      "us",
		Destination: "se",
		Weight:      45.0,
	}
	err := svc.Create(ctx, s)
	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Empty(t, s.ID)
}

func TestFindShipment(t *testing.T) {
	store := MockShipmentStore{
		CreateFunc: func(ctx context.Context, s *store.Shipment) error {
			s.ID = uuid.NewString()
			mockDB[s.ID] = s
			return nil
		},
		FindFunc: func(ctx context.Context, id string) (*store.Shipment, error) {
			return mockDB[id], nil
		},
	}
	svc := New(&store)
	ctx := context.Background()
	customerID := uuid.NewString()

	s := &entities.Shipment{
		CustomerID:  customerID,
		Origin:      "us",
		Destination: "se",
		Weight:      45.0,
	}
	err := svc.Create(ctx, s)
	require.NoError(t, err)
	require.NotEmpty(t, s.ID)
	assert.Equal(t, float64(1250), s.Price)

	res, err := svc.Find(ctx, s.ID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, res.ID, s.ID)
	assert.Equal(t, res.CustomerID, s.CustomerID)
	assert.Equal(t, res.Origin, s.Origin)
	assert.Equal(t, res.Destination, s.Destination)
	assert.Equal(t, res.Weight, s.Weight)
	assert.Equal(t, res.Price, s.Price)
}

func TestFindShipmentWithStoreError(t *testing.T) {
	store := MockShipmentStore{
		FindFunc: func(ctx context.Context, id string) (*store.Shipment, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(&store)
	ctx := context.Background()
	res, err := svc.Find(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindShipmentWithInvalidMock(t *testing.T) {
	store := MockShipmentStore{}
	svc := New(&store)
	ctx := context.Background()
	res, err := svc.Find(ctx, uuid.NewString())
	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Empty(t, res)
}

func TestFindCustomerShipments(t *testing.T) {
	store := MockShipmentStore{
		CreateFunc: func(ctx context.Context, s *store.Shipment) error {
			s.ID = uuid.NewString()
			mockDB[s.ID] = s
			return nil
		},
		FindCustomerShipmentsFunc: func(ctx context.Context, customerID string) ([]*store.Shipment, error) {
			var res []*store.Shipment
			for _, v := range mockDB {
				if v.CustomerID == customerID {
					res = append(res, v)
				}
			}
			return res, nil
		},
	}
	svc := New(&store)
	ctx := context.Background()
	customerID := uuid.NewString()

	// create for a random user
	s := &entities.Shipment{
		CustomerID:  uuid.NewString(),
		Origin:      "us",
		Destination: "se",
		Weight:      float64(gofakeit.Number(20, 100)),
	}
	err := svc.Create(ctx, s)
	require.NoError(t, err)

	for i := 1; i <= 4; i++ {
		s := &entities.Shipment{
			CustomerID:  customerID,
			Origin:      "us",
			Destination: "se",
			Weight:      float64(gofakeit.Number(20, 100)),
		}
		err := svc.Create(ctx, s)
		require.NoError(t, err)
	}
	res, err := svc.FindCustomerShipments(ctx, customerID)

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, 4)
}

func TestFindCustomerShipmentWithDBError(t *testing.T) {
	store := MockShipmentStore{
		FindCustomerShipmentsFunc: func(ctx context.Context, customerID string) ([]*store.Shipment, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := New(&store)
	ctx := context.Background()
	res, err := svc.FindCustomerShipments(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindCustomerShipmentWithnilMock(t *testing.T) {
	store := MockShipmentStore{}
	svc := New(&store)
	ctx := context.Background()

	res, err := svc.FindCustomerShipments(ctx, uuid.NewString())
	require.EqualError(t, err, errMockNotInitialized.Error())
	require.Empty(t, res)
}
