package utilities

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// Retrieves an environment variable, producing an error if the variable is not found
func GetRequiredEnv(key string) (string, error) {
	value, found := os.LookupEnv(key)
	if !found {
		return "", fmt.Errorf("environment variable %s required, but not found", key)
	}
	return value, nil
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
