package pricing

import (
	"context"
	"fmt"
	"testing"

	"shipments/domains/core"
	"shipments/requests"

	"github.com/cucumber/godog"
)

type ctxResult float64

var (
	req requests.Pricing
)

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^I have country "([^"]*)" and destination country of "([^"]*)" and weight of (\d+\.\d+)$`, arrange)
	ctx.When(`I get price`, action)
	ctx.Then(`^I should see the price of (\d+\.\d+)$`, assert)
}

func arrange(ctx context.Context, country, destination string, weight float64) (context.Context, error) {
	req.Origin = country
	req.Destination = destination
	req.Weight = weight
	return ctx, nil
}

func action(ctx context.Context) (context.Context, error) {
	result, err := core.PriceFromSize(ctx, req.Weight, req.Origin, req.Destination)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, "result", ctxResult(result)), nil
}

func assert(ctx context.Context, expected float64) (context.Context, error) {
	result := ctx.Value("result").(ctxResult)
	if float64(result) != expected {
		return ctx, fmt.Errorf("expected %f, got %f", expected, result)
	}
	return ctx, nil
}

func TestFeatures(t *testing.T) {
	suite := &godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/pricing.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status received from suite")
	}
}
