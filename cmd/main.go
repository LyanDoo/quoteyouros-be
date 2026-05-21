package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/quoteyouros/backend/internal/config"
	"github.com/quoteyouros/backend/internal/handler"
	"github.com/quoteyouros/backend/internal/infrastructure/postgres"
	"github.com/quoteyouros/backend/internal/middleware"
	blogrepo "github.com/quoteyouros/backend/internal/repository/blog"
	profilerepo "github.com/quoteyouros/backend/internal/repository/profile"
	projectrepo "github.com/quoteyouros/backend/internal/repository/project"
	authrepo "github.com/quoteyouros/backend/internal/repository/user"
	authuc "github.com/quoteyouros/backend/internal/usecase/auth"
	bloguc "github.com/quoteyouros/backend/internal/usecase/blog"
	profileuc "github.com/quoteyouros/backend/internal/usecase/profile"
	projectuc "github.com/quoteyouros/backend/internal/usecase/project"
	"github.com/quoteyouros/backend/pkg/fileupload"
	applogger "github.com/quoteyouros/backend/pkg/logger"
)

func main() {
	// Load configuration
	applogger.Info("main: loading configuration")
	cfg := config.LoadConfig()
	applogger.Info("main: configuration loaded successfully", "port", cfg.Server.Port)

	// Connect to database
	applogger.Info("main: connecting to database", "host", cfg.Database.Host, "port", cfg.Database.Port)
	db := postgres.NewConnection(&cfg.Database)
	defer postgres.Close(db)
	applogger.Info("main: database connection established")

	// Initialize Fiber app
	applogger.Info("main: initializing Fiber app")
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
	applogger.Info("main: middleware configured")

	// Initialize repositories
	applogger.Debug("main: initializing repositories")
	userRepository := authrepo.NewUserRepository(db)
	blogRepository := blogrepo.NewBlogRepository(db)
	projectRepository := projectrepo.NewProjectRepository(db)
	profileRepository := profilerepo.NewProfileRepository(db)

	// Initialize file upload service
	applogger.Debug("main: initializing file upload service")
	fileUploadService := fileupload.New(fileupload.ResumeStoragePath)

	// Initialize use cases
	applogger.Debug("main: initializing use cases")
	authUseCase := authuc.New(userRepository, cfg.JWT.Secret, cfg.JWT.Expiration)
	blogUseCase := bloguc.New(blogRepository)
	projectUseCase := projectuc.New(projectRepository)
	profileUseCase := profileuc.New(profileRepository, fileUploadService)

	// Initialize handlers
	applogger.Debug("main: initializing handlers")
	authHandler := handler.NewAuthHandler(authUseCase)
	blogHandler := handler.NewBlogHandler(blogUseCase)
	projectHandler := handler.NewProjectHandler(projectUseCase)
	profileHandler := handler.NewProfileHandler(profileUseCase)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		applogger.Debug("health: check endpoint called")
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "QuoteYourOS Backend is running",
		})
	})

	// API v1 Routes
	applogger.Debug("main: setting up routes")
	api := app.Group("/api")

	// Public routes
	setupPublicRoutes(api, authHandler, blogHandler, projectHandler, profileHandler)

	// Protected routes (admin)
	jwtMiddleware := middleware.JWTAuth(cfg.JWT.Secret)
	setupProtectedRoutes(api, authHandler, blogHandler, projectHandler, profileHandler, jwtMiddleware)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	applogger.Info("main: starting server", "address", addr, "environment", cfg.Server.Env)

	if err := app.Listen(addr); err != nil {
		applogger.Error("main: server failed to start", "error", err.Error())
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setupPublicRoutes(app fiber.Router, authHandler *handler.AuthHandler, blogHandler *handler.BlogHandler, projectHandler *handler.ProjectHandler, profileHandler *handler.ProfileHandler) {
	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Blog routes
	blog := app.Group("/blog")
	blog.Get("", blogHandler.GetAllBlogPosts)
	blog.Get("/:id", blogHandler.GetBlogPost)

	// Projects routes
	projects := app.Group("/projects")
	projects.Get("", projectHandler.GetAllProjects)
	projects.Get("/:id", projectHandler.GetProject)

	// Profile routes
	profile := app.Group("/profile")
	profile.Get("/about", profileHandler.GetProfile)
	profile.Get("/resume/download", profileHandler.DownloadResume)

	// Contact routes
	app.Post("/contact", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "POST /api/contact - Not yet implemented"})
	})
}

func setupProtectedRoutes(app fiber.Router, authHandler *handler.AuthHandler, blogHandler *handler.BlogHandler, projectHandler *handler.ProjectHandler, profileHandler *handler.ProfileHandler, jwtMiddleware fiber.Handler) {
	// Auth protected routes
	auth := app.Group("/auth", jwtMiddleware)
	auth.Get("/me", authHandler.GetCurrentUser)
	auth.Post("/logout", authHandler.Logout)

	// Admin blog routes (protected)
	blog := app.Group("/blog", jwtMiddleware)
	blog.Post("", blogHandler.CreateBlogPost)
	blog.Put("/:id", blogHandler.UpdateBlogPost)
	blog.Delete("/:id", blogHandler.DeleteBlogPost)

	// Admin projects routes (protected)
	projects := app.Group("/projects", jwtMiddleware)
	projects.Post("", projectHandler.CreateProject)
	projects.Put("/:id", projectHandler.UpdateProject)
	projects.Delete("/:id", projectHandler.DeleteProject)

	// Admin profile routes (protected)
	profile := app.Group("/profile", jwtMiddleware)
	profile.Put("/about", profileHandler.UpdateProfile)
	profile.Post("/resume", profileHandler.UploadResume)

	// Admin messages routes (protected)
	messages := app.Group("/messages", jwtMiddleware)
	messages.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "GET /api/messages - Not yet implemented"})
	})
	messages.Delete("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "DELETE /api/messages/:id - Not yet implemented"})
	})
}
