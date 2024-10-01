package testhelpers

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"

	mysqlModule "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// SetupTestDB setup test db across the entire systems
func SetupTestDB(ctx context.Context, container *mysqlModule.MySQLContainer) (*gorm.DB, error) {
	if container == nil {
		fmt.Printf("\n\nContainer is empty. Not Starting Test Container\n\n")
	}
	var (
		dsn string
		err error
	)

	fmt.Printf("\n\nEnvironment: %s\n\n", os.Getenv("ENVIRONMENT"))

	if os.Getenv("ENVIRONMENT") == "dagger" {
		dsn = "mysql://shipment_user:rootpassword@shipmentsdb:3306/shipments"
	} else {
		dsn, err = container.ConnectionString(ctx)
	}
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n\nDSN: %s\n\n", dsn)
	dsn = dsn + "?charset=utf8mb4&parseTime=True&loc=Local"

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
