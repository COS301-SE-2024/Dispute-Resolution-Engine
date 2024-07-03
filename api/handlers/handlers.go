package handlers

import (
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
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

func new(db *gorm.DB) Handler {
	return Handler{DB: db}
}

func NewAuthHandler(db *gorm.DB) Auth {
	return Auth{new(db)}
}

func NewUserHandler(db *gorm.DB) User {
	return User{new(db)}
}

func NewDisputeHandler(db *gorm.DB) Dispute {
	return Dispute{new(db)}
}

func NewUtilitiesHandler(db *gorm.DB) Utility {
	return Utility{new(db)}
}
