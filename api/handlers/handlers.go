package handlers

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	RDB *redis.Client
}

func New(db *gorm.DB, rdb *redis.Client) Handler {
	return Handler{db, rdb}
}
