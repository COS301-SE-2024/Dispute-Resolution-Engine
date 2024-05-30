package main

import (
	"api/api"
)

func main() {
	server := api.NewServer(":8080")
	server.Run()
}
