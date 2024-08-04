package main

import (
	"fmt"
	"orchestrator/db"
	"orchestrator/env"
	"orchestrator/utilities"
	"orchestrator/workflow"
)

var requiredEnvVariables = []string{
	// PostGres-related variables
	"DATABASE_URL",
	"DATABASE_PORT",
	"DATABASE_USER",
	"DATABASE_PASSWORD",
	"DATABASE_NAME",
}

func main() {
	fmt.Println("Hello, World!")
	logger := utilities.NewLogger().LogWithCaller()
	env.LoadFromFile(".env", "api.env")

	for _, key := range requiredEnvVariables {
		env.Register(key)
	}

	DB, err := db.Init()
	if err != nil {
		logger.WithError(err).Fatal("Couldn't initialize database connection")
	}
	fmt.Println(DB)
	fmt.Println(workflow.StateAppealCommunique)
}
