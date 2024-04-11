package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			// TODO: 不使用固定密钥，改为从配置文件中读取
			return []byte("%YU3#oZvc*e%23vN"), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Error parsing token", "details": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			role := (*claims)["role"].(string)
			if !contains(requiredRoles, role) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions", "required_roles": requiredRoles, "token_role": role})
				c.Abort()
				return
			}
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
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
