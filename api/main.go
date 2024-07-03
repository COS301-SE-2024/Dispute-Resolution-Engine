package main

import (
	"api/db"
	_ "api/docs" // This is important to import your generated docs package
	"api/handlers"
	"api/middleware"
	"api/redisDB"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Tries to load env files. If an error occurs, it will ignore the file and log the error
func loadEnvFile(files ...string) {
	for _, path := range files {
		if err := godotenv.Load(path); err != nil {
			log.Printf("Error loading env file: %v\n", err.Error())
		} else {
			log.Printf("Loaded env file: %v\n", path)
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
	loadEnvFile(".env", "api.env")

	DB := db.Init()

	redisClient := redisDB.InitRedis()

	authHandler := handlers.NewAuthHandler(DB, redisClient)
	userHandler := handlers.NewUserHandler(DB, redisClient)
	disputeHandler := handlers.NewDisputeHandler(DB, redisClient)
	utilityHandler := handlers.NewUtilitiesHandler(DB, redisClient)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	//setup handlers
	utilGroup := router.Group("/utils")
	utilGroup.GET("/countries", utilityHandler.GetCountries)

	authGroup := router.Group("/auth")
	handlers.SetupAuthRoutes(authGroup, authHandler)

	userGroup := router.Group("/user")
	userGroup.Use(middleware.JWTMiddleware)
	handlers.SetupUserRoutes(userGroup, userHandler)

	disputeGroup := router.Group("/disputes")
	handlers.SetupDisputeRoutes(disputeGroup, disputeHandler)

	archiveGroup := router.Group("/archive")
	handlers.SetupArchiveRoutes(archiveGroup, disputeHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("API server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
