package testhelpers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os/exec"
	"path"
)

// Package packages the current application into a container
func Package(ctx context.Context, host string) (testcontainers.Container, error) {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, err
	}
	basePath := string(cmdOut)

	fmt.Printf("\n\nStart container up testing!!! %s\n\n", basePath)
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				//Context:    "../",
				Context:    basePath + "/",
				Dockerfile: path.Join(basePath, "Dockerfile"),
				Repo:       "lordrahl",
				Tag:        "testing",
				//PrintBuildLog: true,
			},
			//Cmd:          []string{"./shipments"},
			ExposedPorts: []string{"8080/tcp"},
			Env: map[string]string{
				"ENVIRONMENT": "blackbox",
				"DB_HOST":     host,
				"DB_PORT":     "3306",
				"DB_NAME":     "shipments",
				"DB_USER":     "root",
				"DB_PASSWORD": "rootpassword",
			},
			WaitingFor: wait.ForHTTP("/").WithPort("8080/tcp"),
		},

		Started: true,
	})
}
