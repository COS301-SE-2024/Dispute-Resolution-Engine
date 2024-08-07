package handlers

import (
	"api/middleware"

	"gorm.io/gorm"
	"api/env"

)

type Handler struct {
	DB *gorm.DB
	EnvReader env.Env
	jwt middleware.Jwt
}

type Auth struct {
	Handler
}

type User struct {
	Handler
}

type Utility struct {
	Handler
}

type Notification struct {
	Handler
}

type Expert struct {
	Handler
}

func new(db *gorm.DB) Handler {
	return Handler{DB: db, EnvReader: env.NewEnvLoader(), jwt: middleware.NewJwtMiddleware()}
}

func NewAuthHandler(db *gorm.DB) Auth {
	return Auth{new(db)}
}

func NewUserHandler(db *gorm.DB) User {
	return User{new(db)}
}

func NewUtilitiesHandler(db *gorm.DB) Utility {
	return Utility{new(db)}
}

func NewNotificationHandler(db *gorm.DB) Notification {
	return Notification{new(db)}
}

func NewExpertHandler(db *gorm.DB) Expert {
	return Expert{new(db)}
}
