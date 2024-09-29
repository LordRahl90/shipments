package main

import (
	"context"

	"dagger/shipments/internal/dagger"
)

type Shipments struct {
}

// Test initializes the test
func (m *Shipments) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	if err := dag.Testcontainers().StartDocker(ctx); err != nil {
		return "", err
	}
	mysql := dag.Container().
		From("mysql:8.0").
		WithEnvVariable("MYSQL_DATABASE", "shipments").
		WithEnvVariable("MYSQL_USER", "shipment_user").
		WithEnvVariable("MYSQL_PASSWORD", "rootpassword").
		WithEnvVariable("MYSQL_ROOT_PASSWORD", "rootpassword").
		WithExposedPort(3306).
		AsService()

	return m.BuildEnv(ctx, source).
		//With(dag.Docker()).
		WithEnvVariable("ENVIRONMENT", "dagger").
		WithServiceBinding("shipmentsdb", mysql).
		//WithExec([]string{"ls", "-al"}).
		//WithExec([]string{"docker"}).
		//WithExec([]string{"docker", "version"}).
		WithExec([]string{"make", "test"}).
		Stdout(ctx)
}

// Lint initializes the lint
func (m *Shipments) Lint(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(ctx, source).
		WithExec([]string{"make", "lint"}).
		Stdout(ctx)

}

// BuildEnv initializes the build environment
func (m *Shipments) BuildEnv(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src")
	//.
	//With(dag.
	//	Testcontainers().
	//	Setup)
	//WithEnvVariable("DOCKER_HOST", "tcp://docker:2375").
	//WithServiceBinding("docker", dag.Docker().Daemon().Service())

}
