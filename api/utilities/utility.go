package utilities

import (
	"api/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonSaltlen = 12288
	argonMemory  = 4
	argonThreads = 32
	argonKeylen  = 16
)

// Hashes the passed-in password with a randomly generated salt
func HashPassword(password string) (hash, salt []byte, err error) {
	salt, err = RandomSalt(16)
	if err != nil {
		return
	}
	hash = HashPasswordWithSalt(password, salt)
	return hash, salt, err
}

// Hashes a password with the passed-in salt
func HashPasswordWithSalt(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeylen)
}

// Generates a random salt of the specified length
func RandomSalt(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func Compare(hash, salt, password []byte) bool {
	// Generate hash for comparison.
	passHash := HashPasswordWithSalt(string(password), salt)

	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(hash, passHash) {
		return false
	}
	return true
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func GetCurrentTime() time.Time {
	return time.Now()
}

func GetCurrentTimePtr() *time.Time {
	t := time.Now()
	return &t
}

func GenerateVerifyEmailToken() string {
	rand.Seed(time.Now().UnixNano())
	const length = 6 // Define the length of the token
	token := make([]byte, length)
	for i := 0; i < length; i++ {
		token[i] = byte('0' + rand.Intn(10))
	}
	return string(token)
}

func WriteToFile(data string, filepath string) error {
	// Open the file with append and create options, set the permission to 0644
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// Ensure the file is closed when the function exits
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	// Write the data with a newline prefix
	_, err = file.WriteString("\n" + data)
	if err != nil {
		print("Failed to write to file")
		return err
	}
	print("Successfully written to file")
	return nil
}

func ReadFromFile(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func RemoveFromFile(filepath string, data string) (bool, error) {
	lines, err := ReadFromFile(filepath)
	if err != nil {
		return false, err
	}
	if len(lines) == 0 {
		return false, nil
	}
	if ContainsString(lines, data) {
		var newLines []string
		for _, line := range lines {
			if line != data {
				newLines = append(newLines, line)
			}
		}
		err := os.WriteFile(filepath, []byte(strings.Join(newLines, "\n")), 0644)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func ContainsString(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func InternalError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong"})
}
