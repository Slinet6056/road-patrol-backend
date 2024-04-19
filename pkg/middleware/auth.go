package middleware

import (
	"net/http"
	"strings"

	"github.com/Slinet6056/road-patrol-backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Error parsing token", "details": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			if role, ok := (*claims)["role"].(string); ok {
				if !contains(requiredRoles, role) {
					c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Insufficient permissions", "required_roles": requiredRoles, "token_role": role})
					c.Abort()
					return
				}
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Role claim must be a non-empty string"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token"})
			c.Abort()
			return
		}
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
