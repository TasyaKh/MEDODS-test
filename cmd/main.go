package main

import (
	"log"
	"test-task/internal/api/routes"
	"test-task/internal/config"
	"test-task/internal/database"

	_ "test-task/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API
// @version         1.0
// @description     Тестовый API для авторизации
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в виде: Bearer {access_token}

func main() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg)

	if err != nil {
		log.Fatal("failed connect to db:", err)
	}

	// documentation path
	router.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))

	routes.SetupRoutes(router, cfg, db)

	port := cfg.Port
	log.Printf("server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
