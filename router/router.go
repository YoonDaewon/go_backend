package router

import (
	"os"

	"go_backend/controller"
	"go_backend/database"
	"go_backend/repository"
	"go_backend/usecase"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes and returns the gin engine
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize dependencies
	// Choose repository based on environment variable or default to in-memory
	var userRepo repository.UserRepository
	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "postgres":
		if database.PostgresDB != nil {
			userRepo = repository.NewPostgresUserRepository()
		} else {
			userRepo = repository.NewUserRepository() // Fallback to in-memory
		}
	case "mongodb":
		if database.MongoDB != nil {
			userRepo = repository.NewMongoUserRepository()
		} else {
			userRepo = repository.NewUserRepository() // Fallback to in-memory
		}
	default:
		// Default to in-memory or PostgreSQL if available
		if database.PostgresDB != nil {
			userRepo = repository.NewPostgresUserRepository()
		} else {
			userRepo = repository.NewUserRepository()
		}
	}

	userUsecase := usecase.NewUserUsecase(userRepo)
	userController := controller.NewUserController(userUsecase)

	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// API routes
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", userController.CreateUser)
			users.GET("", userController.GetAllUsers)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
	}

	return r
}
