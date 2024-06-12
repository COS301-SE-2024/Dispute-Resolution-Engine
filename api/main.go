package main

import (
	// "api/old_api"
	// "api/storage"
	"api/db"
	"api/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	DB := db.Init()
	h := handlers.New(DB)
	router := mux.NewRouter()

	router.HandleFunc("/createAcc", h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/login", h.LoginUser).Methods(http.MethodPost)
	log.Println("API server is running on port 8080")
	http.ListenAndServe(":8080", router)
	// store, err := storage.NewPostgresStore()
	// if err != nil {
	// 	log.Fatalf("could not create storage: %v", err)
	// }

	// if err := store.Init(); err != nil {
	// 	log.Fatalf("could not init storage: %v", err)
	// }

	// // store.Ping()

	// // users, err := store.GetAllUsers()

	// // if err != nil {
	// // 	log.Fatalf("could not get users: %v", err)
	// // }

	// // for _, user := range users {
	// // 	log.Printf("user: %v", user)
	// // }

	// server := api.NewServer(":8080", store)
	// server.Run()

}
