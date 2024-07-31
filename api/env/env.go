package env

import (
	"api/utilities"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var environment = map[string]string{}

var envInstance Env
var once sync.Once

type EnvLoader struct{
	environment map[string]string
	logger      utilities.Logger
}

func NewEnvLoader() Env {
	once.Do(func() {
		envInstance = createEnvLoader()
	})
	return envInstance
}

func createEnvLoader() Env {
	return &EnvLoader{
		environment: make(map[string]string),
		logger:      *utilities.NewLogger().LogWithCaller(),
	}
}

// Tries to load env files. If an error occurs, it will ignore the file and log the error
func (r *EnvLoader) LoadFromFile(files ...string) {
	logger := utilities.NewLogger().LogWithCaller()
	for _, path := range files {
		if err := godotenv.Load(path); err != nil {
			logger.WithError(err).Warningf("environment file %s not found. Skipping", path)
		} else {
			logger.Infof("Loaded env file %s", path)
		}
	}
}

// Registers the environment variable in the registry. If not found, it will log the errror and exit the program
func (r *EnvLoader) Register(key string) {
	logger := utilities.NewLogger().LogWithCaller()

	value, found := os.LookupEnv(key)
	if !found {
		logger.Fatalf("environment variable %s required, but not found", key)
	}
	environment[key] = value
}

// Registers the environment variable in the registry. If not found, the value will be the passed-in fallback
func (r *EnvLoader) RegisterDefault(key, fallback string) {
	logger := utilities.NewLogger().LogWithCaller()

	if value, found := os.LookupEnv(key); found {
		environment[key] = value
	} else {
		logger.Warningf("environment variable %s not found. Using fallback value", key)
		environment[key] = fallback
	}
}

// Retrieves the passed-in variable from the registry, returning an error if the variable was not found
func (r *EnvLoader) Get(key string) (string, error) {
	logger := utilities.NewLogger().LogWithCaller()
	if value, found := environment[key]; found {
		return value, nil
	}
	logger.Errorf("could not find environment variable: %s", key)
	return "", fmt.Errorf("could not find environment variable: %s. Did you remember to register it?", key)
}
