package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"shipments/domains/tracing"
	"shipments/servers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	otelLogger = otelslog.NewHandler("shipments")
	jsonLogger = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	gin.SetMode(gin.ReleaseMode)

	logger := slogmulti.Fanout(otelLogger, jsonLogger)
	slog.SetDefault(slog.New(logger))

	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" {
		if err := godotenv.Load(".envs/.env"); err != nil {
			log.Fatal(err)
		}
	}
	slog.Info("starting server in " + env + " mode")
	tp, err := tracing.TraceProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}
	lp, err := tracing.LoggerProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}

	//if err := db.Use(otelgorm.NewPlugin()); err != nil {
	//	log.Fatal(err)
	//}

	errChan := make(chan error, 1)

	server, err := servers.New(db)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		slog.Info("starting server on port 8080")
		if err := server.Router.Run("0.0.0.0:8080"); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		log.Fatal(err)
	case <-ctx.Done():
		log.Println("shutting down server")
	}

	slog.InfoContext(context.TODO(), "starting cleanup")
	if err := tp.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}

	if err := lp.Shutdown(context.TODO()); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
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
