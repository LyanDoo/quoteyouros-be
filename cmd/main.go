package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/quoteyouros/backend/internal/config"
	"github.com/quoteyouros/backend/internal/infrastructure/postgres"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db := postgres.NewConnection(&cfg.Database)
	defer postgres.Close(db)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "QuoteYourOS Backend API v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.CORS.Origins, ","),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders: "Content-Type,Authorization,X-Requested-With",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "QuoteYourOS Backend is running",
		})
	})

	// API v1 Routes
	api := app.Group("/api")

	// Public routes
	setupPublicRoutes(api)

	// Protected routes (admin)
	setupProtectedRoutes(api)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("✓ Starting server on %s (environment: %s)\n", addr, cfg.Server.Env)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupPublicRoutes(app fiber.Router) {
	// Blog routes
	blog := app.Group("/blog")
	blog.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/blog - Not yet implemented"})
	})
	blog.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/blog/:id - Not yet implemented"})
	})

	// Projects routes
	projects := app.Group("/projects")
	projects.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/projects - Not yet implemented"})
	})
	projects.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/projects/:id - Not yet implemented"})
	})

	// Contact routes
	app.Post("/contact", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "POST /api/contact - Not yet implemented"})
	})

	// Profile routes
	profile := app.Group("/profile")
	profile.Get("/about", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/profile/about - Not yet implemented"})
	})
	profile.Get("/resume", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/profile/resume - Not yet implemented"})
	})
	profile.Get("/resume/download", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/profile/resume/download - Not yet implemented"})
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "POST /api/auth/login - Not yet implemented"})
	})
	auth.Post("/logout", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "POST /api/auth/logout - Not yet implemented"})
	})
	auth.Get("/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/auth/me - Not yet implemented"})
	})
}

func setupProtectedRoutes(app fiber.Router) {
	// Admin blog routes (protected)
	blog := app.Group("/blog")
	blog.Post("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "POST /api/blog - Not yet implemented"})
	})
	blog.Put("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "PUT /api/blog/:id - Not yet implemented"})
	})
	blog.Delete("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "DELETE /api/blog/:id - Not yet implemented"})
	})

	// Admin projects routes (protected)
	projects := app.Group("/projects")
	projects.Post("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "POST /api/projects - Not yet implemented"})
	})
	projects.Put("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "PUT /api/projects/:id - Not yet implemented"})
	})
	projects.Delete("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "DELETE /api/projects/:id - Not yet implemented"})
	})

	// Admin messages routes (protected)
	messages := app.Group("/messages")
	messages.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/messages - Not yet implemented"})
	})
	messages.Delete("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "DELETE /api/messages/:id - Not yet implemented"})
	})
}
