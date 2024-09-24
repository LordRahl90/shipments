package main

import (
	"context"

	"dagger/shipments/internal/dagger"
)

type Shipments struct{}

// Test initializes the test
func (m *Shipments) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	if _, err := dag.Testcontainers().DockerService().Start(ctx); err != nil {
		return "", err
	}
	//mysql := dag.Container().
	//	From("mysql:8.0").
	//	WithExposedPort(3306)
	return m.BuildEnv(source).
		WithExec([]string{"make", "test"}).
		Stdout(ctx)
}

func (m *Shipments) Lint(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{"make", "lint"}).
		Stdout(ctx)

}

// BuildEnv initializes the build environment
func (m *Shipments) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"ls", "-al"}).
		With(dag.Testcontainers().Setup)
}
