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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
	DB := db.Init()
	redisClient := redisDB.InitRedis()
    h := handlers.New(DB, redisClient)

	router := gin.Default()
    router.Use(cors.New(cors.Config{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Content-Type", "Authorization"},
    }))

	//setup handlers
    utilGroup := router.Group("/utils")
    utilGroup.GET("/countries", h.GetCountries)

	authGroup := router.Group("/auth")
	handlers.SetupAuthRoutes(authGroup, h)

	userGroup := router.Group("/user")
	userGroup.Use(middleware.JWTMiddleware)
	handlers.SetupUserRoutes(userGroup, h)

	disputeGroup := router.Group("/disputes")
	handlers.SetupDisputeRoutes(disputeGroup, h)

	archiveGroup := router.Group("/archive")
    handlers.SetupArchiveRoutes(archiveGroup, h)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("API server is running on port 8080")
	http.ListenAndServe(":8080", router)
}

