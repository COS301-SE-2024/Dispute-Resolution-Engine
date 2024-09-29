package adminanalytics

import (
	"api/env"
	"api/middleware"

	"gorm.io/gorm"
)

type AdminAnalyticsDBModel interface {
}

type AdminAnalyticsHandler struct {
	DB        AdminAnalyticsDBModel
	EnvReader env.Env
	JWT       middleware.Jwt
}

type AdminAnalyticsDBModelReal struct {
	DB *gorm.DB
	env env.Env
}

func NewAdminAnalyticsHandler(db *gorm.DB, envReader env.Env) AdminAnalyticsHandler{
	return AdminAnalyticsHandler{
		DB: 	  AdminAnalyticsDBModelReal{DB: db, env: envReader},
		EnvReader: envReader,
		JWT:      middleware.NewJwtMiddleware(),
	}
}

