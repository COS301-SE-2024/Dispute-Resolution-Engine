package main

import (
    "api/db"
    "api/handlers"
    "log"
    "net/http"

    _ "api/docs" // This is important to import your generated docs package

    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
    "github.com/gorilla/mux"
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
    h := handlers.New(DB)
    router := mux.NewRouter()

    router.HandleFunc("/createAcc", h.CreateUser).Methods(http.MethodPost)
    router.HandleFunc("/login", h.LoginUser).Methods(http.MethodPost)
    router.HandleFunc("/utils/countries", h.GetCountries).Methods(http.MethodGet)

    // Swagger setup
    setupSwaggerDocs(router)

    log.Println("API server is running on port 8080")
    http.ListenAndServe(":8080", router)
}

func setupSwaggerDocs(router *mux.Router) {
    // Create a new gin engine
    ginRouter := gin.Default()
    ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Serve the gin engine on a specific route in the main mux router
    router.PathPrefix("/swagger/").Handler(ginRouter)
}
