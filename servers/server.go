package servers

import (
	"errors"
	"fmt"
	"net/http"

	"shipments/domains/core"
	"shipments/domains/customers"
	"shipments/domains/customers/store"
	"shipments/domains/entities"
	"shipments/domains/shipments"
	cStore "shipments/domains/shipments/store"
	"shipments/domains/tracing"
	"shipments/requests"
	"shipments/responses"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
)

var (
	customerService customers.ICustomerService
	shipmentService shipments.IShipmentService
)

// Server contains the essential configs for the server
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

// New returns a new server implementation
func New(db *gorm.DB) (*Server, error) {
	router := gin.New()
	router.Use(DefaultStructuredLogger(),
		otelgin.Middleware("shipments"))
	custStore, err := store.New(db)
	if err != nil {
		return nil, err
	}
	shipStore, err := cStore.New(db)
	if err != nil {
		return nil, err
	}
	customerService = customers.New(custStore)
	shipmentService = shipments.New(shipStore)

	router.GET("/", func(ctx *gin.Context) {
		_, span := tracing.Tracer().Start(ctx.Request.Context(), "Welcome")
		defer span.End()
		fmt.Printf("\n\nID: %s\n\n", span.SpanContext().TraceID().String())
		ctx.JSON(200, gin.H{
			"message": "Welcome to Shipments API service",
			"docs":    "https://www.getpostman.com/collections/497742b6deae56e91248",
		})
	})

	router.POST("/new", createShipment)
	router.GET("/history/:email", shipmentHistory)
	router.GET("/pricing", pricingDetails)
	return &Server{
		DB:     db,
		Router: router,
	}, nil
}

func pricingDetails(ctx *gin.Context) {
	traceCtx, span := tracing.Tracer().Start(ctx.Request.Context(), "Pricing Details")
	defer span.End()

	var p requests.Pricing
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	res, err := core.PriceFromSize(traceCtx, p.Weight, p.Origin, p.Destination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	result := &responses.Pricing{
		Weight:         p.Weight,
		Origin:         p.Origin,
		Destination:    p.Destination,
		WeightCategory: core.WeightFromSize(p.Weight).String(),
		Price:          res,
	}
	ctx.JSON(http.StatusOK, result)
}

func createShipment(ctx *gin.Context) {
	traceCtx, span := tracing.Tracer().Start(ctx.Request.Context(), "Create Shipment")
	defer span.End()
	var req requests.Shipment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	cust := &entities.Customer{
		Name:  req.Name,
		Email: req.Email,
	}
	err := customerService.FindOrCreate(traceCtx, cust)
	if err != nil {
		span.RecordError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	shipment := &entities.Shipment{
		CustomerID:  cust.ID,
		Origin:      req.Origin,
		Destination: req.Destination,
		Weight:      req.Weight,
	}

	err = shipmentService.Create(ctx, shipment)
	if err != nil {
		span.RecordError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	res := responses.Shipment{
		Success:    true,
		Reference:  shipment.ID,
		CustomerID: cust.ID,
		Price:      shipment.Price,
	}

	ctx.JSON(http.StatusCreated, res)
}

func shipmentHistory(ctx *gin.Context) {
	traceCtx, span := tracing.Tracer().Start(ctx.Request.Context(), "Shipment History")
	defer span.End()
	email := ctx.Param("email")
	if email == "" {
		span.RecordError(errors.New("invalid email provided"))
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "invalid email provided",
		})
		return
	}
	cust, err := customerService.FindByEmail(traceCtx, email)
	if err != nil {
		span.RecordError(err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	res, err := shipmentService.FindCustomerShipments(traceCtx, cust.ID)
	if err != nil {
		span.RecordError(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
