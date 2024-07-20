package utilities

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	hashTime    = 1
	hashMemory  = 12288
	hashThreads = 4
	hashKeylen  = 32
	hashSaltlen = 16
)

// Hashes the passed-in password with a randomly generated salt
func HashPassword(password string) (hash, salt []byte, err error) {
	salt, err = RandomSalt(hashSaltlen)
	if err != nil {
		return
	}
	hash = HashPasswordWithSalt(password, salt)
	return hash, salt, nil
}

// Hashes a password with the passed-in salt
func HashPasswordWithSalt(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, hashTime, hashMemory, hashThreads, hashKeylen)
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
