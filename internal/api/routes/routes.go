package routes

import (
	"test-task/internal/api/handlers"
	"test-task/internal/api/middleware"
	"test-task/internal/config"
	"test-task/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, db *gorm.DB) {
	sessionService := services.NewSessionService(db, cfg)
	userService := services.NewUserService(db)
	authService := services.NewAuthService(cfg, db, sessionService, userService)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.GET("/login/:guid", authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}

	moddleawar := middleware.AuthMiddleware(authService, userService, sessionService)

	api.GET("/me", moddleawar, userHandler.GetMe)

}
