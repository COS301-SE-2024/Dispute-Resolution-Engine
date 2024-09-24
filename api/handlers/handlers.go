package handlers

import (
	"api/auditLogger"
	"api/middleware"

	"api/env"

	"gorm.io/gorm"
)

type Handler struct {
	DB                       *gorm.DB
	EnvReader                env.Env
	Jwt                      middleware.Jwt
	DisputeProceedingsLogger auditLogger.DisputeProceedingsLoggerInterface
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
type Archive struct {
	Handler
}

type Ticket struct {
	Handler
}

func new(db *gorm.DB) Handler {
	envReader := env.NewEnvLoader()
	return Handler{DB: db, EnvReader: envReader, Jwt: middleware.NewJwtMiddleware(), DisputeProceedingsLogger: auditLogger.NewDisputeProceedingsLogger(db, envReader)}
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

func NewArchiveHandler(db *gorm.DB) Archive {
	return Archive{new(db)}
}

func NewTicketHandler(db *gorm.DB) Ticket {
	return Ticket{new(db)}
}
