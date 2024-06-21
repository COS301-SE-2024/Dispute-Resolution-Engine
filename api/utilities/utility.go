package utilities

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

func (a *Argon2idHash) HashPassword(password string) *HashSalt {
	salt, err := RandomSalt(16)
	if err != nil {
		return nil
	}
	hashSalt, err := a.GenerateHash([]byte(password), salt)
	if err != nil {
		return nil
	}
	return hashSalt
}

type Argon2idHash struct {
	time    uint32
	memory  uint32
	threads uint8
	keylen  uint32
	saltlen uint32
}

type HashSalt struct {
	Hash []byte
	Salt []byte
}

func NewArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keylen uint32) *Argon2idHash {
	return &Argon2idHash{time, memory, threads, keylen, saltLen}
}

func RandomSalt(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (a *Argon2idHash) GenerateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	if len(salt) == 0 {
		salt, err = RandomSalt(a.saltlen)
	}
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keylen)
	return &HashSalt{Hash: hash, Salt: salt}, nil

}

func (a *Argon2idHash) Compare(hash, salt, password []byte) bool {
	// Generate hash for comparison.
	hashSalt, err := a.GenerateHash(password, salt)
	if err != nil {
		return false
	}
	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(hash, hashSalt.Hash) {
		return false
	}
	return true
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
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
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
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
		return err
	}
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
