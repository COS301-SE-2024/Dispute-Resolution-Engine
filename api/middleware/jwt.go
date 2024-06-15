package middleware

// Add JWT-related imports
import (
	"api/models"
	"api/utilities"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// JWT secret key
var jwtSecretKey = []byte("your_secret_key")

// Claims struct to store user data in JWT
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
	User models.User `json:"user"`
}

// GenerateJWT generates a JWT token
func GenerateJWT(user models.User) (string, error) {
    claims := &Claims{
        Email: user.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
            IssuedAt:  time.Now().Unix(),
        },
		User: user,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecretKey)
}

// JWTMiddleware is a middleware to validate JWT token
func JWTMiddleware(next http.Handler) http.Handler {
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
