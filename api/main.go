package main

import (
	"api/api"
	"api/storage"
	"log"
)

func main() {
	store, err := storage.NewPostgresStore()
	if err != nil {
		log.Fatalf("could not create storage: %v", err)
	}

	if err := store.Init(); err != nil {
		log.Fatalf("could not init storage: %v", err)
	}

	server := api.NewServer(":8080", store)
	server.Run()
}
