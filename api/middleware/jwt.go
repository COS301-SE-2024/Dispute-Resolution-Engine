package middleware

// Add JWT-related imports
import (
	"api/env"
	"api/models"
	"api/redisDB"
	"api/utilities"
	"errors"
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

type JwtMiddleware struct {
	EnvLoader env.Env
	Logger    *utilities.Logger
}

var jwtMiddleware Jwt
var once sync.Once

func NewJwtMiddleware() Jwt {
	once.Do(func() {
		jwtMiddleware = createJwtMiddleware()
	})
	return jwtMiddleware
}

func createJwtMiddleware() Jwt {
	return &JwtMiddleware{
		EnvLoader: env.NewEnvLoader(),
		Logger:    utilities.NewLogger().LogWithCaller(),
	}
}

// GenerateJWT generates a JWT token
// GenerateJWT generates a JWT token for the given user
func (j *JwtMiddleware) GenerateJWT(user models.User) (string, error) {
	jwtSec, err := j.EnvLoader.Get("JWT_SECRET")
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
		j.Logger.WithError(err).Error("Failed to sign token")
		return "", err
	}

	j.StoreJWT(user.Email, signedToken)
	j.Logger.Info("Token signed successfully")
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
	logger := j.Logger
	envLoader := j.EnvLoader

	// Check for Orchestrator Header
	orchestratorKey := c.GetHeader("X-Orchestrator-Key")
	expectedOrchestratorKey, err := envLoader.Get("ORCHESTRATOR_KEY")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong"})
		return
	}

	// If the request comes from the orchestrator, skip JWT authentication
	if orchestratorKey == expectedOrchestratorKey {
		logger.Info("Request from orchestrator, skipping JWT authentication")
		c.Next()
		return
	}

	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		logger.Error("No Authorization header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		logger.Error("Invalid Authorization header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if tokenString == "" {
		logger.Error("No token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	jwtSecretKey, err := envLoader.Get("JWT_SECRET")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{Error: "Something went wrong"})
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		logger.WithError(err).Error("Error parsing token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
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
func (j *JwtMiddleware) GetClaims(c *gin.Context) (models.UserInfoJWT, error) {
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
		return models.UserInfoJWT{}, errors.New("Missing Authorization header")
	}

	tokenString, _ := strings.CutPrefix(authHeader, "bearer")
	tokenString, _ = strings.CutPrefix(tokenString, "Bearer")
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		logger.Error("No token")
		return models.UserInfoJWT{}, errors.New("Missing Authorization header")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		logger.WithError(err).Error("Error parsing token")
		return models.UserInfoJWT{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		logger.Error("Error getting claims")
		return models.UserInfoJWT{}, errors.New("failed to get claims")
	}
	logger.Info("Claims retrieved successfully")
	return claims.User, nil
}
