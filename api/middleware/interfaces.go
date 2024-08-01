package middleware

import (
	"api/models"

	"github.com/gin-gonic/gin"
)

type Jwt interface {
	GenerateJWT(user models.User) (string, error)
	StoreJWT(email string, jwt string) error
	GetJWT(email string) (string, error)
	JWTMiddleware(c *gin.Context)

	GetClaims(c *gin.Context) (models.UserInfoJWT, error)
}

type Role interface {
	matchKeyToValue(value string) (int, bool)
	RoleMiddleware(reqAuthlevel int) gin.HandlerFunc
}

