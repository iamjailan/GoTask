package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	CustomerIDKey   = "customer_id"
	CustomerRoleKey = "customer_role"
)

func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const bearerPrefix = "Bearer "
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, bearerPrefix) {
			unauthorized(c)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
		if tokenString == "" {
			unauthorized(c)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			unauthorized(c)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		customerID, idOK := claims["sub"].(string)
		if !ok || !idOK || customerID == "" {
			unauthorized(c)
			return
		}
		c.Set(CustomerIDKey, customerID)
		if role, ok := claims["role"].(string); ok {
			c.Set(CustomerRoleKey, role)
		}
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
}
