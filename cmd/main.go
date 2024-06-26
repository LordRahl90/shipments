package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"shipments/domains/tracing"
	"shipments/servers"

	"go.opentelemetry.io/otel"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	otelEndpoint string
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("starting server")
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" {
		if err := godotenv.Load(".envs/.env"); err != nil {
			panic(err)
		}
	}

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	otelEndpoint = os.Getenv("OTLP_ENDPOINT")
	if otelEndpoint == "" {
		log.Fatal("OTLP Endpoint not set")
	}

	//exp, err := tracing.ConsoleExporter()
	exp, err := tracing.TempoExporter(context.Background(), otelEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	tp := tracing.TraceProvider(exp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)

	server, err := servers.New(db)
	if err != nil {
		panic(err)
	}

	slog.Info("server started at port 8080")
	log.Fatal(server.Router.Run("0.0.0.0:8080"))
}

func setupDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
