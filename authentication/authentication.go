package authentication

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is the middleware for authentication and authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.String(http.StatusUnauthorized, "Token is Misssing")
			c.Abort()
			return
		}
		for index, char := range tokenString {
			if char == ' ' {
				tokenString = tokenString[index+1:]
			}
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			c.String(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}
		// Check whether token is expired or not
		check, ok := claims["exp"].(int64)
		if ok && check < time.Now().Unix() {
			c.String(http.StatusUnauthorized, "Expired token")
			c.Abort()
			return
		}

		c.Set("email", claims["email"])
		c.Set("role", claims["role"])
	}
}

// AdminAuth verifies if the user is an admin
func AdminAuth(c *gin.Context) error {
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		return errors.New("only admins can access this endpoint")
	}
	return nil
}

// UserAuth verifies if the user is a regular user
func UserAuth(c *gin.Context) error {
	role, exists := c.Get("role")
	if !exists || role.(string) != "user" {
		return errors.New("user admins can access this endpoint")
	}
	return nil
}

// CommonAuth verifies if the user is either an admin or a regular user
func CommonAuth(c *gin.Context) error {
	role, exists := c.Get("role")
	if !exists || (role.(string) != "user" && role.(string) != "admin") {
		return errors.New("only users and admins have access to this endpoint")
	}
	return nil
}
