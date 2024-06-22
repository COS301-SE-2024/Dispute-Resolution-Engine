package middleware

// Add JWT-related imports
import (
	"api/models"
	"api/redisDB"
	"api/utilities"
	"net/http"
	"os"
	"strings"
	"time"

	"context"

	"github.com/dgrijalva/jwt-go"
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

	claims := &jwt.StandardClaims{
		Subject:   user.Email,                            // Use a unique identifier for the user
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		IssuedAt:  time.Now().Unix(),
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

func StoreJWT(userEmail string, jwt string) error {
	err := redisDB.RDB.Set(context.Background(), userEmail, jwt, 24*time.Hour).Err()
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

// JWTMiddleware is a middleware to validate JWT token
func JWTMiddleware(next http.Handler) http.Handler {
	//stub function
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//     context.Set(r, "user", &Claims{User: models.User{ID: 1, Email: "temp@temp.com"}})
	//     next.ServeHTTP(w, r)
	// })

	//end of stub function
	if jwtSecretKey := os.Getenv("JWT_SECRET"); jwtSecretKey == "" {
		err := godotenv.Load("api.env")
		if err != nil {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				utilities.WriteJSON(w, http.StatusInternalServerError, models.Response{Error: "Error loading environment variables"})
			})
		}
	}

	jwtSecretKey := []byte(os.Getenv("JWT_SECRET"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			utilities.WriteJSON(w, http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
		tokenString := strings.Split(authorizationHeader, " ")[1] // Extract token from "Bearer <token>"
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})
		if err != nil {
			utilities.WriteJSON(w, http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			context.Set(r, "user", claims)
			next.ServeHTTP(w, r)
		} else {
			utilities.WriteJSON(w, http.StatusUnauthorized, models.Response{Error: "Unauthorized"})
			return
		}
	})
}

// return claims
func GetClaims(r *http.Request) *Claims {
	if jwtSecretKey := os.Getenv("JWT_SECRET"); jwtSecretKey == "" {
		err := godotenv.Load("api.env")
		if err != nil {
			return nil
		}
	}
	secret := []byte(os.Getenv("JWT_SECRET"))

	authHeader := r.Header.Get("Authorization")
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
