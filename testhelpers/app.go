package testhelpers

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Package packages the current application into a container
func Package(ctx context.Context) (testcontainers.Container, error) {
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../",
				Dockerfile: "Dockerfile",
			},
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForLog(`{"msg":"server started at port 8080"}`),
		},
		Started: true,
	})
}
