package main

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"testing"
)

type hotdogCounter struct{}

func iEatHotdog(ctx context.Context, toEat int) (context.Context, error) {
	fmt.Printf("\n\nEating %d hotdogs\n\n", toEat)
	available, ok := ctx.Value(hotdogCounter{}).(int)
	if !ok {
		return nil, fmt.Errorf("hotdog counter not found")
	}

	if toEat > available {
		return nil, fmt.Errorf("not enough hotdogs to eat")
	}

	left := available - toEat
	return context.WithValue(ctx, hotdogCounter{}, left), godog.ErrPending
}

func thereAreHotdogs(ctx context.Context, count int) (context.Context, error) {
	fmt.Printf("\n\nInitializing\n\n")
	return context.WithValue(ctx, hotdogCounter{}, count), nil
}
func thereShouldBeLeft(ctx context.Context, left int) (context.Context, error) {
	fmt.Printf("\n\nEvaluated with %d hotdogs left\n\n", left)
	available, ok := ctx.Value(hotdogCounter{}).(int)
	if !ok {
		return nil, fmt.Errorf("hotdog counter not found")
	}

	if available != left {
		return nil, fmt.Errorf("expected %d, got %d", left, available)
	}
	return ctx, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have (\d+) hotdogs$`, thereAreHotdogs)
	ctx.Step(`^I eat (\d+) hotdogs$`, iEatHotdog)
	ctx.Step(`^I should have (\d+) hotdogs left$`, thereShouldBeLeft)
}

func TestFeatures(t *testing.T) {
	suite := &godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/godogs.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status received from suite")
	}
}
