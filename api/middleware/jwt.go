package middleware

// Add JWT-related imports
import (
	"api/env"
	"api/models"
	"api/redisDB"
	"api/utilities"
	"net/http"
	"strings"
	"sync"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Claims struct to store user data in JWT
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
	User models.UserInfoJWT `json:"user"`
}

type JwtMiddleware struct{}

var jwtMiddleware Jwt
var once sync.Once

func NewJwtMiddleware() Jwt {
	once.Do(func() {
		jwtMiddleware = createJwtMiddleware()
	})
	return jwtMiddleware
}

func createJwtMiddleware() Jwt {
	return &JwtMiddleware{}
}



// GenerateJWT generates a JWT token
// GenerateJWT generates a JWT token for the given user
func GenerateJWT(user models.User) (string, error) {
	envReader := env.NewEnvLoader()
	logger := utilities.NewLogger().LogWithCaller()
	jwtSec, err := envReader.Get("JWT_SECRET")
	if err != nil {
		return "", err
	}

	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{ // Use a unique identifier for the user
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
		User: *models.ConvertUserToJWTUser(user),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSec))
	if err != nil {
		logger.WithError(err).Error("Failed to sign token")
		return "", err
	}

	j.StoreJWT(user.Email, signedToken)
	logger.Info("Token signed successfully")
	return signedToken, nil
}

func (j *JwtMiddleware) StoreJWT(email string, jwt string) error {
	logger := utilities.NewLogger().LogWithCaller()
	err := redisDB.RDB.Set(context.Background(), email, jwt, 24*time.Hour).Err()
	if err != nil {
		logger.WithError(err).Error("Failed to store JWT in Redis")
		return err
	}
	logger.Info("JWT stored successfully")
	return nil
}

func (j *JwtMiddleware) GetJWT(userEmail string) (string, error) {
	logger := utilities.NewLogger().LogWithCaller()
	jwt, err := redisDB.RDB.Get(context.Background(), userEmail).Result()
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve JWT from Redis")
		return "", err
	}
	logger.Info("JWT retrieved successfully")
	return jwt, nil
}

func (j *JwtMiddleware) JWTMiddleware(c *gin.Context) {
	logger := utilities.NewLogger().LogWithCaller()
	envLoader := env.NewEnvLoader()

	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		logger.Error("No Authorization header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Check if the Authorization header starts with "Bearer "
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		logger.Error("Invalid Authorization header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Extract token from "Bearer <token>"
	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if tokenString == "" {
		logger.Error("No token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Get JWT secret key
	var jwtSecretKey []byte
	{
		secret, err := envLoader.Get("JWT_SECRET")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong"})
		}
		jwtSecretKey = []byte(secret)
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		logger.WithError(err).Error("Error parsing token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Validate token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		//ctx := context.WithValue(r.Context(), "user", claims)
		userEmail := claims.Email
		jwtFromDB, err := j.GetJWT(userEmail)
		if err != nil {
			logger.WithError(err).Error("Couldn't retrieve JWT from Redis")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
		if jwtFromDB != tokenString {
			logger.Error("Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
		c.Next()
	} else {
		logger.Error("Invalid token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}
}

// return claims
func (j *JwtMiddleware) GetClaims(c *gin.Context) *Claims {
	logger := utilities.NewLogger().LogWithCaller()
	envLoader := env.NewEnvLoader()
	var secret []byte
	{
		s, err := envLoader.Get("JWT_SECRET")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong"})
		}
		secret = []byte(s)
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		logger.Error("No Authorization header")
		return nil
	}

	tokenString := strings.Split(authHeader, " ")[1]
	if tokenString == "" {
		logger.Error("No token")
		return nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		logger.WithError(err).Error("Error parsing token")
		return nil
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		logger.Error("Error getting claims")
		return nil
	}
	logger.Info("Claims retrieved successfully")
	return claims
}
