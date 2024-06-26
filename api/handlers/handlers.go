package handlers

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Handler struct {
	DB  *gorm.DB
	RDB *redis.Client
}

type Auth struct {
	Handler
}

type User struct {
	Handler
}

type Dispute struct {
	Handler
}

type Utility struct {
	Handler
}

func new(db *gorm.DB, rdb *redis.Client) Handler {
	return Handler{db, rdb}
}

func NewAuthHandler(db *gorm.DB, rdb *redis.Client) Auth {
	return Auth{new(db, rdb)}
}

func NewUserHandler(db *gorm.DB, rdb *redis.Client) User {
	return User{new(db, rdb)}
}

func NewDisputeHandler(db *gorm.DB, rdb *redis.Client) Dispute {
	return Dispute{new(db, rdb)}
}

func NewUtilitiesHandler(db *gorm.DB, rdb *redis.Client) Utility {
	return Utility{new(db, rdb)}
}
