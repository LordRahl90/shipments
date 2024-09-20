package testhelpers

import (
	"context"
	"log"

	mysqlModule "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// GetMySQLContainer returns a mysql testcontainer
func GetMySQLContainer(ctx context.Context) *mysqlModule.MySQLContainer {
	container, err := mysqlModule.Run(ctx, "mysql:8.0",
		mysqlModule.WithDatabase("metis"),
		mysqlModule.WithUsername("root"),
		mysqlModule.WithPassword("rootpassword"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return container
}
