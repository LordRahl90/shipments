package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"shipments/domains/entities"
	"shipments/requests"
	"shipments/responses"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	server *Server
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()

	d, err := setupTestDB()
	if err != nil {
		log.Fatal(err)
	}
	db = d
	svr, err := New(db)
	if err != nil {
		panic(err)
	}
	server = svr
	code = m.Run()
}

func TestNewShipment(t *testing.T) {
	shipment := &requests.Shipment{
		Name:        gofakeit.Name(),
		Email:       gofakeit.Email(),
		Origin:      "us",
		Destination: "se",
		Weight:      45.0,
	}

	b, err := json.Marshal(shipment)
	require.NoError(t, err)
	require.NotNil(t, b)

	w := handleRequest(t, http.MethodPost, "/new", b)
	require.Equal(t, http.StatusCreated, w.Code)

	var res responses.Shipment
	err = json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)

	ctx := context.Background()

	cust, err := customerService.Find(ctx, res.CustomerID)
	require.NoError(t, err)
	require.NotEmpty(t, cust)

	assert.Equal(t, cust.Name, shipment.Name)
	assert.Equal(t, cust.Email, shipment.Email)

	shipmentRec, err := shipmentService.Find(ctx, res.Reference)
	require.NoError(t, err)
	require.NotEmpty(t, shipmentRec)

	assert.Equal(t, shipmentRec.Origin, shipment.Origin)
	assert.Equal(t, shipmentRec.Destination, shipment.Destination)
	assert.Equal(t, shipmentRec.Weight, shipment.Weight)
	assert.Equal(t, float64(1250), shipmentRec.Price)
	assert.Equal(t, cust.ID, shipmentRec.CustomerID)
}

func TestCreateWithInvalidJSON(t *testing.T) {
	b := []byte(`
	{
		"name": "Bart Beatty",
		"email": "cordiajacobi@carroll.net",
		"origin": "us",
		"destination": "se",
		"weight": 45,
	}
	`)
	w := handleRequest(t, http.MethodPost, "/new", b)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWithInvalidWeight(t *testing.T) {
	b := []byte(`
	{
		"name": "Bart Beatty",
		"email": "cordiajacobi@carroll.net",
		"origin": "us",
		"destination": "se",
		"weight": 4500000
	}
	`)
	w := handleRequest(t, http.MethodPost, "/new", b)
	require.Equal(t, http.StatusBadRequest, w.Code)
	exp := `{"error":"weight price not defined unknown","success":false}`
	assert.Equal(t, exp, w.Body.String())
}

func TestGetShipmentHistory(t *testing.T) {
	email := gofakeit.Email()

	for i := 1; i <= 5; i++ {
		shipment := &requests.Shipment{
			Name:        gofakeit.Name(),
			Email:       email,
			Origin:      "us",
			Destination: "se",
			Weight:      45.0,
		}

		b, err := json.Marshal(shipment)
		require.NoError(t, err)
		require.NotNil(t, b)

		w := handleRequest(t, http.MethodPost, "/new", b)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	w := handleRequest(t, http.MethodGet, "/history/"+email, nil)
	require.Equal(t, http.StatusOK, w.Code)

	var shipments []*entities.Shipment
	err := json.Unmarshal(w.Body.Bytes(), &shipments)
	require.NoError(t, err)
	require.Len(t, shipments, 5)
}

func TestGetHistoryNonExistentService(t *testing.T) {
	email := gofakeit.Email()
	w := handleRequest(t, http.MethodGet, "/history/"+email, nil)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetHistoryEmptyEmail(t *testing.T) {
	w := handleRequest(t, http.MethodGet, "/history/", nil)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetHistoryWithUserWithNoShipment(t *testing.T) {
	ctx := context.Background()
	customer := &entities.Customer{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	}

	err := customerService.Create(ctx, customer)
	require.NoError(t, err)

	w := handleRequest(t, http.MethodGet, "/history/"+customer.Email, nil)
	require.Equal(t, http.StatusOK, w.Code)

	var shipments []*entities.Shipment
	err = json.Unmarshal(w.Body.Bytes(), &shipments)
	require.NoError(t, err)
	require.Len(t, shipments, 0)
}

func TestPricingDetails(t *testing.T) {
	origin, dest := "us", "se"
	weight := 45.0
	path := fmt.Sprintf("/pricing?origin=%s&destination=%s&weight=%.2f", origin, dest, weight)

	w := handleRequest(t, http.MethodGet, path, nil)
	require.Equal(t, http.StatusOK, w.Code)
	var res responses.Pricing

	err := json.Unmarshal(w.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, origin, res.Origin)
	assert.Equal(t, dest, res.Destination)
	assert.Equal(t, "large", res.WeightCategory)
	assert.Equal(t, float64(weight), res.Weight)
	assert.Equal(t, float64(1250), res.Price)
}

func TestPricingDetailsWithInvalidCountry(t *testing.T) {
	origin, dest := "us", ""
	weight := 45.0
	path := fmt.Sprintf("/pricing?origin=%s&destination=%s&weight=%.2f", origin, dest, weight)

	w := handleRequest(t, http.MethodGet, path, nil)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func handleRequest(t *testing.T, method, path string, payload []byte) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	var (
		req *http.Request
		err error
	)
	if len(payload) > 0 {
		req, err = http.NewRequest(method, path, bytes.NewBuffer(payload))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	require.NoError(t, err)
	server.Router.ServeHTTP(w, req)

	return w
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:password@tcp(127.0.0.1:3306)/shipments?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/shipments?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func cleanup() {
	db.Exec("DELETE FROM shipments")
	db.Exec("DELETE FROM customers")
}
