package main

import (
	"fmt"
	"log"

	"github.com/harri88/eticketing-challenge/config"
	delivery "github.com/harri88/eticketing-challenge/internal/delivery/http"
	"github.com/harri88/eticketing-challenge/internal/infrastructure/client"
	"github.com/harri88/eticketing-challenge/internal/infrastructure/gateway"
	"github.com/harri88/eticketing-challenge/internal/repository"
	"github.com/harri88/eticketing-challenge/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/harri88/eticketing-challenge/docs"
	_ "github.com/lib/pq"
)

// @title Payment Service API
// @version 1.0
// @description Payment processing service for e-ticketing platform
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @host localhost:8081
// @BasePath /
func main() {
	// 1. Load Configuration
	cfg := config.Load()

	log.Printf("Starting Payment Service in %s environment", cfg.Environment)
	log.Printf("Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Printf("Server: %s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Ticket Service: %s (timeout: %ds)", cfg.TicketService.URL, cfg.TicketService.Timeout)
	log.Printf("Payment Gateway: %s", cfg.PaymentGateway.DefaultProvider)

	// 2. Init DB
	db, err := repository.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✓ Database connection established")

	// 3. Init Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 4. Add CORS Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:3001", "http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Idempotency-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// 4. Init Layers (Dependency Injection)
	repo := repository.NewPostgresRepository(db)
	gwFactory := gateway.NewPaymentGatewayFactory()
	ticketClient := client.NewTicketClient(cfg.TicketService.URL)

	payUsecase := usecase.NewPaymentUsecase(repo, gwFactory, ticketClient)
	txnUsecase := usecase.NewTransactionUsecase(repo)

	// 5. Init Handlers
	delivery.NewPaymentHandler(e, payUsecase)
	delivery.NewTransactionHandler(e, txnUsecase)

	// 6. Start Server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	log.Printf("✓ Payment Service started on %s", serverAddr)
	log.Printf("Swagger UI available at http://%s:%d/swagger/index.html", cfg.Server.Host, cfg.Server.Port)

	if err := e.Start(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
