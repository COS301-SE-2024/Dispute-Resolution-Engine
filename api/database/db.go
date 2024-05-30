package database

import (
	"fmt"
)

type DB struct {
	data map[string]string
}

func NewDB() *DB {
	return &DB{
		data: make(map[string]string),
	}
}

func (db *DB) Get(key string) (string, error) {
	value, ok := db.data[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	return value, nil
}

func (db *DB) Set(key, value string) {
	db.data[key] = value
}

func main() {
	// Create a new instance of the mock database
	db := NewDB()

	// Set some values
	db.Set("key1", "value1")
	db.Set("key2", "value2")

	// Get a value
	value, err := db.Get("key1")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value)
	}
}