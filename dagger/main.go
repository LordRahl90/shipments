package main

import (
	"context"

	"dagger/shipments/internal/dagger"
)

type Shipments struct{}

// Test initializes the test
func (m *Shipments) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	//mysql := dag.Container().
	//	From("mysql:8.0").
	//	WithExposedPort(3306)
	return m.BuildEnv(source).
		WithExec([]string{"make", "test"}).
		Stdout(ctx)
}

// BuildEnv initializes the build environment
func (m *Shipments) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"ls", "-al"})
}
