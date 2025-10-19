package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/VictorRibeiroH/ativvo-backend/internal/config"
	"github.com/VictorRibeiroH/ativvo-backend/internal/database"
	"github.com/VictorRibeiroH/ativvo-backend/internal/handlers"
	"github.com/VictorRibeiroH/ativvo-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	log.Println("✅ Configuration loaded")

	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "Ativvo API",
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	allowedOrigins := strings.Split(config.AppConfig.FrontendURL, ",")
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	log.Printf("🌐 CORS enabled for origins: %s", strings.Join(allowedOrigins, ", "))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	api := app.Group("/api")
	
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"env":    config.AppConfig.Environment,
		})
	})

	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	protected := api.Group("", middleware.AuthRequired)
	
	protected.Get("/me", handlers.Me)
	protected.Put("/profile", handlers.UpdateProfile)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("\n🛑 Shutting down gracefully...")
		
		if err := database.Close(); err != nil {
			log.Printf("❌ Error closing database: %v", err)
		}
		
		if err := app.Shutdown(); err != nil {
			log.Printf("❌ Error shutting down server: %v", err)
		}
		
		log.Println("👋 Server stopped")
		os.Exit(0)
	}()

	port := config.AppConfig.Port
	log.Printf("🚀 Server starting on port %s (env: %s)", port, config.AppConfig.Environment)
	
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
