package testhelpers

import (
	"context"
	"fmt"
	"log"

	mysqlModule "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// GetMySQLContainer returns a mysql testcontainer
func GetMySQLContainer(ctx context.Context) *mysqlModule.MySQLContainer {
	fmt.Printf("\n\nStarting Test Container\n\n")
	container, err := mysqlModule.Run(ctx, "mysql:8.0",
		mysqlModule.WithDatabase("shipments"),
		mysqlModule.WithUsername("root"),
		mysqlModule.WithPassword("rootpassword"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return container
}
