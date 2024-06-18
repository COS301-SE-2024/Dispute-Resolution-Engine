package handlers

import "gorm.io/gorm"

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{db}
}
