package store

import (
	"context"
	"log"
	"os"
	"testing"

	"shipments/testhelpers"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	storage ICustomerStore

	container = testhelpers.GetMySQLContainer(context.TODO())
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		if err := container.Terminate(context.TODO()); err != nil {
			log.Fatalf("cannot stop container %s", err)
		}
		os.Exit(code)
	}()
	d, err := setupTestDB(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	db = d
	s, err := New(db)
	if err != nil {
		log.Fatal(err)
	}
	storage = s
	code = m.Run()
}

func TestCreateCustomer(t *testing.T) {
	ctx := context.Background()
	c := newCustomer(t)
	st, err := New(db)
	require.NoError(t, err)
	require.NoError(t, st.Create(ctx, c))
	require.NotEmpty(t, c.ID)
}

func TestFindCustomer(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	cust := make(map[string]*Customer)
	t.Cleanup(func() {
		db.Exec("DELETE FROM customers WHERE id IN (?)", ids)
	})
	for i := 1; i <= 3; i++ {
		c := newCustomer(t)
		require.NoError(t, storage.Create(ctx, c))
		cust[c.ID] = c
		ids = append(ids, c.ID)
	}

	id := ids[1]

	res, err := storage.Find(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, cust[id].Name, res.Name)
	assert.Equal(t, cust[id].Email, res.Email)
}

func TestFindWithInvalidID(t *testing.T) {
	ctx := context.Background()
	res, err := storage.Find(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func TestFindCustomerByEmail(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	cust := make(map[string]*Customer)
	t.Cleanup(func() {
		db.Exec("DELETE FROM customers WHERE id IN (?)", ids)
	})
	for i := 1; i <= 3; i++ {
		c := newCustomer(t)
		require.NoError(t, storage.Create(ctx, c))
		cust[c.ID] = c
		ids = append(ids, c.ID)
	}

	id := ids[1]
	email := cust[id].Email

	res, err := storage.FindByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, cust[id].Name, res.Name)
	assert.Equal(t, cust[id].Email, res.Email)
}

func TestUpdateUserName(t *testing.T) {
	ctx := context.Background()
	c := newCustomer(t)
	require.NoError(t, storage.Create(ctx, c))
	require.NotEmpty(t, c.ID)

	newName := gofakeit.FirstName() + " " + gofakeit.LastName()
	c.Name = newName
	err := storage.Update(ctx, c)
	require.NoError(t, err)

	res, err := storage.Find(ctx, c.ID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, c.Email, res.Email)
	assert.Equal(t, newName, res.Name)
}

func newCustomer(t *testing.T) *Customer {
	t.Helper()
	return &Customer{
		Name:  gofakeit.FirstName() + " " + gofakeit.LastName(),
		Email: gofakeit.Email(),
	}
}

func setupTestDB(ctx context.Context) (*gorm.DB, error) {
	return testhelpers.SetupTestDB(ctx, container)
}

func cleanup() {
	db.Exec("DELETE FROM customers")
}
