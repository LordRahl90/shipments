package blackbox

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"shipments/testhelpers"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

var (
	dbContainer  = testhelpers.GetMySQLContainer(context.TODO())
	shipmentsApp testcontainers.Container

	appLink string
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		if err := dbContainer.Terminate(context.TODO()); err != nil {
			log.Fatal(err)
		}

		if err := shipmentsApp.Terminate(context.TODO()); err != nil {
			log.Fatal(err)
		}
		os.Exit(code)
	}()

	dbHost, err := dbContainer.ContainerIP(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	app, err := testhelpers.Package(context.TODO(), dbHost)
	if err != nil {
		log.Fatal(err)
	}
	shipmentsApp = app

	ip, err := shipmentsApp.ContainerIP(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	appLink = fmt.Sprintf("http://%s:8080", ip)
	fmt.Printf("\n\nApp Link: %s\n\n", appLink)

	code = m.Run()
}

func TestLanding(t *testing.T) {
	ctx := context.TODO()

	logs, err := shipmentsApp.Logs(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, logs)

	b, err := io.ReadAll(logs)
	require.NoError(t, err)

	fmt.Printf("Container logs:\n\n %s\n\n", b)

	mp, err := shipmentsApp.MappedPort(context.TODO(), "8080/tcp")

	fmt.Printf("\n\nMapped Ports: %v\n\n", mp.Port())

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	path := fmt.Sprintf("http://localhost:%s/", mp.Port())
	fmt.Printf("\n\nPath: %s\n\n", path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.StatusCode)

	require.NoError(t, logs.Close())

	//w := handleRequest(t, http.MethodGet, "/", nil)
	//require.Equal(t, http.StatusOK, w.Code)
}
