package main

import (
	"api/db"
	_ "api/docs" // This is important to import your generated docs package
	"api/handlers"
	"api/middleware"
	"api/redisDB"
	"api/utilities"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Tries to load env files. If an error occurs, it will ignore the file and log the error
func loadEnvFile(files ...string) {
	logger := utilities.NewLogger().LogWithCaller()
	for _, path := range files {
		if err := godotenv.Load(path); err != nil {
			logger.WithError(err).Warning("Env file not found")
		} else {
			logger.Info("Loaded env file")
		}
	}
}

// @title Dispute Resolution Engine - v1
// @version 1.0
// @description This is a description.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	logger := utilities.NewLogger().LogWithCaller()
	loadEnvFile(".env", "api.env")

	DB, err := db.Init()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	logger.Info("Connected to database successfully")

	_, err = redisDB.InitRedis()
	if err != nil {
		logger.WithError(err).Fatal("Error initializing Redis")
	}
	logger.Info("Connected to Redis successfully")

	authHandler := handlers.NewAuthHandler(DB)
	userHandler := handlers.NewUserHandler(DB)
	disputeHandler := handlers.NewDisputeHandler(DB)
	archiveHandler := handlers.NewArchiveHandler(DB)
	expertHandler := handlers.NewExpertHandler(DB)
	utilityHandler := handlers.NewUtilitiesHandler(DB)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	router.Static("/filestorage", os.Getenv("FILESTORAGE_ROOT"))

	//setup handlers
	utilGroup := router.Group("/utils")
	utilGroup.GET("/countries", utilityHandler.GetCountries)
	utilGroup.GET("/dispute_statuses", utilityHandler.GetDisputeStatuses)

	authGroup := router.Group("/auth")
	handlers.SetupAuthRoutes(authGroup, authHandler)

	userGroup := router.Group("/user")
	userGroup.Use(middleware.JWTMiddleware)
	handlers.SetupUserRoutes(userGroup, userHandler)

	disputeGroup := router.Group("/disputes")
	handlers.SetupDisputeRoutes(disputeGroup, disputeHandler)

	archiveGroup := router.Group("/archive")
	handlers.SetupArchiveRoutes(archiveGroup, archiveHandler)

	expertGroup := router.Group("/experts")
	handlers.SetupExpertRoutes(expertGroup, expertHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	http.ListenAndServe(":8080", router)
	logger.Info("API started successfully on port 8080")
}
