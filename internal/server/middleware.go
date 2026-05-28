package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const userIDContextKey = "userID"

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeAPIError(c, http.StatusUnauthorized, "missing_token", "Bearer token manquant")
			c.Abort()
			return
		}

		parsedClaims, err := parseClaims(strings.TrimPrefix(header, "Bearer "), s.cfg.JWTSecret)
		if err != nil {
			writeAPIError(c, http.StatusUnauthorized, "invalid_token", "Session invalide ou expirée")
			c.Abort()
			return
		}

		c.Set(userIDContextKey, parsedClaims.UserID)
		c.Next()
	}
}

func currentUserID(c *gin.Context) int64 {
	value, _ := c.Get(userIDContextKey)
	userID, _ := value.(int64)
	return userID
}
