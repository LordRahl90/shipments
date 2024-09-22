package shipments

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"shipments/requests"
	"shipments/responses"
	"shipments/servers"
	"shipments/testhelpers"
	requests2 "shipments/testhelpers/requests"

	"github.com/cucumber/godog"
)

type responseString string

var (
	//req requests.Shipment
	name, email string

	container = testhelpers.GetMySQLContainer(context.TODO())
)

//var svr *servers.Server

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^I am "([^"]*)" With email "([^"]*)"$`, arrange)
	ctx.When(`I create a shipment with country "([^"]*)" and destination country of "([^"]*)" and weight of (\d+\.\d+)`, action)
	ctx.Then(`^I should see the shipment with price of (\d+\.\d+)$`, assert)
	ctx.Then(`I should see the shipment with a non empty reference`, nonEmptyReference)
}

func arrange(ctx context.Context, n, e string) (context.Context, error) {
	fmt.Printf("Setting user with name and email to %s and %s\n", n, e)
	name = n
	email = e

	return ctx, nil
}

func action(ctx context.Context, country, destination string, weight float64) (context.Context, error) {
	fmt.Printf("\n\nOrigin: %s, Destination: %s, Weight: %f\n\n", country, destination, weight)
	req := requests.Shipment{
		Name:        name,
		Email:       email,
		Origin:      country,
		Destination: destination,
		Weight:      weight,
	}
	db, err := testhelpers.SetupTestDB(ctx, container)
	if err != nil {
		return ctx, err
	}

	svr, err := servers.New(db)
	if err != nil {
		return ctx, err
	}

	path := "/new"
	b, err := json.Marshal(req)
	if err != nil {
		return ctx, err
	}

	fmt.Printf("\n\nRequesting %s\n\n", b)

	w, err := requests2.NewRequest(svr, http.MethodPost, path, b)
	if err != nil {
		return ctx, err
	}

	response := w.Body.String()

	return context.WithValue(ctx, "result", responseString(response)), nil
}

func assert(ctx context.Context, expected float64) (context.Context, error) {
	result := ctx.Value("result").(responseString)
	var res responses.Shipment
	if err := json.Unmarshal([]byte(result), &res); err != nil {
		return ctx, err
	}

	if res.Price != expected {
		return ctx, fmt.Errorf("expected %f, got %f", expected, res.Price)
	}
	return ctx, nil
}

func nonEmptyReference(ctx context.Context) (context.Context, error) {
	result, ok := ctx.Value("result").(responseString)
	if !ok {
		return ctx, fmt.Errorf("empty result")
	}
	var res responses.Shipment
	if err := json.Unmarshal([]byte(result), &res); err != nil {
		return ctx, err
	}
	if res.Reference == "" {
		return ctx, fmt.Errorf("empty reference")
	}
	return ctx, nil
}

func TestFeatures(t *testing.T) {
	err := container.Start(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	suite := &godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/shipment.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status received from suite")
	}

	err = container.Terminate(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
