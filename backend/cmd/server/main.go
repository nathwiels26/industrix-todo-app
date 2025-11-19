package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/industrix-todo-app/backend/internal/config"
	"github.com/industrix-todo-app/backend/internal/database"
	"github.com/industrix-todo-app/backend/internal/handlers"
	"github.com/industrix-todo-app/backend/internal/models"
	"github.com/industrix-todo-app/backend/internal/repository"
	"github.com/industrix-todo-app/backend/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.Category{}, &models.Todo{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Initialize services
	categoryService := services.NewCategoryService(categoryRepo)
	todoService := services.NewTodoService(todoRepo)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	todoHandler := handlers.NewTodoHandler(todoService)

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.POST("", categoryHandler.Create)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		// Todo routes
		todos := api.Group("/todos")
		{
			todos.GET("", todoHandler.GetAll)
			todos.POST("", todoHandler.Create)
			todos.GET("/:id", todoHandler.GetByID)
			todos.PUT("/:id", todoHandler.Update)
			todos.DELETE("/:id", todoHandler.Delete)
			todos.PATCH("/:id/complete", todoHandler.ToggleComplete)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
