package store

import (
	"context"
	"log"
	"os"
	"testing"

	"shipments/testhelpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	store IShipmentStore

	container = testhelpers.GetMySQLContainer(context.TODO())
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		if err := container.Terminate(context.TODO()); err != nil {
			log.Fatal(err)
		}
		os.Exit(code)
	}()
	d, err := setupTestDB()
	if err != nil {
		panic(err)
	}
	db = d
	s, err := New(db)
	if err != nil {
		panic(err)
	}
	store = s
	code = m.Run()
}

func TestCreateNewShipment(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.NewString()
	s := newShipment(t, customerID)
	err := store.Create(ctx, s)
	require.NoError(t, err)
}

func TestShipmentDetails(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.NewString()
	s := newShipment(t, customerID)
	err := store.Create(ctx, s)
	require.NoError(t, err)

	res, err := store.Find(ctx, s.ID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, customerID, res.CustomerID)
	assert.Equal(t, s.Origin, res.Origin)
	assert.Equal(t, s.Destination, res.Destination)
	assert.Equal(t, s.Weight, res.Weight)
	assert.Equal(t, s.Price, res.Price)
}

func TestFindNonExistingShipment(t *testing.T) {
	ctx := context.Background()
	res, err := store.Find(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindCustomerShipments(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.NewString()
	for i := 0; i < 5; i++ {
		s := newShipment(t, customerID)
		require.NoError(t, store.Create(ctx, s))
	}

	res, err := store.FindCustomerShipments(ctx, customerID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, 5)
}

func TestFindNonExistentCustomerShipping(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.NewString()
	res, err := store.FindCustomerShipments(ctx, customerID)
	require.NoError(t, err)
	require.Empty(t, res)
	assert.Len(t, res, 0)
}

func newShipment(t *testing.T, customerID string) *Shipment {
	t.Helper()
	return &Shipment{
		CustomerID:  customerID,
		Origin:      "se",
		Destination: "dk",
		Weight:      250,
		Price:       2500, // doesn't matter, just to test
	}
}

func setupTestDB() (*gorm.DB, error) {
	return testhelpers.SetupTestDB(context.TODO(), container)
}

func cleanup() {
	db.Exec("DELETE FROM shipments")
}
