package main

import (
	"api/db"
	_ "api/docs" // This is important to import your generated docs package
	"api/env"
	"api/handlers"
	"api/handlers/dispute"
	"api/middleware"
	"api/redisDB"
	"api/utilities"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var requiredEnvVariables = []string{

	// PostGres-related variables
	"DATABASE_URL",
	"DATABASE_PORT",
	"DATABASE_USER",
	"DATABASE_PASSWORD",
	"DATABASE_NAME",

	// Redis-related variables
	"REDIS_URL",
	"REDIS_PASSWORD",
	"REDIS_DB",

	// Variables for file storage
	"FILESTORAGE_ROOT",
	"FILESTORAGE_URL",

	// Variables for sending email using SMTP
	"COMPANY_EMAIL",
	"COMPANY_AUTH",

	// Miscellaneous
	"FRONTEND_BASE_URL",
	"JWT_SECRET",
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
	jwt := middleware.NewJwtMiddleware()
	logger := utilities.NewLogger().LogWithCaller()
	envLoader := env.NewEnvLoader()
	envLoader.LoadFromFile(".env", "api.env")

	for _, key := range requiredEnvVariables {
		envLoader.Register(key)
	}

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
	disputeHandler := dispute.NewHandler(DB)
	archiveHandler := handlers.NewArchiveHandler(DB)
	expertHandler := handlers.NewExpertHandler(DB)
	utilityHandler := handlers.NewUtilitiesHandler(DB)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	fileStorageRoot, err := envLoader.Get("FILESTORAGE_ROOT")
	if err != nil {
		return
	}
	router.Static("/filestorage", fileStorageRoot)

	//setup handlers
	utilGroup := router.Group("/utils")
	utilGroup.GET("/countries", utilityHandler.GetCountries)
	utilGroup.GET("/dispute_statuses", utilityHandler.GetDisputeStatuses)

	authGroup := router.Group("/auth")
	handlers.SetupAuthRoutes(authGroup, authHandler)

	userGroup := router.Group("/user")
	userGroup.Use(jwt.JWTMiddleware)
	handlers.SetupUserRoutes(userGroup, userHandler)

	disputeGroup := router.Group("/disputes")
	dispute.SetupRoutes(disputeGroup, disputeHandler)

	archiveGroup := router.Group("/archive")
	handlers.SetupArchiveRoutes(archiveGroup, archiveHandler)

	expertGroup := router.Group("/experts")
	handlers.SetupExpertRoutes(expertGroup, expertHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	http.ListenAndServe(":8080", router)
	logger.Info("API started successfully on port 8080")
}
