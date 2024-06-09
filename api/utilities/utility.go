package utilities

import (
	"bytes"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/argon2"
)

func (a *Argon2idHash) hashPassword(password string) *HashSalt {
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

func (a *Argon2idHash) Compare(hash, salt, password []byte) error {
	// Generate hash for comparison.
	hashSalt, err := a.GenerateHash(password, salt)
	if err != nil {
		return err
	}
	// Compare the generated hash with the stored hash.
	// If they don't match return error.
	if !bytes.Equal(hash, hashSalt.Hash) {
		return errors.New("hash doesn't match")
	}
	return nil
}