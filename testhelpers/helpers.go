package testhelpers

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	mysqlModule "github.com/testcontainers/testcontainers-go/modules/mysql"
)

// SetupTestDB setup test db across the entire systems
func SetupTestDB(ctx context.Context, container *mysqlModule.MySQLContainer) (*gorm.DB, error) {
	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}
	dsn = dsn + "?charset=utf8mb4&parseTime=True&loc=Local"

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
