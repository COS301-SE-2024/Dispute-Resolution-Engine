package middleware

// Add JWT-related imports
import (
	"api/models"
	"api/redisDB"
	"net/http"
	"os"
	"strings"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Claims struct to store user data in JWT
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
	User models.User `json:"user"`
}

// GenerateJWT generates a JWT token
// GenerateJWT generates a JWT token for the given user
func GenerateJWT(user models.User) (string, error) {
	jwtSec := os.Getenv("JWT_SECRET")

	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{ // Use a unique identifier for the user
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
		User: user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSec))
	if err != nil {
		return "", err
	}

	// Store the token somewhere (for example, in memory map activeTokens)
	StoreJWT(user.Email, signedToken)

	return signedToken, nil
}

func StoreJWT(email string, jwt string) error {
	err := redisDB.RDB.Set(context.Background(), email, jwt, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetJWT(userEmail string) (string, error) {
	jwt, err := redisDB.RDB.Get(context.Background(), userEmail).Result()
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func JWTMiddleware(c *gin.Context) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Check if the Authorization header starts with "Bearer "
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Extract token from "Bearer <token>"
	tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Get JWT secret key
	jwtSecretKey := []byte(os.Getenv("JWT_SECRET"))

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}

	// Validate token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		//ctx := context.WithValue(r.Context(), "user", claims)
		userEmail := claims.Email
		jwtFromDB, err := GetJWT(userEmail)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
		if jwtFromDB != tokenString {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
        c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
		return
	}
}

// return claims
func GetClaims(c *gin.Context) *Claims {
	if jwtSecretKey := os.Getenv("JWT_SECRET"); jwtSecretKey == "" {
		err := godotenv.Load("api.env")
		if err != nil {
			return nil
		}
	}
	secret := []byte(os.Getenv("JWT_SECRET"))

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil
	}

	tokenString := strings.Split(authHeader, " ")[1]
	if tokenString == "" {
		return nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil
	}
	return claims
}
